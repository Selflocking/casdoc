package casdoc

import (
	"fmt"
	"path"
	"strings"
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestPolishDocs(t *testing.T) {
	err := q.getFileList(path.Join(repoPath, "/docs/"))
	if err != nil {
		panic(err)
	}

	//var err error
	//q.push("/home/yunshu/Studio/GitHub/casdoor-website/docs/basic/core-concepts.md")

	counter := 0
	totalItems := len(q.item)

	for {
		f := q.pop()
		counter++
		logger = log.WithField("rate", fmt.Sprintf("%d/%d", counter, totalItems))

		logger.Info("now polish: ", strings.TrimPrefix(f, repoPath))
		err = polish(f)

		if err != nil {
			q.addToFailedList(f)
			logger.Errorf("error: %v\n", err)
		}
		if q.empty() {
			break
		}
	}
}

func TestTranslateDocs(t *testing.T) {
	langs := []string{"zh", "fr", "de", "ko", "ru", "ja"}

	err := q.getFileList(path.Join(repoPath, "/docs/"))
	if err != nil {
		panic(err)
	}

	counter := 0
	totalItems := len(q.item)

	for _, lang := range langs {
		logger = log.WithField("lang", lang)
		for {
			f := q.pop()
			counter++
			logger = logger.WithField("rate", fmt.Sprintf("%d/%d", counter, totalItems))

			logger.Info("now translate: ", strings.TrimPrefix(f, repoPath))
			err = translate(f, lang)

			if err != nil {
				q.addToFailedList(f)
				logger.Errorf("error: %v\n", err)
			}
			if q.empty() {
				break
			}
		}
	}
}
