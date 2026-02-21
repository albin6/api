package domain

import (
	"time"
)

type Student struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	FullName      string    `json:"full_name" gorm:"not null"`
	Email         string    `json:"email" gorm:"unique;not null"`
	Phone         string    `json:"phone" gorm:"not null"`
	ProgramStatus bool      `json:"program_status" gorm:"index"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
