package service

import "blog/internal/models"

type Auth interface {
	Register(*models.User) (*models.UserWithToken, error)
	Login(*models.User) (*models.UserWithToken, error)
	ChangePassword(user *models.User, oldPass, newPass string) (*models.User, error)
	GetUserById(id string) (*models.User, error)
}

type Session interface {
	CreateSession(session *models.Session, expire int) (string, error)
	GetSessionByID(sessionID string) (*models.Session, error)
	DeleteSessionByID(sessionID string) error
}
