package models

import (
	"errors"
	"fmt"
	"time"
	"unicode"
)

const (
	MinUsernameLength = 6
	MinPasswordLength = 8
	MaxUsernameLength = 15
	MaxPasswordLength = 12
)

type User struct {
	Id         int64     `json:"id"`
	Username   string    `json:"username"`
	Password   string    `json:"-"`
	CreateTime time.Time `json:"create_time"`
	IsAdmin    bool      `json:"is_admin"`
}

type UserForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"is_admin"`
}

type PayloadKey string

const PayloadContextKey PayloadKey = "payload"

type JwtPayload struct {
	Username string
	IsAdmin  bool
}

func isEnglishLetter(c rune) bool {
	return ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z')
}

func (form *UserForm) Validate() error {
	runedUsername := []rune(form.Username)
	runedPassword := []rune(form.Password)
	if len(runedUsername) < MinUsernameLength || len(runedUsername) > MaxUsernameLength {
		return fmt.Errorf("username length must be from %d to %d characters", MinUsernameLength, MaxUsernameLength)
	}
	if len(runedPassword) < MinPasswordLength || len(runedPassword) > MaxPasswordLength {
		return fmt.Errorf("password length must be from %d to %d characters", MinPasswordLength, MaxPasswordLength)
	}

	for _, sym := range runedUsername {
		if !unicode.IsDigit(sym) && !isEnglishLetter(sym) {
			return errors.New("username can only include symbols: A-Z, a-z, 0-9")
		}
	}
	return nil

}
