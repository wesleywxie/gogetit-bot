package bot

import (
	"github.com/wesleywxie/gogetit/internal/config"
	"github.com/wesleywxie/gogetit/internal/util"
	"go.uber.org/zap"
	tb "gopkg.in/tucnak/telebot.v3"
	"time"
)

var (
	B *tb.Bot
)

func init() {
	if config.RunMode == config.TestMode {
		return
	}

	poller := &tb.LongPoller{Timeout: 10 * time.Second}

	spamProtected := tb.NewMiddlewarePoller(poller, func(upd *tb.Update) bool {
		if !isUserAllowed(upd) {
			// 检查用户是否可以使用bot
			return false
		}

		if !CheckAdmin(upd) {
			return false
		}
		return true
	})

	zap.S().Infow("init telegram bot",
		"token", config.BotToken,
		"endpoint", config.TelegramEndpoint,
	)

	// create bot
	var err error

	B, err = tb.NewBot(tb.Settings{
		URL:    config.TelegramEndpoint,
		Token:  config.BotToken,
		Poller: spamProtected,
		Client: util.HttpClient,
	})

	if err != nil {
		zap.S().Fatal(err)
		return
	}
}

//Start bot
func Start() {
	if config.RunMode != config.TestMode {
		zap.S().Infof("bot start %s", config.AppVersionInfo())
		B.Start()
	}
}