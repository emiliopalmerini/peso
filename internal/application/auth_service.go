package application

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"peso/internal/domain/session"
	"peso/internal/domain/user"
	"peso/internal/interfaces"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailAlreadyExists = errors.New("email already registered")
	ErrAuthUserNotFound   = errors.New("user not found")
	ErrSessionExpired     = errors.New("session expired")
	ErrNoPassword         = errors.New("user has no password set")
	ErrInvalidEmail       = errors.New("invalid email format")
)

type AuthService struct {
	userRepo    interfaces.UserRepository
	sessionRepo interfaces.SessionRepository
}

func NewAuthService(userRepo interfaces.UserRepository, sessionRepo interfaces.SessionRepository) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
	}
}

func (s *AuthService) Register(name, email, password string) (*user.User, *session.Session, error) {
	email = strings.TrimSpace(strings.ToLower(email))

	if !isValidEmail(email) {
		return nil, nil, ErrInvalidEmail
	}

	exists, err := s.userRepo.EmailExists(email)
	if err != nil {
		return nil, nil, err
	}
	if exists {
		return nil, nil, ErrEmailAlreadyExists
	}

	userID := uuid.New().String()
	u, err := user.NewUserWithPassword(userID, name, email, password)
	if err != nil {
		return nil, nil, err
	}

	if err := s.userRepo.Save(u); err != nil {
		return nil, nil, err
	}

	sess, err := session.NewSession(u.ID())
	if err != nil {
		return nil, nil, err
	}

	if err := s.sessionRepo.Save(sess); err != nil {
		return nil, nil, err
	}

	return u, sess, nil
}

func (s *AuthService) Login(email, password string) (*user.User, *session.Session, error) {
	email = strings.TrimSpace(strings.ToLower(email))

	u, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, nil, ErrInvalidCredentials
	}

	if !u.HasPassword() {
		return u, nil, ErrNoPassword
	}

	if !u.VerifyPassword(password) {
		return nil, nil, ErrInvalidCredentials
	}

	sess, err := session.NewSession(u.ID())
	if err != nil {
		return nil, nil, err
	}

	if err := s.sessionRepo.Save(sess); err != nil {
		return nil, nil, err
	}

	return u, sess, nil
}

func (s *AuthService) SetPassword(email, password string) (*user.User, *session.Session, error) {
	email = strings.TrimSpace(strings.ToLower(email))

	u, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, nil, ErrAuthUserNotFound
	}

	if err := u.SetPassword(password); err != nil {
		return nil, nil, err
	}

	if err := s.userRepo.Save(u); err != nil {
		return nil, nil, err
	}

	sess, err := session.NewSession(u.ID())
	if err != nil {
		return nil, nil, err
	}

	if err := s.sessionRepo.Save(sess); err != nil {
		return nil, nil, err
	}

	return u, sess, nil
}

func (s *AuthService) Logout(token string) error {
	return s.sessionRepo.DeleteByToken(token)
}

func (s *AuthService) ValidateSession(token string) (*user.User, error) {
	sess, err := s.sessionRepo.FindByToken(token)
	if err != nil {
		return nil, ErrSessionExpired
	}

	if sess.IsExpired() {
		s.sessionRepo.DeleteByToken(token)
		return nil, ErrSessionExpired
	}

	u, err := s.userRepo.FindByID(sess.UserID())
	if err != nil {
		return nil, ErrAuthUserNotFound
	}

	return u, nil
}

func (s *AuthService) CleanupExpiredSessions() error {
	return s.sessionRepo.DeleteExpired()
}

func isValidEmail(email string) bool {
	if email == "" {
		return false
	}
	atIndex := strings.Index(email, "@")
	if atIndex < 1 {
		return false
	}
	dotIndex := strings.LastIndex(email, ".")
	if dotIndex < atIndex+2 || dotIndex >= len(email)-1 {
		return false
	}
	return true
}
