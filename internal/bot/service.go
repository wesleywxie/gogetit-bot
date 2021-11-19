package bot

import (
	"fmt"
	"github.com/wesleywxie/gogetit/internal/config"
	"github.com/wesleywxie/gogetit/internal/model"
	"go.uber.org/zap"
	tb "gopkg.in/tucnak/telebot.v3"
	"regexp"
	"strconv"
	"strings"
)


// IsUserAllowed check user is allowed to use bot
func isUserAllowed(upd *tb.Update) bool {
	if upd == nil {
		return false
	}

	var userID int64

	if upd.Message != nil {
		userID = upd.Message.Sender.ID
	} else if upd.Callback != nil {
		userID = upd.Callback.Sender.ID
	} else {
		return false
	}

	if len(config.AllowUsers) == 0 {
		return true
	}

	for _, allowUserID := range config.AllowUsers {
		if allowUserID == userID {
			return true
		}
	}

	zap.S().Infow("user not allowed", "userID", userID)
	return false
}


func registerChannel(c tb.Context, channel *tb.Chat, index int) (err error) {
	msg, _ := B.Send(c.Chat(),"处理中...")
	channelID := channel.ID
	messageIndex := uint(index)

	subscribe, err := model.RegisterChannel(channelID, c.Chat().ID, messageIndex)

	if err != nil {
		msg, err = B.Edit(msg, fmt.Sprintf("%s，订阅失败", err))
	}

	zap.S().Infof("%d subscribe [%d] at %d", c.Chat().ID, subscribe.ChannelID, subscribe.MessageIndex)


	msg, err = B.Edit(msg, fmt.Sprintf("频道 %v 订阅成功", channel.Title))

	return err
}


// CheckAdmin check user is admin of group/channel
func CheckAdmin(upd *tb.Update) bool {

	if upd.Message != nil {
		if HasAdminType(upd.Message.Chat.Type) {
			adminList, _ := B.AdminsOf(upd.Message.Chat)
			for _, admin := range adminList {
				if admin.User.ID == upd.Message.Sender.ID {
					return true
				}
			}

			return false
		}

		return true
	} else if upd.Callback != nil {
		if HasAdminType(upd.Callback.Message.Chat.Type) {
			adminList, _ := B.AdminsOf(upd.Callback.Message.Chat)
			for _, admin := range adminList {
				if admin.User.ID == upd.Callback.Sender.ID {
					return true
				}
			}

			return false
		}

		return true
	}
	return false
}

// HasAdminType check if the message is sent in the group/channel environment
func HasAdminType(t tb.ChatType) bool {
	hasAdmin := []tb.ChatType{tb.ChatGroup, tb.ChatSuperGroup, tb.ChatChannel, tb.ChatChannelPrivate}
	for _, n := range hasAdmin {
		if t == n {
			return true
		}
	}
	return false
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