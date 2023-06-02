package casdoc

import (
	"context"
	"fmt"
	rate "golang.org/x/time/rate"
	"io/fs"
	"path/filepath"
)

type workQueue struct {
	item           []string
	tokenLimiter   *rate.Limiter
	requestLimiter *rate.Limiter
	failed         []string
}

// 获得文件列表
func (q *workQueue) getFileList(docDir string) error {
	err := filepath.Walk(docDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Println("error!")
		}
		var name = info.Name()
		var size = len(name)
		if !info.IsDir() && (name[size-3:] == "mdx" || name[size-2:] == "md") {
			q.push(path)
		}
		return err
	})
	return err
}

// 进入任务队列
func (q *workQueue) push(path string) {
	q.item = append(q.item, path)
}

// 弹出任务
func (q *workQueue) pop() string {
	removed := q.item[0]
	q.item = q.item[1:]
	return removed
}

// 等待合适时间后在执行
func (q *workQueue) wait(ctx context.Context, tokenNumber int) (err error) {
	err = q.requestLimiter.Wait(ctx)
	if err != nil {
		return
	}
	err = q.tokenLimiter.WaitN(ctx, tokenNumber)
	if err != nil {
		return
	}
	return nil
}

// 任务是否为空
func (q *workQueue) empty() bool {
	return len(q.item) == 0
}

// 队列头
func (q *workQueue) front() string {
	return q.item[0]
}

// 队列尾
func (q *workQueue) back() string {
	return q.item[len(q.item)-1]
}

// 添加到失败列表
func (q *workQueue) addToFailedList(path string) {
	q.failed = append(q.failed, path)
}
