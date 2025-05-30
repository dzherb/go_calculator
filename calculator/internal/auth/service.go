package auth

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dzherb/go_calculator/calculator/internal/repository"
	"github.com/dzherb/go_calculator/calculator/pkg/security"
	"github.com/jackc/pgx/v5"
)

type Service interface {
	Login(username, password string) (AccessPayload, error)
	Register(username, password string) (AccessPayload, error)
}

type AccessPayload struct {
	Token string    `json:"access_token"`
	User  repo.User `json:"user"`
}

func NewService() Service {
	return &ServiceImpl{
		userVal: DefaultUsernameValidator,
		passVal: DefaultPasswordValidator,
	}
}

type ServiceImpl struct {
	userVal Validator
	passVal Validator
}

func (s *ServiceImpl) Login(username, password string) (AccessPayload, error) {
	const errFmt = "failed to login: %w"

	// Validate credentials early to avoid unnecessary DB calls
	err := s.validateCredentials(username, password)
	if err != nil {
		return AccessPayload{}, fmt.Errorf(errFmt, err)
	}

	ur := repo.NewUserRepository()

	user, err := ur.GetByCredentials(username, password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = fmt.Errorf("invalid credentials")
		}

		return AccessPayload{}, fmt.Errorf(errFmt, err)
	}

	return s.issueToken(user, "failed to login: %w")
}

func (s *ServiceImpl) Register(
	username,
	password string,
) (AccessPayload, error) {
	const errFmt = "failed to register: %w"

	err := s.validateCredentials(username, password)
	if err != nil {
		return AccessPayload{}, err
	}

	ur := repo.NewUserRepository()

	user, err := ur.Create(repo.User{
		Username: username,
		Password: password,
	})

	if err != nil {
		if strings.Contains(err.Error(), "violates unique constraint") {
			err = fmt.Errorf("username is already taken")
		}

		return AccessPayload{}, fmt.Errorf(errFmt, err)
	}

	return s.issueToken(user, errFmt)
}

func (s *ServiceImpl) validateCredentials(username, password string) error {
	if err := s.userVal.Validate(username); err != nil {
		return fmt.Errorf("invalid username: %w", err)
	}

	if err := s.passVal.Validate(password); err != nil {
		return fmt.Errorf("invalid password: %w", err)
	}

	return nil
}

func (s *ServiceImpl) issueToken(
	user repo.User,
	context string,
) (AccessPayload, error) {
	token, err := security.IssueAccessToken(user.ID)
	if err != nil {
		return AccessPayload{}, fmt.Errorf(context, err)
	}

	return AccessPayload{Token: token, User: user}, nil
}
