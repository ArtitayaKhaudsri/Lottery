package models

import "time"

type Draw struct {
	DrawID   uint64    `gorm:"primaryKey"`
	DrawDate time.Time `gorm:"not null"`
}

func (Draw) TableName() string {
	return "draws"
}

type Prize struct {
	PrizeID       uint64 `gorm:"primaryKey;column:prize_id"`
	DrawID        uint64 `gorm:"column:draw_id;not null"`
	PrizeType     string `gorm:"column:prize_type"`
	WinningNumber string `gorm:"column:winning_number"`

	Draw *Draw `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}

func (Prize) TableName() string {
	return "prizes"
}

type RedisPrize struct {
	Number    string `json:"number"`
	PrizeType string `json:"prizeType"`
}

type PrizeResponse struct {
	Date   time.Time `json:"date"`
	Number string    `json:"number"`
	Type   string    `json:"type"`
}
