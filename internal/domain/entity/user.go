package entity

import "time"

type User struct {
	ID        int64
	Username  string
	Password  string
	Email     string
	Status    UserStatus
	CreatedAt time.Time
	UpdatedAt time.Time
	DeleteAt  time.Time
}

type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
	UserStatusBanned   UserStatus = "banned"
)
