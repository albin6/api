package domain

import (
	"time"
)

type OutcomeStatus string

const (
	OutcomeSelected OutcomeStatus = "SELECTED"
	OutcomeRejected OutcomeStatus = "REJECTED"
)

func (o OutcomeStatus) String() string {
	return string(o)
}

type MeetingOutcome struct {
	ID             uint          `json:"id" gorm:"primaryKey"`
	MeetingID      uint          `json:"meeting_id" gorm:"not null;unique;index"`
	Meeting        Meeting       `json:"meeting" gorm:"foreignKey:MeetingID"`
	Status         OutcomeStatus `json:"status" gorm:"type:varchar(20);not null"`
	Remarks        string        `json:"remarks" gorm:"type:text;not null"`
	RecordingURL   string        `json:"recording_url" gorm:"not null"`
	NextFollowUpAt *time.Time    `json:"next_follow_up_at"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
}
