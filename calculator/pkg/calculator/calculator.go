package calc

func Calculate(expression string) (float64, error) {
	exp, err := NewExpression(expression)
	if err != nil {
		return 0, err
	}

	err = simpleEvaluation(exp)
	if err != nil {
		return 0, err
	}

	res, err := exp.GetResult()

	return res, err
}
