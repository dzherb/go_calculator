package orchestrator_test

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/dzherb/go_calculator/calculator/internal/orchestrator"
	"github.com/dzherb/go_calculator/calculator/pkg/security"
)

var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	userID := r.Context().Value(orchestrator.UserIDKey).(uint64)
	w.Write([]byte(strconv.FormatUint(userID, 10))) //nolint:errcheck
})

func initSecurity() {
	security.Init(security.Config{
		SecretKey:      "secret",
		AccessTokenTTL: time.Hour,
	})
}

func TestAuthMiddlewareSuccess(t *testing.T) {
	initSecurity()

	userID := uint64(25)

	req := httptest.NewRequest(http.MethodGet, "/", nil)

	token, err := security.IssueAccessToken(userID)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()

	orchestrator.AuthRequired(handler).ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("got status %v, want %v", status, http.StatusOK)
	}

	if rr.Body.String() != strconv.FormatUint(userID, 10) {
		t.Errorf(
			"got body %v, want %v",
			rr.Body.String(),
			strconv.FormatUint(userID, 10),
		)
	}
}

func TestAuthMiddlewareFail(t *testing.T) {
	initSecurity()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer fake_token")

	rr := httptest.NewRecorder()

	orchestrator.AuthRequired(handler).ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("got status %v, want %v", status, http.StatusOK)
	}
}
