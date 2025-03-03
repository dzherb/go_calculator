package calculator

func Calculate(expression string) (float64, error) {
	tokens, err := Tokenize(expression)
	if err != nil {
		return 0, err
	}

	// Переводим токены в обратную польскую нотацию (RPN)
	rpnOrganizedTokens := shuntingYard(tokens)

	ast := buildAST(rpnOrganizedTokens)

	exp := NewExpression(ast.(*operatorNode))

	err = simpleEvaluation(exp)
	if err != nil {
		return 0, err
	}

	res, err := exp.GetResult()

	return res, err
}
