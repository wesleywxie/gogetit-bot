package task

import (
	"github.com/wesleywxie/gogetit/internal/model"
	"go.uber.org/atomic"
)

// LivestreamUpdateTask Livestream 更新任务
type LivestreamUpdateTask struct {
	observerList []LivestreamUpdateObserver
	isStop       atomic.Bool
}

// NewLivestreamTask new NewLivestreamTask
func NewLivestreamTask() *LivestreamUpdateTask {
	return &LivestreamUpdateTask{
		observerList: []LivestreamUpdateObserver{},
	}
}

// LivestreamUpdateObserver Livestream update observer
type LivestreamUpdateObserver interface {
	update([]*model.Subscription)
	errorUpdate(*model.Subscription)
	id() string
}
