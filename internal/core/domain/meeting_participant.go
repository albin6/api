package domain

import (
	"time"
)

type MeetingParticipant struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	MeetingID uint      `json:"meeting_id" gorm:"not null;index"`
	Meeting   Meeting   `json:"-" gorm:"foreignKey:MeetingID"`
	UserID    string    `json:"user_id" gorm:"type:varchar(36);not null;index"`
	User      User      `json:"user" gorm:"foreignKey:UserID"`
	CreatedAt time.Time `json:"created_at"`
}
