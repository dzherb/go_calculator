package calculator

import (
	"fmt"
	"strconv"
)

type node interface {
	String() string
}

type numberNode struct {
	value float64
}

func (n *numberNode) String() string {
	return fmt.Sprintf("%.2f", n.value)
}

type operatorNode struct {
	operator     string
	left         node
	right        node
	parent       *operatorNode // Ссылка на родителя
	isProcessing bool
	isProcessed  bool
}

func (o *operatorNode) String() string {
	return fmt.Sprintf(
		"(%s %s %s)",
		o.left.String(),
		o.operator,
		o.right.String(),
	)
}

func (o *operatorNode) nextReadyForProcessingNode() (*operatorNode, bool) {
	if o.isProcessed || o.isProcessing {
		return nil, false
	}

	// Проверяем, можно ли вычислить этот узел
	_, leftOK := o.left.(*numberNode)
	_, rightOK := o.right.(*numberNode)

	if leftOK && rightOK {
		o.isProcessing = true
		return o, true
	}

	// Рекурсивный поиск
	if leftOp, ok := o.left.(*operatorNode); ok {
		if node, ok := leftOp.nextReadyForProcessingNode(); ok {
			return node, true
		}
	}

	if rightOp, ok := o.right.(*operatorNode); ok {
		if node, ok := rightOp.nextReadyForProcessingNode(); ok {
			return node, true
		}
	}

	return nil, false
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

	root := stack[0]
	if r, ok := root.(*operatorNode); ok {
		addParents(r)
	}

	return root
}

func addParents(node *operatorNode) {
	if leftOp, ok := node.left.(*operatorNode); ok {
		leftOp.parent = node
		addParents(leftOp)
	}
	if rightOp, ok := node.right.(*operatorNode); ok {
		rightOp.parent = node
		addParents(rightOp)
	}
}
