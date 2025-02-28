package calculator

import (
	"context"
	"fmt"
	"time"
)

func Calculate(expression string) (float64, error) {
	tokens, err := Tokenize(expression)
	if err != nil {
		return 0, err
	}

	// Переводим токены в обратную польскую нотацию (RPN)
	rpnOrganizedTokens := shuntingYard(tokens)

	ast := buildAST(rpnOrganizedTokens)

	resultChan := make(chan float64)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go ast.evaluate(ctx, resultChan)

	select {
	case res := <-resultChan:
		return res, nil
	case <-time.After(time.Millisecond * 300):
		cancel()
		return 0, fmt.Errorf("calculation timed out")
	}
}
