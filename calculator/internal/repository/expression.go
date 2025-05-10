package repo

import (
	"context"
	"time"

	"github.com/dzherb/go_calculator/internal/storage"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

type ExpressionStatus string

const (
	ExpressionNew        ExpressionStatus = "new"
	ExpressionProcessing ExpressionStatus = "processing"
	ExpressionSucceed    ExpressionStatus = "succeed"
	ExpressionAborted    ExpressionStatus = "aborted"
	ExpressionFailed     ExpressionStatus = "failed"
)

type Expression struct {
	ID         uint64           `json:"id"`
	UserID     uint64           `json:"user_id"`
	Status     ExpressionStatus `json:"status"`
	Expression string           `json:"expression"`
	Result     *float64         `json:"result"`
	CreatedAt  time.Time        `json:"created_at"`
	UpdatedAt  time.Time        `json:"updated_at"`
}

type ExpressionRepository interface {
	Get(id uint64) (Expression, error)
	Create(expression Expression) (Expression, error)
	Update(expression Expression) (Expression, error)
	GetForUser(userID uint64) ([]Expression, error)
	Unprocessed() ([]Expression, error)
}

type ExpressionRepositoryImpl struct {
	db storage.Connection
}

func NewExpressionRepository() ExpressionRepository {
	return &ExpressionRepositoryImpl{
		db: storage.Conn(),
	}
}

func NewExpressionRepositoryFromTx(tx pgx.Tx) ExpressionRepository {
	return &ExpressionRepositoryImpl{
		db: tx,
	}
}

func (er *ExpressionRepositoryImpl) Get(id uint64) (Expression, error) {
	expr := Expression{}
	err := pgxscan.Get(
		context.Background(),
		er.db,
		&expr,
		`SELECT id, user_id, status, expression, result, created_at, updated_at
		FROM expressions
		WHERE id = $1;`,
		id,
	)

	if err != nil {
		return Expression{}, err
	}

	return expr, nil
}

func (er *ExpressionRepositoryImpl) Create(
	expr Expression,
) (Expression, error) {
	err := pgxscan.Get(
		context.Background(),
		er.db,
		&expr,
		`INSERT INTO expressions (user_id, status, expression)
		VALUES ($1, 'new', $2)
		RETURNING id, user_id, status, expression, result, created_at, updated_at;`,
		expr.UserID,
		expr.Expression,
	)

	if err != nil {
		return Expression{}, err
	}

	return expr, nil
}

func (er *ExpressionRepositoryImpl) Update(
	expr Expression,
) (Expression, error) {
	err := pgxscan.Get(
		context.Background(),
		er.db,
		&expr,
		`UPDATE expressions
		SET status = $2, result = $3
		WHERE id = $1
		RETURNING id, user_id, status, expression, result, created_at, updated_at;`,
		expr.ID,
		expr.Status,
		expr.Result,
	)

	if err != nil {
		return Expression{}, err
	}

	return expr, nil
}

func (er *ExpressionRepositoryImpl) GetForUser(
	userID uint64,
) ([]Expression, error) {
	var exprs []Expression
	err := pgxscan.Select(
		context.Background(),
		er.db,
		&exprs,
		`SELECT id, user_id, status, expression, result, created_at, updated_at
		FROM expressions
		WHERE user_id = $1;`,
		userID,
	)

	if err != nil {
		return nil, err
	}

	return exprs, nil
}

func (er *ExpressionRepositoryImpl) Unprocessed() ([]Expression, error) {
	var exprs []Expression
	err := pgxscan.Select(
		context.Background(),
		er.db,
		&exprs,
		`SELECT id, user_id, status, expression, result, created_at, updated_at
		FROM expressions
		WHERE status IN ('new', 'processing');`,
	)

	if err != nil {
		return nil, err
	}

	return exprs, nil
}
