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

func listenCmdCtr(c tb.Context) error {
	channel, index := GetChannelAndMessageIndexFromMessage(c)
	zap.S().Debugf("Received channel: %v and index: %v", channel, index)
	if channel != "" {
		channelInfo, err := B.ChatByUsername("-100" + channel)
		if err != nil {
			zap.S().Warnf("failed to subscribe channel=%v, error=%v", channel, err)
			return c.Send("订阅失败，请检查后重新调用' /listen @ChannelID messageIndex or /listen https://t.me/c/{ChannelID}/{messageIndex}' 命令")
		}

		registerChannel(c, channelInfo, index)
	} else {
		return c.Send("频道订阅请使用' /listen @ChannelID messageIndex or /listen https://t.me/c/{ChannelID}/{messageIndex}' 命令")
	}

	return nil
}

func helpCmdCtr(c tb.Context) error {
	message := `
命令：
/listen @{ChannelID} {MessageIndex} or https://t.me/c/{ChannelID}/{MessageIndex} 
/help 帮助
/version 查看当前bot版本
详细使用方法请看：https://github.com/wesleywxie/gogetit
`
	return c.Send(message)
}

func versionCmdCtr(c tb.Context) error {
	return c.Send(config.AppVersionInfo())
}
