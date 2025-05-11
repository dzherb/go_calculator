package repo_test

import (
	"strings"
	"testing"
	"time"

	"github.com/dzherb/go_calculator/calculator/internal/repository"
	"github.com/dzherb/go_calculator/calculator/internal/storage"
)

func testUser() repo.User {
	return repo.User{
		Username: "test",
		Password: "test_pass",
	}
}

func TestUserRepository_Create(t *testing.T) {
	storage.TestWithTransaction(t)

	now := time.Now().Add(-time.Second * 10)
	ur := repo.NewUserRepository()

	userToCreate := testUser()

	user, err := ur.Create(userToCreate)
	if err != nil {
		t.Fatal(err)
	}

	if user.ID == 0 {
		t.Error("user ID is zero")
	}

	if user.Username != userToCreate.Username {
		t.Errorf(
			"expected username %q, got %q",
			userToCreate.Username,
			user.Username,
		)
	}

	if user.Password == userToCreate.Password {
		t.Errorf("expected password to be hashed")
	}

	if user.CreatedAt.Before(now) {
		t.Errorf("created_at %s is earlier than expected", user.CreatedAt)
	}

	if user.UpdatedAt.Before(now) {
		t.Errorf("updated_at %s is earlier than expected", user.UpdatedAt)
	}
}

func TestUserConstraints(t *testing.T) {
	firstUser := testUser()

	users := []repo.User{
		{
			Username: firstUser.Username,
			Password: "test_pass",
		},
	}

	for _, u := range users {
		t.Run(u.Username, func(t *testing.T) {
			storage.TestWithTransaction(t)

			ur := repo.NewUserRepository()

			_, err := ur.Create(firstUser)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			_, err = ur.Create(u)
			if err == nil {
				t.Error("expected error")
				return
			}

			if !strings.Contains(err.Error(), "unique constraint") {
				t.Errorf("expected unique constraint error, got: %v", err)
			}
		})
	}
}

func TestUserRepository_GetByCredentials(t *testing.T) {
	storage.TestWithTransaction(t)

	ur := repo.NewUserRepository()

	created, err := ur.Create(testUser())
	if err != nil {
		t.Fatal(err)
	}

	got, err := ur.GetByCredentials(created.Username, "test_pass")
	if err != nil {
		t.Error(err)
		return
	}

	if got != created {
		t.Errorf("expected user to be %v, got %v", created, got)
	}
}

func TestUserRepository_GetByCredentials2(t *testing.T) {
	storage.TestWithTransaction(t)

	ur := repo.NewUserRepository()

	created, err := ur.Create(testUser())
	if err != nil {
		t.Fatal(err)
	}

	_, err = ur.GetByCredentials(created.Username, "wrong_pass")
	if err == nil {
		t.Error("expected error")
	}
}

func TestUserRepository_Get(t *testing.T) {
	storage.TestWithTransaction(t)

	ur := repo.NewUserRepository()

	created, err := ur.Create(testUser())
	if err != nil {
		t.Fatal(err)
	}

	got, err := ur.Get(created.ID)

	if err != nil {
		t.Fatal(err)
	}

	if got != created {
		t.Errorf("expected user %+v, got %+v", created, got)
	}
}
