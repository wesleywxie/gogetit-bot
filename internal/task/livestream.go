package task

import (
	"github.com/wesleywxie/gogetit/internal/cmd"
	"github.com/wesleywxie/gogetit/internal/config"
	"github.com/wesleywxie/gogetit/internal/model"
	"go.uber.org/atomic"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"time"
)

func init() {
	task := NewLivestreamTask()
	task.Register(&telegramBotUpdateObserver{})
	registerTask(task)
}

// Register 注册更新订阅者
func (t *LivestreamUpdateTask) Register(observer LivestreamUpdateObserver) {
	t.observerList = append(t.observerList, observer)
}

// Deregister 注销更新订阅者
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

// Stop 停止
func (t *LivestreamUpdateTask) Stop() {
	t.isStop.Store(true)
}

// Start 启动
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

			subscriptions, err := model.GetSubscriptions()

			if err != nil {
				zap.S().Errorf("Failed to get subscriptions from db, error:%v", err)
			}

			for _, subscription := range subscriptions {
				zap.S().Debugf("Checking subscription[%d], %v", subscription.ID, subscription.KOL)
				// Check if the recoding is ON or OFF
				if subscription.Streaming == false {
					// Check if the broadcast is ON or OFF
					content, _ := makeGetRequest(subscription.Link)
					zap.S().Debugw("Check liveness...",
						"url", subscription.Link,
						"content", content,
					)
					// TODO better mechanism to check liveness
					broadcasting := len(content) > 0

					if broadcasting {
						// record
						record := make(chan string)
						go cmd.Recording(subscription, record)

						// upload
						go cmd.Upload(record)
					}
				}
			}

			time.Sleep(time.Duration(config.UpdateInterval) * time.Minute)
		}
	}()
}

func makeGetRequest(url string) (content string, err error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		zap.S().Error(err)
		return
	}
	req.Header.Set("User-Agent", config.UserAgent)
	resp, err := client.Do(req)
	if err != nil {
		zap.S().Error(err)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		zap.S().Error(err)
		return
	}

	content = string(body)
	return
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
