package session

import (
	"blog/internal/config"
	"blog/internal/models"
	"blog/internal/repository"
	"blog/internal/service"
)

type sessionService struct {
	sessionRepository repository.RedisSession
	cfg               *config.Config
}

func NewSessionService(sessionRepository repository.RedisSession, cfg *config.Config) service.Session {
	return &sessionService{
		sessionRepository: sessionRepository,
		cfg:               cfg,
	}
}

func (s *sessionService) CreateSession(session *models.Session, expire int) (string, error) {
	return s.sessionRepository.CreateSession(session, expire)
}

func (s *sessionService) GetSessionByID(sessionID string) (*models.Session, error) {
	return s.sessionRepository.GetSessionByID(sessionID)
}

func (s *sessionService) DeleteSessionByID(sessionID string) error {
	return s.sessionRepository.DeleteSessionByID(sessionID)
}
