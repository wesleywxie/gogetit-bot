package model

import (
	"github.com/jinzhu/gorm"
)

type Subscribe struct {
	ID                 uint `gorm:"primary_key;AUTO_INCREMENT"`
	UserID             int64
	ChannelID          int64
	EnableNotification int
	MessageIndex       uint
	EditTime
}

func RegisterChannel(channelID int64, userID int64, messageIndex uint) (*Subscribe, error) {
	var subscribe Subscribe

	if err := db.Where("channel_id=? and user_id=?", channelID, userID).Find(&subscribe).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			subscribe.ChannelID = channelID
			subscribe.UserID = userID
			subscribe.MessageIndex = messageIndex
			subscribe.EnableNotification = 1
			db.Create(&subscribe)
			return &subscribe, nil
		}
		return nil, err
	}

	return &subscribe, nil
}
