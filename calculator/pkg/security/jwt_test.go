package security_test

import (
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/dzherb/go_calculator/calculator/pkg/security"
	"github.com/golang-jwt/jwt/v5"
)

type token struct {
	iat float64
	sub string
	exp float64
}

func TestMain(m *testing.M) {
	security.Init(security.Config{
		SecretKey:      "secret",
		AccessTokenTTL: time.Hour,
	})

	os.Exit(m.Run())
}

func TestTokenIssue(t *testing.T) {
	token, err := security.IssueAccessToken(23)
	if err != nil {
		t.Error(err)
		return
	}

	userID, err := security.ValidateToken(token)
	if err != nil {
		t.Error(err)
		return
	}

	if userID != 23 {
		t.Errorf("got sub %d, expected 23", userID)
	}
}

func TestTokenValidation(t *testing.T) {
	cases := []struct {
		token   token
		isValid bool
	}{
		{
			token: token{
				iat: security.TimeToFloat64(time.Now()),
				sub: "1",
				exp: security.TimeToFloat64(time.Now().Add(time.Second * 100)),
			},
			isValid: true,
		},
		{
			token: token{
				iat: security.TimeToFloat64(time.Now()),
				sub: "1",
				exp: security.TimeToFloat64(time.Now().Add(-time.Second * 100)),
			},
			isValid: false,
		},
		{
			token: token{
				iat: security.TimeToFloat64(time.Now().Add(time.Second * 100)),
				sub: "1",
				exp: security.TimeToFloat64(time.Now().Add(time.Second * 100)),
			},
			isValid: false,
		},
		{
			token: token{
				iat: security.TimeToFloat64(time.Now()),
				sub: "not_a_number",
				exp: security.TimeToFloat64(time.Now().Add(time.Second * 100)),
			},
			isValid: false,
		},
	}

	for _, c := range cases {
		tokenEncoded, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"iat": c.token.iat,
			"sub": c.token.sub,
			"exp": c.token.exp,
		}).
			SignedString(*security.SecretKey)

		if err != nil {
			t.Error(err)
		}

		userID, err := security.ValidateToken(tokenEncoded)
		if err != nil {
			if c.isValid {
				t.Error("unexpected error:", err)
			}

			continue
		}

		if !c.isValid {
			t.Error("token is unexpectedly valid")
			continue
		}

		expectedID, err := strconv.ParseUint(c.token.sub, 10, 64)
		if err != nil {
			t.Error(err)
			continue
		}

		if userID != expectedID {
			t.Errorf("got sub %d, want %d", userID, expectedID)
		}
	}
}
