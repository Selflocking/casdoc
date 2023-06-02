package casdoc

import (
	"github.com/sashabaranov/go-openai"
	"os"
)

// OPENAI
var authToken = ""

const RPM = 3
const TPM = 40000

// repo path
const repoPath = "/home/yunshu/Studio/GitHub/casdoor-website/"

func init() {
	if authToken == "" {
		key, exist := os.LookupEnv("OPENAI_SECRET_KEY")
		if exist {
			authToken = key
		} else {
			panic("请将token添加到环境变量中或者硬编码到程序中")
		}
	}
}

// i18n 目录
const i18nPathPrefix = "/i18n/"
const i18nPathSuffix = "/docusaurus-plugin-content-docs/current/"

// /i18n/%two_letters_code%/docusaurus-plugin-content-docs/current/**/%original_file_name%
// 润色 Prompt
const polishPrompt = `You are given an mdx document written in English, and you are tasked with polishing it by correcting typos, improving sentence fluency, and so on, without changing the structure or style of the document. Please only return the polished document.`

var polishRequest = openai.ChatCompletionRequest{
	Model: openai.GPT3Dot5Turbo,
	Messages: []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: polishPrompt,
		},
		{
			Role: openai.ChatMessageRoleUser,
			Content: `:::note

content

:::`,
		},
		{
			Role: openai.ChatMessageRoleAssistant,
			Content: `:::note

polished content

:::`,
		},
	},
}

// 翻译 prompt
// 中文
const chinesePrompt = `Translate the given document in mdx format into Chinese. Please be careful not to modify the structure of the document and do not translate technical terms such as Casbin, Casdoor, SSO, Swagger, URL, etc. Only provide the final result.`

var chineseRequest = openai.ChatCompletionRequest{
	Model: openai.GPT3Dot5Turbo,
	Messages: []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: chinesePrompt,
		},
	},
}

// Français
const frenchPrompt = ``

var frenchRequest = openai.ChatCompletionRequest{
	Model: openai.GPT3Dot5Turbo,
	Messages: []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: frenchPrompt,
		},
	},
}

// Deutsch
const germanPrompt = ``

var germanRequest = openai.ChatCompletionRequest{
	Model: openai.GPT3Dot5Turbo,
	Messages: []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: germanPrompt,
		},
	},
}

// 한국어
const koreanPrompt = ``

var koreanRequest = openai.ChatCompletionRequest{
	Model: openai.GPT3Dot5Turbo,
	Messages: []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: koreanPrompt,
		},
	},
}

// Русский
const russianPrompt = ``

var russianRequest = openai.ChatCompletionRequest{
	Model: openai.GPT3Dot5Turbo,
	Messages: []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: russianPrompt,
		},
	},
}

// 日本語
const japanesePrompt = ``

var japaneseRequest = openai.ChatCompletionRequest{
	Model: openai.GPT3Dot5Turbo,
	Messages: []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: japanesePrompt,
		},
	},
}
