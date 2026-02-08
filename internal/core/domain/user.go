package domain

import (
	"time"
)

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"not null"`
	Email     string    `json:"email" gorm:"unique;not null"`
	Phone     string    `json:"phone" gorm:"not null"`
	Password  string    `json:"password" gorm:"not null"`
	Role      Role      `json:"role" gorm:"type:varchar(20);default:'MEMBER'"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
