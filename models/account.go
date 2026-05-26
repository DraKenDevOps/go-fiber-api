package models

import (
	"time"

	"gorm.io/gorm"
)

type Account struct {
	AccountID uint           `gorm:"primaryKey;column:account_id" json:"account_id"`
	Balance   float64        `gorm:"column:balance" json:"balance"` // DECIMAL(10,2)
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

type AccReqBody struct {
	Balance *float64 `gorm:"column:balance" json:"balance"`
}

func (Account) TableName() string {
	return "accounts"
}

func (AccReqBody) TableName() string {
	return "accounts"
}
