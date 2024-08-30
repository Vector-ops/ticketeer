package models

import (
	"time"

	"gorm.io/gorm"
)

type UserRole string

const (
	ADMIN    UserRole = "admin"
	HOST     UserRole = "host"
	ATTENDEE UserRole = "attendee"
)

type User struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	Email     string    `json:"email" grom:"text;not null"`
	Role      UserRole  `json:"role" gorm:"text;default:attendee"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) AfterCreate(db *gorm.DB) (err error) {
	if u.ID == 1 {
		db.Model(u).Update("role", HOST)
	}
	return
}
