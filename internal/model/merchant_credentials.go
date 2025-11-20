package model

import "time"

type MerchantCredentials struct {
	Id         string    `gorm:"type:varchar(100);primary_key"`
	MerchantId string    `gorm:"type:varchar(100);not null"`
	ApiKey     string    `gorm:"type:varchar(100);not null"`
	CreatedAt  time.Time `gorm:"type:timestamp;not null"`
}
