package repo_test

import (
	"testing"
	"time"

	repo "github.com/dzherb/go_calculator/internal/repository"
	"github.com/dzherb/go_calculator/internal/storage"
	"github.com/jackc/pgx/v5"
)

func createTestUser() (repo.User, error) {
	user, err := repo.NewUserRepository().Create(testUser())
	if err != nil {
		return repo.User{}, err
	}

	return user, nil
}

func TestExpressionRepository_Create(t *testing.T) {
	storage.TestWithTransaction(t)

	user, err := createTestUser()
	if err != nil {
		t.Fatal(err)
	}

	expr := repo.Expression{
		UserID:     user.ID,
		Expression: "2+4/2",
	}

	now := time.Now().Add(-time.Second * 10)
	er := repo.NewExpressionRepository()

	createdExpr, err := er.Create(expr)
	if err != nil {
		t.Errorf("error while creating expression: %v", err)
		return
	}

	if createdExpr.UserID != user.ID {
		t.Errorf(
			"createdExpr.UserID = %v, want %v",
			createdExpr.UserID,
			user.ID,
		)
	}

	if createdExpr.Expression != expr.Expression {
		t.Errorf(
			"createdExpr.Expression = %v, want %v",
			createdExpr.Expression,
			expr.Expression,
		)
	}

	if createdExpr.Status != repo.ExpressionNew {
		t.Errorf(
			"createdExpr.Status = %v, want %v",
			createdExpr.Status,
			repo.ExpressionNew,
		)
	}

	if createdExpr.CreatedAt.Before(now) {
		t.Errorf("created_at %s is earlier than expected", user.CreatedAt)
	}

	if createdExpr.UpdatedAt.Before(now) {
		t.Errorf("updated_at %s is earlier than expected", user.UpdatedAt)
	}
}

func float64Ptr(v float64) *float64 {
	return &v
}

func TestExpressionRepository_Update(t *testing.T) { //nolint:gocognit
	storage.TestWithTransaction(t)

	user, err := createTestUser()
	if err != nil {
		t.Fatal(err)
	}

	expr := repo.Expression{
		UserID:     user.ID,
		Expression: "2+4/2",
	}

	er := repo.NewExpressionRepository()

	expr, err = er.Create(expr)
	if err != nil {
		t.Fatal()
	}

	cases := []struct {
		expr repo.Expression
	}{
		{
			repo.Expression{
				ID:     expr.ID,
				Status: repo.ExpressionProcessing,
			},
		},
		{
			repo.Expression{
				ID:     expr.ID,
				UserID: 123,
				Status: repo.ExpressionAborted,
			},
		},
		{
			repo.Expression{
				ID:     expr.ID,
				Status: repo.ExpressionSucceed,
				Result: float64Ptr(4),
			},
		},
	}

	for _, c := range cases {
		t.Run(string(c.expr.Status), func(t *testing.T) {
			err = storage.WithTransaction(
				t.Context(),
				storage.Conn(),
				func(tx pgx.Tx) error {
					er = repo.NewExpressionRepositoryFromTx(tx)

					updated, err := er.Update(c.expr)
					if err != nil {
						return err
					}

					if updated.ID != expr.ID {
						t.Errorf(
							"updated.ID = %v, want %v",
							updated.ID,
							expr.ID,
						)
					}

					if updated.UserID != expr.UserID {
						t.Errorf(
							"updated.UserID = %v, want %v",
							updated.ID,
							expr.ID,
						)
					}

					if updated.Expression != expr.Expression {
						t.Errorf(
							"updated.Expression = %v, want %v",
							updated.Expression,
							expr.Expression,
						)
					}

					if updated.Status != c.expr.Status {
						t.Errorf(
							"updated.Status = %v, want %v",
							updated.Status,
							c.expr.Status,
						)
					}

					return nil
				},
			)

			if err != nil {
				t.Errorf("error while updating expression: %v", err)
			}
		})
	}
}

func TestExpressionRepository_Get(t *testing.T) {
	storage.TestWithTransaction(t)

	user, err := createTestUser()
	if err != nil {
		t.Fatal(err)
	}

	expr := repo.Expression{
		UserID:     user.ID,
		Expression: "2+4/2",
	}
	er := repo.NewExpressionRepository()

	expr, err = er.Create(expr)
	if err != nil {
		t.Fatal()
	}

	got, err := er.Get(expr.ID)
	if err != nil {
		t.Errorf("error while getting expression: %v", err)
		return
	}

	if got != expr {
		t.Errorf("got = %v, want %v", got, expr)
	}
}

func TestExpressionRepository_GetForUser(t *testing.T) {
	storage.TestWithTransaction(t)

	ur := repo.NewUserRepository()

	u1, err := ur.Create(testUser())
	if err != nil {
		t.Fatal(err)
	}

	u2, err := ur.Create(repo.User{
		Username: "user2",
		Password: "pass",
	})
	if err != nil {
		t.Fatal(err)
	}

	er := repo.NewExpressionRepository()

	expr1, err := er.Create(repo.Expression{
		UserID:     u1.ID,
		Expression: "2+4/2",
	})
	if err != nil {
		t.Fatal()
	}

	_, err = er.Create(repo.Expression{
		UserID:     u2.ID,
		Expression: "1-10",
	})
	if err != nil {
		t.Fatal()
	}

	res, err := er.GetForUser(u1.ID)
	if err != nil {
		t.Errorf("error while getting expressions for user: %v", err)
		return
	}

	if len(res) != 1 {
		t.Errorf("len(res) = %v, want %v", len(res), 1)
	}

	if res[0] != expr1 {
		t.Errorf("res = %v, want %v", res, expr1)
	}
}
