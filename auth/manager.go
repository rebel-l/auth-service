package auth

import (
	"errors"
	"fmt"
	"net/http"
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

	// ErrNoStore indicates that no storage for the tokens was setup.
	ErrNoStore = errors.New("no token store set")

	// ErrUnexpectedSigningMethod indicates that the signing method of the token is wrong.
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")

	// ErrInvalidToken indicates that the was manipulated and is not valid anymore.
	ErrInvalidToken = errors.New("token is invalid")
)

// TokenGenerator defines an interfaces to generate JWT tokens.
type TokenGenerator interface {
	GenerateTokens(user *usermodel.User) (map[string]*Token, error)
}

type TokenManager interface {
	DeleteTokens(request *http.Request) error
}

type TokenValidator interface {
	IsAccessTokenValid(header http.Header) bool
}

// Manager take care on token handling.
type Manager struct {
	method  jwt.SigningMethod
	secrets map[string]string
	store   Storage
}

// NewManager returns an instance of a manager to handle JWT tokens.
func NewManager(accessSecret, refreshSecret string, store Storage) *Manager {
	return &Manager{
		method: jwt.SigningMethodHS256,
		secrets: map[string]string{
			TokenTypeAccess:  accessSecret,
			TokenTypeRefresh: refreshSecret,
		},
		store: store,
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

	if err := m.save(tokens); err != nil {
		return nil, err
	}

	return tokens, nil
}

// DeleteTokens deletes given tokens in request from storage.
func (m *Manager) DeleteTokens(request *http.Request) error {
	tokenID, err := m.ExtractAccessTokenID(request.Header)
	if err != nil {
		return fmt.Errorf("failed to delete access token: %w", err)
	}

	// TODO: delete refresh token too
	if err := m.store.Del(tokenID).Err(); err != nil {
		return fmt.Errorf("failed to delete tokens from store: %w", err)
	}

	return nil
}

// IsAccessTokenValid returns true if the access token can be extracted from header and is valid.
func (m *Manager) IsAccessTokenValid(header http.Header) bool {
	bearerToken := extractToken(header)

	token, err := m.verifyToken(bearerToken, TokenTypeAccess)
	if err != nil {
		return false
	}

	if _, ok := token.Claims.(jwt.Claims); !ok || !token.Valid {
		return false
	}

	return true
}

// ExtractAccessTokenID returns token from request header.
func (m *Manager) ExtractAccessTokenID(header http.Header) (string, error) {
	bearerToken := extractToken(header)

	token, err := m.verifyToken(bearerToken, TokenTypeAccess)
	if err != nil {
		return "", fmt.Errorf("failed to extract access token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", fmt.Errorf("%s %w", TokenTypeAccess, ErrInvalidToken)
	}

	id, ok := claims["id"].(string)
	if !ok {
		return "", fmt.Errorf("%s %w: no ID", TokenTypeAccess, ErrInvalidToken)
	}

	return id, nil
}

func (m *Manager) verifyToken(bearerToken string, tokenType string) (*jwt.Token, error) {
	return jwt.Parse(bearerToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%w: %v", ErrUnexpectedSigningMethod, token.Header["alg"])
		}

		return []byte(m.secrets[tokenType]), nil
	})
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

	expires := time.Now()

	switch tokenType {
	case TokenTypeAccess:
		expires = expires.Add(LifetimeAccessToken)

		claims["authorized"] = true
	case TokenTypeRefresh:
		expires = expires.Add(LifetimeRefreshToken)
	}

	claims["exp"] = expires.Unix()

	secret, ok := m.secrets[tokenType]
	if !ok || secret == "" {
		return nil, fmt.Errorf("%s token: %w", tokenType, ErrNoTokenSecret)
	}

	return NewToken(id, expires, jwt.NewWithClaims(m.method, claims), secret)
}

func (m *Manager) save(tokens map[string]*Token) error {
	if m.store == nil {
		return ErrNoStore
	}

	for k, v := range tokens {
		exp := time.Until(v.Expires)

		res := m.store.Set(v.ID.String(), v.JWT, exp)
		if res != nil && res.Err() != nil {
			return fmt.Errorf("failed to store token %s/%d: %w", k, len(tokens), res.Err())
		}
	}

	return nil
}
