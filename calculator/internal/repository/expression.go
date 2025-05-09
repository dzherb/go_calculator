package repo

import (
	"time"
)

type ExpressionStatus string

const (
	New        ExpressionStatus = "new"
	Processing ExpressionStatus = "processing"
	Succeed    ExpressionStatus = "succeed"
	Aborted    ExpressionStatus = "aborted"
	Failed     ExpressionStatus = "failed"
)

type Expression struct {
	ID         int              `json:"id"`
	UserID     int              `json:"user_id"`
	Status     ExpressionStatus `json:"status"`
	Expression string           `json:"expression"`
	Result     *float64         `json:"result"`
	CreatedAt  time.Time        `json:"created_at"`
	UpdatedAt  time.Time        `json:"updated_at"`
}
