package entity

import (
	"time"
)

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

func CreateUser(username, password, email string, status UserStatus) *User {
	return &User{
		Username: username,
		Password: password,
		Email:    email,
		Status:   status,
	}
}

func (u *User) Update(username, password, email string, status UserStatus) *User {
	u.Username = username
	u.Password = password
	u.Email = email
	u.Status = status
	return u
}
