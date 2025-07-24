package token

// PasetoMaker implements Paseto v2‑local (symmetric‑key) token generation.
// All comments are in English.

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/o1egl/paseto"
)

// ---- Sentinel errors ----
var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

// Payload is embedded inside the Paseto token.
// You can extend it with roles, permissions, etc.

type Payload struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	IssuedAt  time.Time `json:"iat"`
	ExpiredAt time.Time `json:"exp"`
}

// CtxKey is the context key for storing payloads in context.Context.
const CtxKey = "payload"

// IsExpired checks if the token is expired.
func (p *Payload) IsExpired() bool {
	return time.Now().After(p.ExpiredAt)
}

// Maker defines operations for token generation and verification.

type Maker interface {
	CreateToken(userID uuid.UUID, duration time.Duration) (string, *Payload, error)
	VerifyToken(token string) (*Payload, error)
}

// PasetoMaker is a concrete implementation of Maker using Paseto v2‑local.

type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

// NewPasetoMaker creates a PasetoMaker with the provided 32‑byte key.
func NewPasetoMaker(key string) (*PasetoMaker, error) {
	if len(key) != 32 {
		return nil, fmt.Errorf("invalid key size: must be 32, got %d", len(key))
	}
	return &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(key),
	}, nil
}

// CreateToken generates a new token for a specific userID and duration.
func (m *PasetoMaker) CreateToken(userID uuid.UUID, duration time.Duration) (string, *Payload, error) {
	payload := &Payload{
		ID:        uuid.New(),
		UserID:    userID,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	token, err := m.paseto.Encrypt(m.symmetricKey, payload, nil)
	return token, payload, err
}

// VerifyToken decrypts the token and validates its payload.
func (m *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	var payload Payload
	if err := m.paseto.Decrypt(token, m.symmetricKey, &payload, nil); err != nil {
		return nil, ErrInvalidToken
	}
	if payload.IsExpired() {
		return nil, ErrExpiredToken
	}
	return &payload, nil
}
