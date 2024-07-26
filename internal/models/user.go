package models

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type User struct {
	ID       uuid.UUID `json:"id" db:"id" validate:"omitempty"`
	Email    string    `json:"email,omitempty" db:"email" validate:"required,omitempty,lte=60,email"`
	Password string    `json:"password,omitempty" db:"password" validate:"required,omitempty,gte=6"`
}

type UserWithToken struct {
	User  *User  `json:"user"`
	Token string `json:"token"`
}

func (u *User) Prepare() error {
	u.Email = strings.ToLower(u.Email)

	u.Password = strings.TrimSpace(u.Password)

	if err := u.HashPassword(); err != nil {
		return err
	}

	return nil
}

func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) SanitizePassword() {
	u.Password = ""
}

func (u *User) ComparePasswords(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}
