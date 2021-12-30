package model

import (
	"errors"
	"github.com/wesleywxie/gogetit/internal/config"
	"strings"
)

type Subscription struct {
	ID        uint `gorm:"primary_key;AUTO_INCREMENT"`
	UserID    int64
	Category  string
	KOL       string
	Link      string
	Interval  int
	Streaming bool
	EditTime
}

func (s *Subscription) Unsubscribe() error {
	if s.ID == 0 {
		return errors.New("can't delete 0 subscribe")
	}

	return db.Delete(&s).Error
}

func SubscribeLiveStream(userID int64, url string) (subscription Subscription, err error) {
	if err := db.Where("user_id=? and link=?", userID, url).Find(&subscription).Error; err != nil {
		if err.Error() == "record not found" {
			KOL, category := processUrl(url)
			subscription.UserID = userID
			subscription.KOL = KOL
			subscription.Category = category
			subscription.Link = url
			subscription.Interval = config.UpdateInterval
			subscription.Streaming = false
			if db.Create(&subscription).Error == nil {
				return subscription, nil
			}
		}
	}
	return
}

func GetSubscriptions() ([]Subscription, error) {
	var subscriptions []Subscription

	err := db.Find(&subscriptions).Error

	return subscriptions, err
}

func GetSubscriptionsByID(subscriptionID uint) (Subscription, error) {
	var subscription Subscription

	err := db.Where("id=?", subscriptionID).Find(&subscription).Error

	return subscription, err
}

func GetSubscriptionsByUserID(userID int64) ([]Subscription, error) {
	var subscriptions []Subscription

	err := db.Where("user_id=?", userID).First(&subscriptions).Error

	return subscriptions, err
}

func GetSubscriptionsByUserIDAndURL(userID int64, url string) (Subscription, error) {
	var subscription Subscription

	err := db.Where("user_id=? and link=?", userID, url).First(&subscription).Error

	return subscription, err
}

// UpdateStreamingStatus update streaming status for a specific subscription
func UpdateStreamingStatus(subscriptionID uint, streaming bool) (subscription Subscription, err error) {
	subscription, err = GetSubscriptionsByID(subscriptionID)
	subscription.Streaming = streaming
	db.Save(&subscription)
	return
}

func processUrl(url string) (KOL, category string) {
	if strings.Index(url, "chaturbate") > 0 {
		url = strings.TrimSuffix(url, "/")
		category = "chaturbate"
		KOL = url[strings.LastIndex(url, "/")+1:]
	}
	return
}
