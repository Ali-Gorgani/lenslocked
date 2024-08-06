package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/Ali-Gorgani/lenslocked/rand"
)

const (
	DefaultResetDuration = 1 * time.Hour
)

type PasswordReset struct {
	ID     int
	UserID int
	// Token is only set when PasswordReset is being created.
	Token     string
	TokenHash string
	ExpiresAt time.Time
}

type PasswordResetService struct {
	DB *sql.DB
	// BytesPerToken specifies the number of bytes to be used when generating PasswordReset tokens.
	// If the value is less than MinBytesPerToken, MinBytesPerToken will be used instead.
	BytesPerToken int
	// Duration is a mount of time which a PasswordReset token is valid.
	// Defaults to DefaultResetDuration.
	Duration time.Duration
}

func (service *PasswordResetService) Create(email string) (*PasswordReset, error) {
	email = strings.ToLower(email)
	var userID int
	row := service.DB.QueryRow(`
		SELECT id
		FROM users
		WHERE email = $1`, email)
	err := row.Scan(&userID)
	if err != nil {
		return nil, fmt.Errorf("Create passwordReset: %w", err)
	}

	bytesPerToken := service.BytesPerToken
	if bytesPerToken < MinBytesPerToken {
		bytesPerToken = MinBytesPerToken
	}
	token, err := rand.String(bytesPerToken)
	if err != nil {
		return nil, fmt.Errorf("Create passwordResetToken: %w", err)
	}
	tokenHash := service.hash(token)

	duration := service.Duration
	if duration <= 0 {
		duration = DefaultResetDuration
	}
	expiresAt := time.Now().Add(duration)

	passwordReset := PasswordReset{
		Token:     token,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
	}

	row = service.DB.QueryRow(`
		INSERT INTO password_resets (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3) ON CONFLICT (user_id) DO
		UPDATE
		SET token_hash = $2, expires_at = $3
		RETURNING id`, userID, passwordReset.TokenHash, passwordReset.ExpiresAt)
	err = row.Scan(&passwordReset.ID)
	if err != nil {
		return nil, fmt.Errorf("Create passwordReset: %w", err)
	}

	return &passwordReset, nil
}

func (service *PasswordResetService) Consume(token string) (*User, error) {
	tokenHash := service.hash(token)
	row := service.DB.QueryRow(`
		SELECT password_resets.id,
			password_resets.expires_at,
			users.id,
			users.email,
			users.password_hash
		FROM password_resets
			JOIN users ON password_resets.user_id = users.id
		WHERE token_hash = $1`, tokenHash)
	var user User
	var pwReset PasswordReset
	err := row.Scan(&pwReset.ID, &pwReset.ExpiresAt, &user.ID, &user.Email, &user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("Consume passwordReset: %w", err)
	}

	if time.Now().After(pwReset.ExpiresAt) {
		return nil, fmt.Errorf("token expired: %v", token)
	}

	err = service.delete(pwReset.ID)
	if err != nil {
		return nil, fmt.Errorf("Consume passwordReset: %w", err)
	}

	return &user, nil
}

func (service *PasswordResetService) delete(id int) error {
	_, err := service.DB.Exec(`
		DELETE FROM password_resets
		WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("Delete passwordReset: %w", err)
	}
	return nil
}

func (service *PasswordResetService) hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}
