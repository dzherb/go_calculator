package storage_test

import (
	"os"
	"testing"

	"github.com/dzherb/go_calculator/calculator/internal/storage"
)

func TestMain(m *testing.M) {
	code := storage.RunTestsWithTempDB(m.Run)

	os.Exit(code)
}

func TestTempDB(t *testing.T) {
	pool := storage.ActivePool()

	dbName := pool.Config().ConnConfig.Database
	if dbName != "go_test" {
		t.Errorf("db name is %s, expected %s", dbName, "go_test")
	}

	err := pool.Ping(t.Context())
	if err != nil {
		t.Error(err)
	}
}
