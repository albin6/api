package domain

import (
	"time"
)

type FollowUpReminder struct {
	ID         uint            `json:"id" gorm:"primaryKey"`
	FollowUpID uint            `json:"follow_up_id" gorm:"not null;index"`
	FollowUp   StudentFollowUp `json:"follow_up" gorm:"foreignKey:FollowUpID"`
	RemindAt   time.Time       `json:"remind_at" gorm:"not null;index"`
	Sent       bool            `json:"sent" gorm:"default:false;index"`
	SentAt     *time.Time      `json:"sent_at"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
}
