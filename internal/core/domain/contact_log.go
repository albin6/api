package domain

import (
	"time"
)

type ContactLog struct {
	ID          uint            `json:"id" gorm:"primaryKey"`
	FollowUpID  uint            `json:"follow_up_id" gorm:"not null;index"`
	FollowUp    StudentFollowUp `json:"-" gorm:"foreignKey:FollowUpID"`
	Successful  bool            `json:"successful" gorm:"not null"`
	Notes       string          `json:"notes" gorm:"type:text"`
	ContactedAt time.Time       `json:"contacted_at" gorm:"not null"`
	CreatedAt   time.Time       `json:"created_at"`
}
