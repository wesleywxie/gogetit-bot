package bot

import (
	"github.com/wesleywxie/gogetit-bot/internal/config"
	"github.com/wesleywxie/gogetit-bot/internal/util"
	"go.uber.org/zap"
	tb "gopkg.in/telebot.v3"
	"log"
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
	log.Printf("init telegram bot, token=%v, endpoint=%v", config.BotToken, config.TelegramEndpoint)

	// create bot
	var err error

	B, err = tb.NewBot(tb.Settings{
		URL:    config.TelegramEndpoint,
		Token:  config.BotToken,
		Poller: spamProtected,
		Client: util.HttpClient,
	})

	if err != nil {
		log.Fatal(err)
		return
	}
}

//Start bot
func Start() {
	if config.RunMode != config.TestMode {
		zap.S().Infof("bot start %s", config.AppVersionInfo())
		setCommands()
		setHandle()
		B.Start()
	}
}

func setCommands() {
	// 设置bot命令提示信息
	commands := []tb.Command{
		{Text: "dl", Description: "下载"},
		{Text: "sub", Description: "订阅直播源"},
		{Text: "unsub", Description: "退订直播源"},
		{Text: "list", Description: "列出当前订阅的直播源"},
		{Text: "help", Description: "使用帮助"},
		{Text: "version", Description: "bot版本"},
	}

	zap.S().Debugf("set bot command %+v", commands)

	if err := B.SetCommands(commands); err != nil {
		zap.S().Errorw("set bot commands failed", "error", err.Error())
	}
}

func setHandle() {
	B.Handle("/start", startCmdCtr)
	B.Handle("/dl", dlCmdCtr)
	B.Handle("/sub", subCmdCtr)
	B.Handle("/unsub", unsubCmdCtr)
	B.Handle("/list", listCmdCtr)
	B.Handle("/help", helpCmdCtr)
	B.Handle("/version", versionCmdCtr)
}
