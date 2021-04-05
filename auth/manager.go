package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"

	"github.com/rebel-l/auth-service/user/usermodel"
)

const (
	// TokenTypeAccess identifies the access token.
	TokenTypeAccess = "access"

	// TokenTypeRefresh identifies the refresh token.
	TokenTypeRefresh = "refresh"

	// LifetimeAccessToken defines the lifetime for the access token.
	LifetimeAccessToken = time.Minute * 15

	// LifetimeRefreshToken defines the lifetime for the refresh token.
	LifetimeRefreshToken = time.Hour * 24 * 7
)

var (
	// ErrNoTokenSecret indicates that the secret for the JWT is not defined.
	ErrNoTokenSecret = errors.New("secret is not set or empty")

	// ErrNoUser indicates that the user is not defined.
	ErrNoUser = errors.New("user should be not nil")
)

// TokenGenerator defines an interfaces to generate JWT tokens.
type TokenGenerator interface {
	GenerateTokens(user *usermodel.User) (map[string]*Token, error)
}

// Manager take care on token handling.
type Manager struct {
	method  jwt.SigningMethod
	secrets map[string]string

	// TODO: we need a store ... redis?
}

// NewManager returns an instance of a manager to handle JWT tokens.
func NewManager(accessSecret, refreshSecret string) *Manager {
	return &Manager{
		method: jwt.SigningMethodHS256,
		secrets: map[string]string{
			TokenTypeAccess:  accessSecret,
			TokenTypeRefresh: refreshSecret,
		},
	}
}

// GenerateTokens returns an access and a refresh token. It is important that the secrets are defined (not empty string)
// on the Manager otherwise you'll get a ErrNoTokenSecret error for security reasons.
// It expects a user model as parameter. It returns ErrNoUser id user is nil.
func (m *Manager) GenerateTokens(user *usermodel.User) (map[string]*Token, error) {
	if user == nil {
		return nil, ErrNoUser
	}

	tokens := make(map[string]*Token)

	var err error

	tokens[TokenTypeAccess], err = m.createToken(TokenTypeAccess, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create access token: %w", err)
	}

	tokens[TokenTypeRefresh], err = m.createToken(TokenTypeRefresh, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create refresh token: %w", err)
	}

	return tokens, nil
}

func (m *Manager) createToken(tokenType string, user *usermodel.User) (*Token, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("failed to create UUID for token: %w", err)
	}

	claims := jwt.MapClaims{
		"id":     id.String(),
		"userID": user.ID.String(),
		"type":   tokenType,
	}

	var expires int64

	switch tokenType {
	case TokenTypeAccess:
		expires = time.Now().Add(LifetimeAccessToken).Unix() // TODO: maybe duration only?

		claims["authorized"] = true
	case TokenTypeRefresh:
		expires = time.Now().Add(LifetimeRefreshToken).Unix()
	}

	claims["exp"] = expires

	secret, ok := m.secrets[tokenType]
	if !ok || secret == "" {
		return nil, fmt.Errorf("%s token: %w", tokenType, ErrNoTokenSecret)
	}

	return NewToken(id, expires, jwt.NewWithClaims(m.method, claims), secret)
}
