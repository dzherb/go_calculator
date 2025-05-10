package calc_test

import (
	"testing"

	"github.com/dzherb/go_calculator/calculator/pkg/calculator"
)

func TestCalculator(t *testing.T) {
	testCasesSuccess := []struct {
		name           string
		expression     string
		expectedResult float64
	}{
		{
			name:           "simple",
			expression:     "1+1",
			expectedResult: 2,
		},
		{
			name:           "priority",
			expression:     "2+2*2",
			expectedResult: 6,
		},
		{
			name:           "priority with brackets",
			expression:     "(2+2)*2",
			expectedResult: 8,
		},
		{
			name:           "dot at the end",
			expression:     "1./2",
			expectedResult: 0.5,
		},
		{
			name:           "whitespace",
			expression:     "1 -   0.5 + 3-(3  *2)",
			expectedResult: -2.5,
		},
		{
			name:           "unary minus",
			expression:     "-2",
			expectedResult: -2,
		},
		{
			name:           "complex expression",
			expression:     "-1.2*(3-2) / 10 + (-4 - 3.5)",
			expectedResult: -7.62,
		},
		{
			name:           "negative numbers order",
			expression:     "55-12/(3-4)-17",
			expectedResult: 50,
		},
	}

	for _, testCase := range testCasesSuccess {
		t.Run(testCase.name, func(t *testing.T) {
			val, err := calc.Calculate(testCase.expression)
			if err != nil {
				t.Fatalf(
					"successful case %s returns error: %s",
					testCase.expression,
					err.Error(),
				)
			}

			if val != testCase.expectedResult {
				t.Fatalf("%f should be equal %f", val, testCase.expectedResult)
			}
		})
	}

	testCasesFail := []struct {
		name        string
		expression  string
		expectedErr error
	}{
		{
			name:       "unexpected char",
			expression: "1+a",
		},
		{
			name:       "Operator at the end",
			expression: "1+1*",
		},
		{
			name:       "doubling Operator",
			expression: "2+2**2",
		},
		{
			name:       "brackets don't match",
			expression: "(2+2))-(2+3)",
		},
		{
			name:       "empty expression",
			expression: "",
		},
		{
			name:       "no operators",
			expression: "1 23 7.4",
		},
	}

	for _, testCase := range testCasesFail {
		t.Run(testCase.name, func(t *testing.T) {
			val, err := calc.Calculate(testCase.expression)
			if err == nil {
				t.Fatalf(
					"expression %s is invalid but result  %f was obtained",
					testCase.expression,
					val,
				)
			}
		})
	}
}
