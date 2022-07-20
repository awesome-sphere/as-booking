package models

import (
	"github.com/shopspring/decimal"
)

type SeatTier string

const (
	Standard SeatTier = "STANDARD"
	Premium  SeatTier = "PREMIUM"
)

type SeatType struct {
	ID    int64      `json:"id" gorm:"primaryKey;autoincrement;"`
	Price decimal.Decimal `json:"price"`
	Type  SeatTier        `json:"seat_type" sql:"type:seat_type"`
}
