package bot

import (
	"fmt"

	"github.com/wesleywxie/gogetit-bot/internal/cmd"
	"github.com/wesleywxie/gogetit-bot/internal/config"
	"github.com/wesleywxie/gogetit-bot/internal/model"
	"go.uber.org/zap"
	tb "gopkg.in/telebot.v3"
)

func startCmdCtr(c tb.Context) error {
	user, _ := model.FindOrCreateUserByTelegramID(c.Chat().ID)
	zap.S().Infof("/start user_id: %d telegram_id: %d", user.ID, user.TelegramID)
	return c.Send(fmt.Sprintf("你好，欢迎使用%v。", config.ProjectName))
}

func dlCmdCtr(c tb.Context) (err error) {
	url := GetHyperlinkFromMessage(c.Message())

	zap.S().Debugw("Received ytb download command",
		"url", url,
	)

	gen := make(chan string)
	download := make(chan string)

	msg, err := B.Send(c.Chat(), "正在下载...")

	// generate filename
	go cmd.GetFilename(c, msg, url, gen)

	// execute download and store
	go cmd.ExecDownload(c, msg, url, gen, download)

	// upload with rclone
	go cmd.Sync(c, msg, download)

	return
}

func subCmdCtr(c tb.Context) (err error) {

	url := GetHyperlinkFromMessage(c.Message())

	zap.S().Debugw("Received live stream recording subscription link",
		"url", url,
	)

	if url != "" {
		err = subscribeLiveStream(c, url)
	}

	return
}

func unsubCmdCtr(c tb.Context) (err error) {
	url := GetHyperlinkFromMessage(c.Message())

	if url != "" {
		subscription, err := model.GetSubscriptionsByUserIDAndURL(c.Chat().ID, url)

		if err != nil {
			if err.Error() == "record not found" {
				_, err = B.Send(c.Chat(), "未订阅该源")

			} else {
				_, err = B.Send(c.Chat(), "退订失败")
			}
			return err

		}

		err = subscription.Unsubscribe()
		if err == nil {
			_, _ = B.Send(
				c.Chat(),
				fmt.Sprintf("退订 [%s](%s) 成功", subscription.KOL, subscription.Link),
				&tb.SendOptions{
					DisableWebPagePreview: true,
					ParseMode:             tb.ModeMarkdown,
				},
			)
			zap.S().Infof("%d for  unsubscribe [%s]%s", c.Chat().ID, subscription.KOL, subscription.Link)
		} else {
			_, err = B.Send(c.Chat(), err.Error())
		}
		return err

	}
	_, err = B.Send(c.Chat(), "退订请使用' /unsub URL ' 命令")
	return
}

func listCmdCtr(c tb.Context) (err error) {
	user, err := model.FindOrCreateUserByTelegramID(c.Chat().ID)
	if err != nil {
		_, err = B.Send(c.Chat(), fmt.Sprintf("内部错误 list@1"))
		return
	}

	subscriptions, err := user.GetSubscriptions()
	if err != nil {
		_, err = B.Send(c.Chat(), fmt.Sprintf("内部错误 list@2"))
		return
	}

	rspMessage := "当前订阅列表：\n"
	if len(subscriptions) == 0 {
		rspMessage = "订阅列表为空"
	} else {
		for _, subscription := range subscriptions {
			rspMessage = rspMessage + fmt.Sprintf("[[%d]] [%s](%s)\n", subscription.ID, subscription.KOL, subscription.Link)
		}
	}
	_, err = B.Send(c.Chat(), rspMessage, &tb.SendOptions{
		DisableWebPagePreview: true,
		ParseMode:             tb.ModeMarkdown,
	})
	return
}

func helpCmdCtr(c tb.Context) error {
	message := `
命令： 
/dl 下载 url
/help 帮助
/sub 订阅直播源
/unsub 退订直播源
/list 列出当前订阅的直播源
/version 查看当前bot版本
详细使用方法请看：https://github.com/wesleywxie/gogetit-bot
`
	return c.Send(message)
}

func versionCmdCtr(c tb.Context) error {
	return c.Send(config.AppVersionInfo())
}
