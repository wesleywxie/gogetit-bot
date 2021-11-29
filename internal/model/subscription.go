package model

type Subscription struct {
	ID                 uint `gorm:"primary_key;AUTO_INCREMENT"`
	UserID             int64
	Category           string
	KOL                string
	Interval           int
	WaitTime           int
	EditTime
}