package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"

	"github.com/Ali-Gorgani/lenslocked/rand"
)

const (
	// The minimum number of bytes to be used for a session token.
	MinBytesPerToken = 32
)

type Session struct {
	ID     int
	UserID int
	// Token is only set when creating a new session. It will be empty when reading from the database.
	// Only the TokenHash will store in the database.
	Token     string
	TokenHash string
}

type SessionService struct {
	DB *sql.DB
	// BytesPerToken specifies the number of bytes to be used when generating session tokens.
	// If the value is less than MinBytesPerToken, MinBytesPerToken will be used instead.
	BytesPerToken int
}

func (service *SessionService) Create(userID int) (*Session, error) {
	bytesPerToken := service.BytesPerToken
	if bytesPerToken < MinBytesPerToken {
		bytesPerToken = MinBytesPerToken
	}
	token, err := rand.String(bytesPerToken)
	if err != nil {
		return nil, fmt.Errorf("Create sessionToken: %w", err)
	}

	tokenHash := service.hash(token)

	session := Session{
		UserID:    userID,
		Token:     token,
		TokenHash: tokenHash,
	}

	row := service.DB.QueryRow(`
		INSERT INTO sessions (user_id, token_hash)
		VALUES ($1, $2) ON CONFLICT (user_id) DO
		UPDATE
		SET token_hash = $2
		RETURNING id
	`, session.UserID, session.TokenHash)
	err = row.Scan(&session.ID)
	if err != nil {
		return nil, fmt.Errorf("Create session: %w", err)
	}
	
	return &session, nil
}

func (service *SessionService) User(token string) (*User, error) {
	tokenHash := service.hash(token)

	row := service.DB.QueryRow(`
		SELECT users.id,
			users.email,
			users.password_hash
		FROM sessions
			JOIN users ON sessions.user_id = users.id
		WHERE sessions.token_hash = $1
	`, tokenHash)

	var user User
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("User: %w", err)
	}

	return &user, nil
}

func (service *SessionService) Delete(token string) error {
	tokenHash := service.hash(token)
	_, err := service.DB.Exec("DELETE FROM sessions WHERE token_hash = $1", tokenHash)
	if err != nil {
		return fmt.Errorf("Delete session: %w", err)
	}
	return nil
}

func (ss *SessionService) hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}
