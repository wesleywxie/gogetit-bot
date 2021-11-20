package bot

import (
	"fmt"
	"github.com/wesleywxie/gogetit/internal/model"
	"go.uber.org/zap"
	tb "gopkg.in/tucnak/telebot.v3"
	"regexp"
	"strconv"
	"strings"
)

func registerChannel(c tb.Context, channel *tb.Chat, index int) (err error) {
	msg, _ := B.Send(c.Chat(),"处理中...")
	channelID := channel.ID
	messageIndex := uint(index)

	subscribe, err := model.RegisterChannel(channelID, c.Chat().ID, channel.Title, messageIndex)

	if err != nil {
		msg, err = B.Edit(msg, fmt.Sprintf("%s，订阅失败", err))
	}

	zap.S().Infof("%d subscribe [%d] at %d", c.Chat().ID, subscribe.ChannelID, subscribe.MessageIndex)
	msg, err = B.Edit(msg, fmt.Sprintf("频道 %v 订阅成功", channel.Title))
	return err
}

func unregisterChannel(c tb.Context, channel *tb.Chat) (err error) {
	msg, _ := B.Send(c.Chat(),"处理中...")
	channelID := channel.ID
	subscribe, err := model.GetSubscribeByChannel(channelID, c.Chat().ID)

	if subscribe == nil {
		zap.S().Warnf("error when unsubscribing channel %d, err=%v", channelID, err.Error())
		_, err = B.Edit(msg, "未订阅该频道")
	} else {
		err = model.UnregisterChannel(channelID, c.Chat().ID)
		if err == nil {
			_, err = B.Edit(
				msg,
				fmt.Sprintf("频道 %v 退订成功！", channel.Title),
			)
			zap.S().Infof("%d unsubscribe [%d]%s", c.Chat().ID, channelID, channel.Title)
		} else {
			_, err = B.Edit(msg, err.Error())
		}
	}
	return err
}

var relaxUrlMatcher = regexp.MustCompile(`^(https?://.*?)($| )`)

// GetChannelAndMessageIndexFromMessage get URL and mention from message
func GetChannelAndMessageIndexFromMessage(c tb.Context) (channel string, index int) {
	message := c.Message()
	channel = ""
	index = 0
	payloads :=  strings.Split(message.Payload, " ")

	switch length := len(payloads); length {
	case 2:
		// e.g.
		// For public channel: @TestFlightCN 11073
		// For private Channel: @1304836281 8498
		channel = payloads[0][1:]
		i, err := strconv.Atoi(payloads[1])
		if err != nil {
			zap.S().Error(err)
		} else {
			index = i
		}
	case 1:
		// Check the payload is URL type
		// e.g.
		// For public channel: https://t.me/TestFlightCN/11073
		// For private Channel: https://t.me/c/1304836281/8498
		payloadMatching := relaxUrlMatcher.FindStringSubmatch(payloads[0])
		if len(payloadMatching) > 0 && payloadMatching[0] != "" {
			payloads = strings.Split(payloadMatching[0], "/")
			channel = payloads[4]
			i, err := strconv.Atoi(payloads[5])
			if err != nil {
				zap.S().Error(err)
			} else {
				index = i
			}
		}
	default:
		zap.S().Error("Received 0 payload")
	}

	return
}