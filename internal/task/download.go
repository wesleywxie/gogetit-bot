package task

import (
	"github.com/wesleywxie/gogetit/internal/config"
	"github.com/wesleywxie/gogetit/internal/model"
	"go.uber.org/atomic"
	"go.uber.org/zap"
	"sync"
)

func init() {
	task := NewDownloadTask()
	task.Register(&botUpdateObserver{})
	registerTask(task)
}


// DownloadTask 下载任务
type DownloadTask struct {
	observerList []UpdateObserver
	isStop       atomic.Bool
}

// NewDownloadTask new DownloadTask
func NewDownloadTask() *DownloadTask {
	return &DownloadTask{
		observerList: []UpdateObserver{},
	}
}

// Name 任务名称
func (t *DownloadTask) Name() string {
	return "DownloadTask"
}

// Register 注册rss更新订阅者
func (t *DownloadTask) Register(observer UpdateObserver) {
	t.observerList = append(t.observerList, observer)
}


// Stop end task
func (t *DownloadTask) Stop() {
	t.isStop.Store(true)
}

// Start run task
func (t *DownloadTask) Start() {
	if config.RunMode == config.TestMode {
		return
	}

	t.isStop.Store(false)
}


// notifyAllObserverUpdate notify all download update observer
func (t *DownloadTask) notifyAllObserverUpdate(subscription *model.Subscribe) {

	wg := sync.WaitGroup{}
	for _, observer := range t.observerList {
		wg.Add(1)
		go func(o UpdateObserver) {
			defer wg.Done()
			o.update(subscription)
		}(observer)
	}
	wg.Wait()
}

// notifyAllObserverErrorUpdate notify all download error update observer
func (t *DownloadTask) notifyAllObserverErrorUpdate(subscription *model.Subscribe) {
	wg := sync.WaitGroup{}
	for _, observer := range t.observerList {
		wg.Add(1)
		go func(o UpdateObserver) {
			defer wg.Done()
			o.errorUpdate(subscription)
		}(observer)
	}
	wg.Wait()
}

// UpdateObserver Update observer
type UpdateObserver interface {
	update(*model.Subscribe)
	errorUpdate(*model.Subscribe)
	id() string
}

type botUpdateObserver struct {
}

func (o *botUpdateObserver) update(subscription *model.Subscribe) {
	zap.S().Debugf("%v receiving [%d]%v update", o.id(), subscription.ID, subscription.Title)
}

func (o *botUpdateObserver) errorUpdate(subscription *model.Subscribe) {
	zap.S().Debugf("%v receiving [%d]%v error update", o.id(), subscription.ID, subscription.Title)
}

func (o *botUpdateObserver) id() string {
	return "botUpdateObserver"
}
