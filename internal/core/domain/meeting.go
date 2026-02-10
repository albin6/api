package domain

import (
	"time"
)

type MeetingStatus string

const (
	MeetingScheduled MeetingStatus = "SCHEDULED"
	MeetingCompleted MeetingStatus = "COMPLETED"
	MeetingCancelled MeetingStatus = "CANCELLED"
)

func (m MeetingStatus) String() string {
	return string(m)
}

type Meeting struct {
	ID          uint            `json:"id" gorm:"primaryKey"`
	FollowUpID  uint            `json:"follow_up_id" gorm:"not null;index"`
	FollowUp    StudentFollowUp `json:"follow_up" gorm:"foreignKey:FollowUpID"`
	ScheduledAt time.Time       `json:"scheduled_at" gorm:"not null;index"`
	MeetingLink string          `json:"meeting_link" gorm:"not null"`
	Status      MeetingStatus   `json:"status" gorm:"type:varchar(20);default:'SCHEDULED'"`
	CreatedBy   uint            `json:"created_by" gorm:"not null"`
	Creator     User            `json:"creator" gorm:"foreignKey:CreatedBy"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}
