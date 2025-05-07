package security

import "time"

type Config struct {
	SecretKey      string
	AccessTokenTTL time.Duration
}
