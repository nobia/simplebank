package token

import "time"

// Maker is an interface for managing tokens
type TokenManager interface {
	CreateToken(username string, duration time.Duration) (string, error)
	VerifyToken(token string) (*Payload, error)
}
