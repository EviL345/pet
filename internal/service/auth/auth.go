package auth

import (
	"blog/internal/config"
	"blog/internal/models"
	"blog/internal/repository"
	"blog/pkg/httperrors"
	"blog/pkg/utils"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"
)

const cacheDuration = 3600

type AuthService struct {
	authPgRepo    repository.PgAuth
	authRedisRepo repository.RedisAuth
	cfg           *config.Config
}

func NewService(authPgRepo repository.PgAuth, authRedisRepo repository.RedisAuth, cfg *config.Config) *AuthService {
	return &AuthService{authPgRepo: authPgRepo, authRedisRepo: authRedisRepo, cfg: cfg}
}

func (s *AuthService) Register(user *models.User) (*models.UserWithToken, error) {
	isExists, err := s.authPgRepo.GetUserByEmail(user.Email)
	if err != nil {
		return nil, err
	}

	if isExists != nil {
		return nil, httperrors.NewRestError(http.StatusBadRequest, "User with given email already exists", nil)
	}

	if err = user.Prepare(); err != nil {
		return nil, httperrors.NewRestError(http.StatusBadRequest, httperrors.BadRequest.Error(), fmt.Errorf("AuthService.Register.Prepare: %w", err))
	}

	userId := uuid.New()
	user.ID = userId

	createdUser, err := s.authPgRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	createdUser.SanitizePassword()

	token, err := utils.GenerateJWTToken(createdUser, s.cfg)
	if err != nil {
		return nil, fmt.Errorf("AuthService.Register.GenerateJWTToken: %w", err)
	}

	return &models.UserWithToken{User: createdUser, Token: token}, nil
}

func (s *AuthService) Login(user *models.User) (*models.UserWithToken, error) {
	foundedUser, err := s.authPgRepo.GetUserByEmail(user.Email)
	if err != nil {
		return nil, fmt.Errorf("AuthService.Login.GetUserByEmail: %w", err)
	}

	if foundedUser == nil {
		return nil, httperrors.NewRestError(http.StatusBadRequest, "User with given email does not exists", nil)
	}

	if err = foundedUser.ComparePasswords(user.Password); err != nil {
		return nil, httperrors.NewRestError(http.StatusBadRequest, "Invalid password", nil)
	}

	foundedUser.SanitizePassword()

	token, err := utils.GenerateJWTToken(foundedUser, s.cfg)
	if err != nil {
		return nil, fmt.Errorf("AuthService.Login.GenerateJWTToken: %w", err)
	}

	return &models.UserWithToken{User: foundedUser, Token: token}, nil
}

func (s *AuthService) ChangePassword(user *models.User, oldPass, newPass string) (*models.User, error) {
	if err := user.ComparePasswords(oldPass); err != nil {
		return nil, httperrors.NewRestError(http.StatusBadRequest, "Invalid password", nil)
	}

	user.Password = newPass

	user.Prepare()

	updatedUser, err := s.authPgRepo.UpdatePassword(user)
	if err != nil {
		return nil, fmt.Errorf("AuthService.ChangePassword.UpdatePassword: %w", err)
	}

	return updatedUser, nil
}

func (s *AuthService) GetUserById(id string) (*models.User, error) {
	cachedUser, err := s.authRedisRepo.GetUserById(id)
	if err != nil {
		log.Println(err)
	}

	if cachedUser != nil {
		return cachedUser, nil
	}

	user, err := s.authPgRepo.GetUserById(id)
	if err != nil {
		return nil, err
	}

	if err = s.authRedisRepo.SetUser(user, 3600); err != nil {
		log.Println(err)
	}

	user.SanitizePassword()

	return user, nil
}
