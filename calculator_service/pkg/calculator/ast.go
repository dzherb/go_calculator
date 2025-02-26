package calculator

import (
	"context"
	"strconv"
)

type node interface {
	evaluate(ctx context.Context, resChan chan float64)
}

type numberNode struct {
	value float64
}

func (n *numberNode) evaluate(ctx context.Context, resChan chan float64) {
	resChan <- n.value
}

type operatorNode struct {
	operator string
	left     node
	right    node
}

func (o *operatorNode) evaluate(ctx context.Context, resChan chan float64) {
	leftValChan := make(chan float64)
	rightValChan := make(chan float64)
	go o.left.evaluate(ctx, leftValChan)
	go o.right.evaluate(ctx, rightValChan)

	var leftVal, rightVal float64

	for _ = range 2 {
		select {
		case leftVal = <-leftValChan:
		case rightVal = <-rightValChan:
		case <-ctx.Done():
			return
		}

	}

	switch o.operator {
	case "+":
		resChan <- leftVal + rightVal
	case "-":
		resChan <- leftVal - rightVal
	case "*":
		resChan <- leftVal * rightVal
	case "/":
		resChan <- leftVal / rightVal
	}
}

// Переводит инфиксное выражение в RPN (обратную польскую нотацию)
func shuntingYard(tokens []token) []token {
	precedence := map[string]int{"+": 1, "-": 1, "*": 2, "/": 2}
	var output []token
	var operators []token

	for _, currToken := range tokens {

		if currToken.tokenType == number {
			output = append(output, currToken)
		} else if currToken.tokenType == operator {
			for len(operators) > 0 {
				top := operators[len(operators)-1]
				if precedence[top.value] >= precedence[currToken.value] {
					output = append(output, top)
					operators = operators[:len(operators)-1]
				} else {
					break
				}
			}
			operators = append(operators, currToken)
		} else if currToken.tokenType == openingBracket {
			operators = append(operators, currToken)
		} else if currToken.tokenType == closingBracket {
			for len(operators) > 0 && operators[len(operators)-1].tokenType != openingBracket {
				output = append(output, operators[len(operators)-1])
				operators = operators[:len(operators)-1]
			}
			if len(operators) > 0 && operators[len(operators)-1].tokenType == openingBracket {
				operators = operators[:len(operators)-1]
			}
		}
	}

	for len(operators) > 0 {
		output = append(output, operators[len(operators)-1])
		operators = operators[:len(operators)-1]
	}

	return output
}

// Строит абстрактое ситактическое дерево на основе
// последовательности токенов в обратной польской нотации
func buildAST(rpnOrganizedTokens []token) node {
	var stack []node

	for _, currToken := range rpnOrganizedTokens {
		if currToken.tokenType == number {
			val, _ := strconv.ParseFloat(currToken.value, 64)
			stack = append(stack, &numberNode{value: val})
		} else {
			right := stack[len(stack)-1]
			left := stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			stack = append(stack, &operatorNode{
				operator: currToken.value,
				left:     left,
				right:    right,
			})
		}
	}

	return stack[0]
}
