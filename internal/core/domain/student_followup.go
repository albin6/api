package domain

import (
	"time"
)

type StudentFollowUp struct {
	ID           uint          `json:"id" gorm:"primaryKey"`
	StudentID    uint          `json:"student_id" gorm:"not null;index"`
	Student      Student       `json:"student" gorm:"foreignKey:StudentID"`
	AssignedTo   uint          `json:"assigned_to" gorm:"not null;index"`
	AssignedUser User          `json:"assigned_user" gorm:"foreignKey:AssignedTo"`
	Stage        FollowUpStage `json:"stage" gorm:"type:varchar(30);not null;index"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
}
