package repo_test

import (
	"os"
	"testing"

	"github.com/dzherb/go_calculator/calculator/internal/storage"
)

func TestMain(m *testing.M) {
	code := storage.RunTestsWithTempDB(
		func() int {
			return storage.RunTestsWithMigratedDB(m.Run)
		},
	)

	os.Exit(code)
}
