package token

import (
	"fmt"
	"time"

	"github.com/o1egl/paseto"
	"golang.org/x/crypto/chacha20poly1305"
)

type PasetoManager struct {
	paseto       *paseto.V2
	symmetrickey []byte
}

func NewPasetoManager(symmetrickey string) (TokenManager, error) {
	if len(symmetrickey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be exactly %d characters", chacha20poly1305.KeySize)
	}

	m := &PasetoManager{
		paseto:       paseto.NewV2(),
		symmetrickey: []byte(symmetrickey),
	}
	return m, nil
}

func (m *PasetoManager) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}
	return m.paseto.Encrypt(m.symmetrickey, payload, nil)
}

func (m *PasetoManager) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	err := m.paseto.Decrypt(token, m.symmetrickey, payload, nil)
	if err != nil {
		return nil, err
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}
