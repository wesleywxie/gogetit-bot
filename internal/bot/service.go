package bot

import (
	"fmt"
	"github.com/wesleywxie/gogetit-bot/internal/config"
	"github.com/wesleywxie/gogetit-bot/internal/model"
	"go.uber.org/zap"
	tb "gopkg.in/telebot.v3"
	"regexp"
)

// IsUserAllowed check user is allowed to use bot
func isUserAllowed(upd *tb.Update) bool {
	if upd == nil {
		return false
	}

	var userID int64

	if upd.Message != nil {
		userID = int64(upd.Message.Sender.ID)
	} else if upd.Callback != nil {
		userID = int64(upd.Callback.Sender.ID)
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

// GetHyperlinkFromMessage get hyperlink from message mention
func GetHyperlinkFromMessage(m *tb.Message) (url string) {
	for _, entity := range m.Entities {
		if entity.Type == tb.EntityURL {
			if url == "" {
				url = m.Text[entity.Offset : entity.Offset+entity.Length]
			}
		}
	}

	var payloadMatching = relaxUrlMatcher.FindStringSubmatch(m.Payload)
	if url == "" && len(payloadMatching) > 0 && payloadMatching[0] != "" {
		url = payloadMatching[0]
	}

	return
}

func subscribeLiveStream(c tb.Context, url string) (err error) {
	msg, err := B.Send(c.Chat(), "处理中...")
	chatID := c.Chat().ID
	subscription, err := model.SubscribeLiveStream(chatID, url)
	zap.S().Infof("%d subscribe url %s", chatID, url)

	if err == nil {
		_, _ = B.Edit(msg, fmt.Sprintf("[%s](%s) 订阅成功", subscription.KOL, subscription.Link),
			&tb.SendOptions{
				DisableWebPagePreview: true,
				ParseMode:             tb.ModeMarkdown,
			})
	} else {
		_, _ = B.Edit(msg, "订阅失败")
	}
	return
}
