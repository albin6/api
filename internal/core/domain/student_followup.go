package domain

import (
	"time"
)

type StudentFollowUp struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	StudentID string `json:"student_id" gorm:"type:varchar(36);not null;index"`
	// Student      Student       `json:"student" gorm:"-"` // Removed local relation
	AssignedTo string `json:"assigned_to" gorm:"type:varchar(36);not null;index"`
	// AssignedUser User          `json:"assigned_user" gorm:"-"` // Removed local relation
	Stage     FollowUpStage `json:"stage" gorm:"type:varchar(30);not null;index"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}
