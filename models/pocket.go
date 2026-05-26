package models

import (
	"time"

	"gorm.io/gorm"
)

type Pocket struct {
	PocketID   uint           `gorm:"primaryKey;column:pocket_id" json:"pocket_id"`
	PocketName string         `gorm:"column:pocket_name" json:"pocket_name"` // VARCHAR(100)
	Amount     float64        `gorm:"column:amount" json:"amount"`           // DECIMAL(10,2)
	Currency   string         `gorm:"column:currency" json:"currency"`       // VARCHAR(3), DEFAULT 'USD'
	AccountID  uint           `gorm:"column:account_id" json:"account_id"`
	CreatedAt  time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

type SavePocket struct {
	PocketName string   `gorm:"column:pocket_name" json:"pocket_name"` // VARCHAR(100)
	Amount     *float64 `gorm:"column:amount" json:"amount"`           // DECIMAL(10,2)
	Currency   string   `gorm:"column:currency" json:"currency"`       // VARCHAR(3), DEFAULT 'USD'
	AccountID  uint     `gorm:"column:account_id" json:"account_id"`
}

func (Pocket) TableName() string {
	return "pockets"
}

func (SavePocket) TableName() string {
	return "pockets"
}
