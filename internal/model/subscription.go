package model

import "github.com/wesleywxie/gogetit/internal/config"

type Subscription struct {
	ID                 uint `gorm:"primary_key;AUTO_INCREMENT"`
	UserID             int64
	Category           string
	KOL                string
	Link 			   string
	Interval           int
	WaitTime           int
	EditTime
}

func SubscribeLiveStream(userID int64, url string) (subscription Subscription, err error) {

	if err := db.Where("user_id=? and link=?", userID, url).Find(&subscription).Error; err != nil {
		if err.Error() == "record not found" {
			subscription.UserID = userID
			subscription.Link = url
			subscription.Interval = config.UpdateInterval
			subscription.WaitTime = config.UpdateInterval
			if db.Create(&subscription).Error == nil {
				return subscription, nil
			}
		}
	}
	return
}
