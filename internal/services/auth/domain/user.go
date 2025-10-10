package domain

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

var (
	ErrInvalidRole = errors.New("invalid role")
	ErrEmptyFields = errors.New("username/email and password cannot be empty")
	ErrEmailFormat = errors.New("invalid email format")
)

type User struct {
	ID           string
	Username     string
	Email        string
	PasswordHash []byte
	Role         Role
}

func NewUser(username, email, plainPassword string, role Role) (User, error) {
	switch {
	case username == "" || email == "" || plainPassword == "":
		return User{}, ErrEmptyFields
	case role != RoleAdmin && role != RoleUser:
		return User{}, ErrInvalidRole
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}
	return User{
		Username:     username,
		Email:        email,
		PasswordHash: hash,
		Role:         role,
	}, nil
}

func (u *User) VerifyPassword(passsword string) bool {
	return bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(passsword)) == nil
}
