package models

import (
	"time"
)

type Deal struct {
	Id           uint `gorm:"primaryKey"`
	Currencies   string
	LastCost     float64
	UpdatedAtUtc time.Time
}

func NewDeal(currencies string, lastCost float64, updatedAtUtc time.Time) *Deal {
	return &Deal{0, currencies, lastCost, updatedAtUtc}
}
