package task

import (
	"github.com/wesleywxie/gogetit/internal/config"
	"github.com/wesleywxie/gogetit/internal/model"
	"go.uber.org/atomic"
	"go.uber.org/zap"
	"time"
)

func init() {
	task := NewLivestreamTask()
	task.Register(&telegramBotUpdateObserver{})
	registerTask(task)
}

// Register 注册rss更新订阅者
func (t *LivestreamUpdateTask) Register(observer LivestreamUpdateObserver) {
	t.observerList = append(t.observerList, observer)
}

// Deregister 注销rss更新订阅者
func (t *LivestreamUpdateTask) Deregister(removeObserver LivestreamUpdateObserver) {
	for i, observer := range t.observerList {
		if observer.id() == removeObserver.id() {
			t.observerList = append(t.observerList[:i], t.observerList[i+1:]...)
			return
		}
	}
}

// Name 任务名称
func (t *LivestreamUpdateTask) Name() string {
	return "LivestreamUpdateTask"
}

// Stop stop task
func (t *LivestreamUpdateTask) Stop() {
}

// Start run task
func (t *LivestreamUpdateTask) Start() {
	if config.RunMode == config.TestMode {
		return
	}

	t.isStop.Store(false)

	go func() {
		for {
			if t.isStop.Load() == true {
				zap.S().Info("LivestreamUpdateTask stopped")
				return
			}

			time.Sleep(time.Duration(config.UpdateInterval) * time.Minute)
		}
	}()
}

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
	update(*model.Subscription)
	errorUpdate(*model.Subscription)
	id() string
}

type telegramBotUpdateObserver struct {
}

func (o *telegramBotUpdateObserver) update(subscription *model.Subscription) {
	zap.S().Debugf("%v receiving [%d]%v update", o.id(), subscription.ID, subscription.KOL)
}

func (o *telegramBotUpdateObserver) errorUpdate(subscription *model.Subscription) {
	zap.S().Debugf("%v receiving [%d]%v error update", o.id(), subscription.ID, subscription.KOL)
}

func (o *telegramBotUpdateObserver) id() string {
	return "telegramBotUpdateObserver"
}
