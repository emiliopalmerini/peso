package user

import (
	"errors"
	"strings"
	"time"
)

type User struct {
	id           UserID
	name         string
	email        string
	passwordHash string
	active       bool
	createdAt    time.Time
	updatedAt    time.Time
}

var (
	ErrEmptyName = errors.New("user name cannot be empty")
)

func NewUser(id, name, email string) (*User, error) {
	userID, err := NewUserID(id)
	if err != nil {
		return nil, err
	}

	trimmedName := strings.TrimSpace(name)
	if trimmedName == "" {
		return nil, ErrEmptyName
	}

	now := time.Now()

	return &User{
		id:           userID,
		name:         trimmedName,
		email:        email,
		passwordHash: "",
		active:       true,
		createdAt:    now,
		updatedAt:    now,
	}, nil
}

func NewUserWithPassword(id, name, email, password string) (*User, error) {
	u, err := NewUser(id, name, email)
	if err != nil {
		return nil, err
	}

	if err := u.SetPassword(password); err != nil {
		return nil, err
	}

	return u, nil
}

func (u *User) ID() UserID {
	return u.id
}

func (u *User) Name() string {
	return u.name
}

func (u *User) Email() string {
	return u.email
}

func (u *User) IsActive() bool {
	return u.active
}

func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

func (u *User) UpdatedAt() time.Time {
	return u.updatedAt
}

func (u *User) Deactivate() {
	u.active = false
	u.updatedAt = time.Now()
}

func (u *User) Activate() {
	u.active = true
	u.updatedAt = time.Now()
}

func (u *User) UpdateEmail(email string) {
	u.email = email
	u.updatedAt = time.Now()
}

func (u *User) UpdateName(name string) error {
	trimmedName := strings.TrimSpace(name)
	if trimmedName == "" {
		return ErrEmptyName
	}

	u.name = trimmedName
	u.updatedAt = time.Now()
	return nil
}

func (u *User) PasswordHash() string {
	return u.passwordHash
}

func (u *User) HasPassword() bool {
	return u.passwordHash != ""
}

func (u *User) SetPassword(plaintext string) error {
	pwd, err := NewPassword(plaintext)
	if err != nil {
		return err
	}
	u.passwordHash = pwd.Hash()
	u.updatedAt = time.Now()
	return nil
}

func (u *User) SetPasswordHash(hash string) {
	u.passwordHash = hash
}

func (u *User) VerifyPassword(plaintext string) bool {
	if u.passwordHash == "" {
		return false
	}
	pwd := NewPasswordFromHash(u.passwordHash)
	return pwd.Verify(plaintext)
}
