package domain

import "time"

type TokenManager interface {
	NewJWT(userID int64, ttl time.Duration) (string, error)
	Parse(accessToken string) (int64, error)
}
