package repository

import "blog/internal/models"

type PgAuth interface {
	CreateUser(user *models.User) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUserById(id string) (*models.User, error)
	UpdatePassword(user *models.User) (*models.User, error)
}

type RedisAuth interface {
	GetUserById(id string) (*models.User, error)
	SetUser(user *models.User, ttl int) error
}

type RedisSession interface {
	CreateSession(session *models.Session, expire int) (string, error)
	GetSessionByID(sessionID string) (*models.Session, error)
	DeleteSessionByID(sessionID string) error
}
