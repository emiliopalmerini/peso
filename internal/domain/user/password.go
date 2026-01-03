package user

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost      = 12
	minPasswordLen  = 8
)

var (
	ErrPasswordTooShort = errors.New("password must be at least 8 characters")
	ErrPasswordMismatch = errors.New("password does not match")
)

type Password struct {
	hash string
}

func NewPassword(plaintext string) (Password, error) {
	if len(plaintext) < minPasswordLen {
		return Password{}, ErrPasswordTooShort
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(plaintext), bcryptCost)
	if err != nil {
		return Password{}, err
	}

	return Password{hash: string(hash)}, nil
}

func NewPasswordFromHash(hash string) Password {
	return Password{hash: hash}
}

func (p Password) Verify(plaintext string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(p.hash), []byte(plaintext))
	return err == nil
}

func (p Password) Hash() string {
	return p.hash
}

func (p Password) IsEmpty() bool {
	return p.hash == ""
}
