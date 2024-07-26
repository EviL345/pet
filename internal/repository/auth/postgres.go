package auth

import (
	"blog/internal/models"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type AuthPgRepo struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *AuthPgRepo {
	return &AuthPgRepo{db: db}
}

func (r *AuthPgRepo) GetUserByEmail(email string) (*models.User, error) {
	query := fmt.Sprint("SELECT (*) FROM users WHERE email = $1")
	user := &models.User{}

	if err := r.db.QueryRowx(query, email).StructScan(user); err != nil {
		return nil, fmt.Errorf("AuthPgRepo.GetUserByEmail.QueryRowx: %w", err)
	}

	return user, nil
}

func (r *AuthPgRepo) CreateUser(user *models.User) (*models.User, error) {
	query := fmt.Sprintf("INSERT INTO users (id, email, password) VALUES ($1, $2, $3) RETURNING (*)")
	createdUser := &models.User{}

	if err := r.db.QueryRowx(query, user.Email, user.Password).StructScan(createdUser); err != nil {
		return nil, fmt.Errorf("AuthPgRepo.CreateUser.QueryRowx: %w", err)
	}

	return createdUser, nil
}

func (r *AuthPgRepo) GetUserById(id string) (*models.User, error) {
	query := fmt.Sprint("SELECT (*) FROM users WHERE id = $1")
	user := &models.User{}

	if err := r.db.QueryRowx(query, id).StructScan(user); err != nil {
		return nil, fmt.Errorf("AuthPgRepo.GetUserById.QueryRowx: %w", err)
	}

	return user, nil
}

func (r *AuthPgRepo) UpdatePassword(user *models.User) (*models.User, error) {
	query := fmt.Sprint("UPDATE users SET password = $1 WHERE id = $2 RETURNING (*)")
	updatedUser := &models.User{}

	if err := r.db.QueryRowx(query, user.Password, user.ID).StructScan(updatedUser); err != nil {
		return nil, fmt.Errorf("AuthPgRepo.UpdatePassword.QueryRowx: %w", err)
	}

	return updatedUser, nil
}
