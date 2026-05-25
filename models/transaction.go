package models

import (
	"time"
)

// Transaction represents the transaction model
type Transaction struct {
	TxnID        uint      `gorm:"primaryKey;column:txn_id" json:"txn_id"`
	FromPocketID uint      `gorm:"column:from_pocket_id" json:"from_pocket_id"`
	ToPocketID   uint      `gorm:"column:to_pocket_id" json:"to_pocket_id"`
	Amount       float64   `gorm:"column:amount" json:"amount"` // DECIMAL(10,2)
	TxnDate      time.Time `gorm:"column:txn_date;autoCreateTime" json:"txn_date"`
}

// TableName specifies the table name for Transaction model
func (Transaction) TableName() string {
	return "transactions"
}
