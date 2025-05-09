package security

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey []byte

var accessTokenTTL time.Duration

func Init(cfg Config) {
	secretKey = []byte(cfg.SecretKey)
	accessTokenTTL = cfg.AccessTokenTTL
}

func IssueAccessToken(userID int) (string, error) {
	now := time.Now().UTC()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iat": timeToFloat64(now),
		"sub": strconv.Itoa(userID),
		"exp": timeToFloat64(now.Add(accessTokenTTL)),
	})

	return token.SignedString(secretKey)
}

func ValidateToken(tokenString string) (int, error) {
	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		},
		jwt.WithIssuedAt(),
		jwt.WithExpirationRequired(),
	)

	if err != nil {
		return 0, err
	}

	sub, err := token.Claims.GetSubject()
	if err != nil {
		return 0, err
	}

	userID, err := strconv.Atoi(sub)
	if err != nil {
		return 0, err
	}

	return userID, nil
}
