package domain

import (
	"time"

	"github.com/ryuyb/fusion/internal/pkg/errors"
	"github.com/ryuyb/fusion/internal/pkg/util"
)

type User struct {
	ID        int64
	Username  string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func CreateUser(username, email, password string) (*User, error) {
	hashPassword, err := util.HashPassword(password)
	if err != nil {
		return nil, errors.Internal(err)
	}
	return &User{
		Username: username,
		Email:    email,
		Password: hashPassword,
	}, nil
}

func (u *User) Update(username, email, password string) (*User, error) {
	hashPassword, err := util.HashPassword(password)
	if err != nil {
		return nil, errors.Internal(err)
	}
	u.Username = username
	u.Email = email
	u.Password = hashPassword
	return u, nil
}
