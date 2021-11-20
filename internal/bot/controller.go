package bot

import (
	"fmt"
	"github.com/wesleywxie/gogetit/internal/config"
	"github.com/wesleywxie/gogetit/internal/model"
	"go.uber.org/zap"
	tb "gopkg.in/tucnak/telebot.v3"
)

func startCmdCtr(c tb.Context) error {
	user, _ := model.FindOrCreateUserByTelegramID(c.Chat().ID)
	zap.S().Infof("/start user_id: %d telegram_id: %d", user.ID, user.TelegramID)
	return c.Send(fmt.Sprintf("你好，欢迎使用%v。", config.ProjectName))
}

func subCmdCtr(c tb.Context) error {
	channel, index := GetChannelAndMessageIndexFromMessage(c)
	zap.S().Debugf("about to subscribe to channel=%v with index=%v", channel, index)
	channelInfo, err := B.ChatByUsername("-100" + channel)
	if err != nil {
		zap.S().Warnf("failed to subscribe channel=%v, error=%v", channel, err)
		return c.Send("订阅失败，请检查后重新调用' /sub @ChannelID messageIndex or /sub https://t.me/c/{ChannelID}/{messageIndex}' 命令")
	}

	err = registerChannel(c, channelInfo, index)
	if err != nil {
		zap.S().Warnf("failed to subscribe channel=%v, error=%v", channel, err)
		return c.Send("订阅失败，请检查后重新调用' /sub @ChannelID messageIndex or /sub https://t.me/c/{ChannelID}/{messageIndex}' 命令")
	}
	return nil
}

func unsubCmdCtr(c tb.Context) error {
	channel, index := GetChannelAndMessageIndexFromMessage(c)
	zap.S().Debugf("about to unsubscribe from channel=%v with index=%v", channel, index)
	channelInfo, err := B.ChatByUsername("-100" + channel)
	if err != nil {
		zap.S().Warnf("failed to subscribe channel=%v, error=%v", channel, err)
		return c.Send("退订失败，请检查后重新调用' /unsub @ChannelID or /unsub https://t.me/c/{ChannelID}/' 命令")
	}
	err = unregisterChannel(c, channelInfo)
	if err != nil {
		zap.S().Warnf("failed to subscribe channel=%v, error=%v", channel, err)
		return c.Send("退订失败，请检查后重新调用' /unsub @ChannelID or /unsub https://t.me/c/{ChannelID}/' 命令")
	}

	return nil
}

func listCmdCtr(c tb.Context) error {
	userID := c.Chat().ID
	user, err := model.FindOrCreateUserByTelegramID(userID)
	if err != nil {
		zap.S().Warnw("error while retrieving subscriptions",
			"userID", userID,
			"error", err.Error())
		return c.Send(fmt.Sprintf("当前订阅列表为空"))
	}

	subscriptions, err := user.GetSubscriptions()
	if err != nil {
		zap.S().Warnw("error while retrieving subscriptions",
			"userID", userID,
			"error", err.Error())
		return c.Send(fmt.Sprintf("当前订阅列表为空"))
	}

	rspMessage := ""
	if len(subscriptions) == 0 {
		rspMessage = "当前订阅列表为空"
	} else {
		rspMessage = "当前订阅列表：\n"
		for _, sub := range subscriptions {
			rspMessage += fmt.Sprintf("[%d] - %s\n", sub.ID, sub.Title)
		}
	}
	return c.Send(rspMessage)
}

func helpCmdCtr(c tb.Context) error {
	message := `
命令：
/sub @{ChannelID} {MessageIndex} or https://t.me/c/{ChannelID}/{MessageIndex} 
/unsub @{ChannelID} or https://t.me/c/{ChannelID}/ 
/list 列出所有订阅 
/help 帮助
/version 查看当前bot版本
详细使用方法请看：https://github.com/wesleywxie/gogetit
`
	return c.Send(message)
}

func versionCmdCtr(c tb.Context) error {
	return c.Send(config.AppVersionInfo())
}