package auth_test

import (
	"os"
	"testing"
	"time"

	"github.com/dzherb/go_calculator/internal/auth"
	repo "github.com/dzherb/go_calculator/internal/repository"
	"github.com/dzherb/go_calculator/internal/storage"
	"github.com/dzherb/go_calculator/pkg/security"
)

func TestMain(m *testing.M) {
	security.Init(security.Config{
		SecretKey:      "secret",
		AccessTokenTTL: time.Hour,
	})

	code := storage.RunTestsWithTempDB(func() int {
		return storage.RunTestsWithMigratedDB(m.Run)
	})
	os.Exit(code)
}

func TestLogin(t *testing.T) {
	storage.TestWithTransaction(t)

	ur := repo.NewUserRepository()
	user, err := ur.Create(repo.User{
		Username: "test",
		Password: "test_pass",
	})

	if err != nil {
		t.Fatal(err)
	}

	as := auth.NewService()

	got, err := as.Login(user.Username, "test_pass")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if got.User != user {
		t.Errorf("got user %v, want %v", got.User, user)
		return
	}

	userID, err := security.ValidateToken(got.Token)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if got.User.ID != userID {
		t.Errorf("got userID %v, want %v", userID, got.User.ID)
	}
}

func TestRegister(t *testing.T) {
	storage.TestWithTransaction(t)

	as := auth.NewService()

	res, err := as.Register("user", "strongPass11")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	ur := repo.NewUserRepository()

	userFromDB, err := ur.Get(res.User.ID)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if res.User != userFromDB {
		t.Errorf("got user %v, want %v", userFromDB, res.User)
	}

	userID, err := security.ValidateToken(res.Token)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if res.User.ID != userID {
		t.Errorf("got userID %v, want %v", userID, res.User.ID)
	}
}
