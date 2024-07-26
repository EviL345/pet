package session

import (
	"blog/internal/models"
	"context"
	"fmt"
	"github.com/go-redis/cache/v9"
	"github.com/google/uuid"
	"time"
)

const basePrefix = "api-session"

type SessionRepo struct {
	redis *cache.Cache
}

func NewSessionRepository(redisClient *cache.Cache) *SessionRepo {
	return &SessionRepo{
		redis: redisClient,
	}
}

func (s *SessionRepo) CreateSession(session *models.Session, ttl int) (string, error) {
	session.SessionID = uuid.New().String()
	sessionKey := s.createKey(session.SessionID)

	if err := s.redis.Set(&cache.Item{
		Key:   sessionKey,
		Value: session,
		TTL:   time.Duration(ttl) * time.Second,
	}); err != nil {
		return "", fmt.Errorf("SessionRepo.CreateSession.redis.Set: %w", err)
	}

	return sessionKey, nil
}

func (s *SessionRepo) GetSessionByID(sessionID string) (*models.Session, error) {

	session := &models.Session{}
	if err := s.redis.Get(context.Background(), sessionID, &session); err != nil {
		return nil, fmt.Errorf("SessionRepo.GetSessionByID.redis.Get: %w", err)
	}

	return session, nil
}
func (s *SessionRepo) DeleteSessionByID(sessionID string) error {

	if err := s.redis.Delete(context.Background(), sessionID); err != nil {
		return fmt.Errorf("SessionRepo.DeleteSessionByID.redis.Delete: %w", err)
	}

	return nil

}

func (s *SessionRepo) createKey(sessionID string) string {
	return fmt.Sprintf("%s: %s", basePrefix, sessionID)
}
