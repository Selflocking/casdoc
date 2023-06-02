package casdoc

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	openai "github.com/sashabaranov/go-openai"
	log "github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

var client *openai.Client
var q workQueue
var logger *log.Entry

func init() {
	client = openai.NewClient(authToken)
	q = workQueue{}
	requestLimit := rate.Every(time.Minute / RPM)
	tokenLimit := rate.Every(time.Minute / TPM)
	q.requestLimiter = rate.NewLimiter(requestLimit, RPM)
	q.tokenLimiter = rate.NewLimiter(tokenLimit, TPM)
}

// 得到文档内容
func getDocContext(path string) *string {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	res := string(b)
	return &res
}

func gpt(req openai.ChatCompletionRequest, c *string) (ans *string, totalTokens int, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	req.Messages = append(req.Messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: *c,
	})
	resp, err := client.CreateChatCompletion(ctx, req)

	if err != nil {
		logger.Errorf("chat failed\n")
		return
	}

	ans = &resp.Choices[0].Message.Content
	totalTokens = resp.Usage.TotalTokens
	return
}

// 文档润色
func polish(path string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	docContext := getDocContext(path)
	/*
	*	refer to : https://platform.openai.com/tokenizer
	*	A helpful rule of thumb is that one token generally corresponds to ~4 characters
	*	of text for common English text. This translates to roughly ¾ of a word (so 100
	*	tokens ~= 75 words).
	**/
	tokenNumber := len(*docContext) / 4

	/*
	*	https://platform.openai.com/docs/models/gpt-3-5
	*	gpt-3.5-turbo model's maximum context length is 4096 tokens.
	**/
	var userContext []*string
	if tokenNumber > 2048 {
		// split doc context by h2
		strArr := strings.Split(*docContext, "\n## ")
		strArr[0] = strArr[0] + "\n## " + strArr[1]
		userContext = append(userContext, &strArr[0])
		for i := 2; i < len(strArr); i++ {
			strArr[i] = "## " + strArr[i]
			userContext = append(userContext, &strArr[i])
		}
	} else {
		userContext = append(userContext, docContext)
	}

	var polishedDoc string
	realTokenNumber := 0

	for _, c := range userContext {
		tokenNumber := len(*c) / 4
		logger.Info("cost token: ", tokenNumber)
		err := q.wait(ctx, tokenNumber)
		if err != nil {
			logger.Errorf("failed to polish: %s\n", path)
			return err
		}
		ans, totalTokens, err := gpt(polishRequest, c)
		if err != nil {
			logger.Errorf("failed to polish: %s\n", path)
			return err
		}
		polishedDoc += *ans + "\n"
		realTokenNumber += totalTokens
	}

	logger.Info("total token: ", realTokenNumber)

	err := os.WriteFile(path, []byte(polishedDoc), 0644)

	if err != nil {
		logger.Errorf("unable to write file: %s \n", path)
		return err
	}
	return nil
}

func translate(docPath string, lang string) error {
	var req openai.ChatCompletionRequest
	switch lang {
	case "zh": // 中文
		req = chineseRequest
	case "fr": // Français
		req = frenchRequest
	case "de": // Deutsch
		req = germanRequest
	case "ko": // 한국어
		req = koreanRequest
	case "ru": // Русский
		req = russianRequest
	case "ja": // 日本語
		req = japaneseRequest
	default:
		return errors.New(fmt.Sprint("unknown language: ", lang))
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	docContext := getDocContext(docPath)
	/*
		refer to : https://platform.openai.com/tokenizer
		A helpful rule of thumb is that one token generally corresponds to ~4 characters
		of text for common English text. This translates to roughly ¾ of a word (so 100
		tokens ~= 75 words).
	*/
	tokenNumber := len(*docContext) / 4

	/*
	*	https://platform.openai.com/docs/models/gpt-3-5
	*	gpt-3.5-turbo model's maximum context length is 4096 tokens.
	**/
	var userContext []*string
	if tokenNumber > 2048 {
		logger.Info("Large doc, split it by h2")

		// split doc context by h2
		strArr := strings.Split(*docContext, "\n## ")
		strArr[0] = strArr[0] + "\n## " + strArr[1]
		userContext = append(userContext, &strArr[0])
		for i := 2; i < len(strArr); i++ {
			strArr[i] = "## " + strArr[i]
			userContext = append(userContext, &strArr[i])
		}
		logger.Info("Split doc into ", len(userContext), " parts")
	} else {
		userContext = append(userContext, docContext)
	}

	var polishedDoc string
	realTokenNumber := 0

	for _, c := range userContext {
		tokenNumber := len(*c) / 4
		logger.Info("cost token: ", tokenNumber)
		err := q.wait(ctx, tokenNumber)
		if err != nil {
			logger.Errorf("failed to polish: %s\n", docPath)
			return err
		}

		ans, totalTokens, err := gpt(req, c)
		if err != nil {
			logger.Errorf("failed to polish: %s\n", docPath)
			return err
		}
		polishedDoc += *ans + "\n"
		realTokenNumber += totalTokens
	}

	logger.Info("total token: ", realTokenNumber)

	// 相对路径英文单词做变量名
	relativePath := strings.Split(docPath, path.Join(repoPath, "/docs/"))[1]
	translatedDocPath := path.Join(repoPath, i18nPathPrefix, lang, i18nPathSuffix, relativePath)

	_ = os.MkdirAll(filepath.Dir(translatedDocPath), 0755)
	err := os.WriteFile(translatedDocPath, []byte(polishedDoc), 0644)
	if err != nil {
		logger.Errorf("unable to write file: %s", docPath)
		return err
	}
	return nil
}
