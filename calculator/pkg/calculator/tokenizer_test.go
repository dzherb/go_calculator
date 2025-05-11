package calc_test

import (
	"reflect"
	"testing"

	"github.com/dzherb/go_calculator/calculator/pkg/calculator"
)

func TestTokenizer(t *testing.T) {
	testCases := []struct {
		expression  string
		expected    []calc.Token
		expectError bool
	}{
		{
			expression:  "1",
			expected:    []calc.Token{{Value: "1", TokenType: calc.Number}},
			expectError: false,
		},
		{
			expression: "1+2-3",
			expected: []calc.Token{
				{Value: "1", TokenType: calc.Number},
				{Value: "+", TokenType: calc.Operator},
				{Value: "2", TokenType: calc.Number},
				{Value: "-", TokenType: calc.Operator},
				{Value: "3", TokenType: calc.Number},
			},
			expectError: false,
		},
		{
			expression: "1+2-3.4",
			expected: []calc.Token{
				{Value: "1", TokenType: calc.Number},
				{Value: "+", TokenType: calc.Operator},
				{Value: "2", TokenType: calc.Number},
				{Value: "-", TokenType: calc.Operator},
				{Value: "3.4", TokenType: calc.Number},
			},
			expectError: false,
		},
		{
			expression: "1+(242-3.4)/33",
			expected: []calc.Token{
				{Value: "1", TokenType: calc.Number},
				{Value: "+", TokenType: calc.Operator},
				{Value: "(", TokenType: calc.OpeningBracket},
				{Value: "242", TokenType: calc.Number},
				{Value: "-", TokenType: calc.Operator},
				{Value: "3.4", TokenType: calc.Number},
				{Value: ")", TokenType: calc.ClosingBracket},
				{Value: "/", TokenType: calc.Operator},
				{Value: "33", TokenType: calc.Number},
			},
			expectError: false,
		},
		{
			expression: "45*4+(2.42-3.4)/33",
			expected: []calc.Token{
				{Value: "45", TokenType: calc.Number},
				{Value: "*", TokenType: calc.Operator},
				{Value: "4", TokenType: calc.Number},
				{Value: "+", TokenType: calc.Operator},
				{Value: "(", TokenType: calc.OpeningBracket},
				{Value: "2.42", TokenType: calc.Number},
				{Value: "-", TokenType: calc.Operator},
				{Value: "3.4", TokenType: calc.Number},
				{Value: ")", TokenType: calc.ClosingBracket},
				{Value: "/", TokenType: calc.Operator},
				{Value: "33", TokenType: calc.Number},
			},
			expectError: false,
		},
		{
			expression:  "1+(242-3.4.3)/33",
			expected:    nil,
			expectError: true,
		},
		{
			expression:  "1+(242-.3)/33",
			expected:    nil,
			expectError: true,
		},
		{
			expression:  "45*4+(2.42-3.4)/3..3",
			expected:    nil,
			expectError: true,
		},
		{
			expression:  "45*4+(2.42-3.4.)/33",
			expected:    nil,
			expectError: true,
		},
		{
			expression:  "1+(242-3y)/33",
			expected:    nil,
			expectError: true,
		},
		{
			expression:  ".1+(242-3)/33",
			expected:    nil,
			expectError: true,
		},
		{
			expression:  ".",
			expected:    nil,
			expectError: true,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.expression, func(t *testing.T) {
			res, err := calc.TokenizeInternal(testCase.expression)
			if err != nil && !testCase.expectError {
				t.Fatalf(
					"expression %s is valid but error returned: %s",
					testCase.expression,
					err.Error(),
				)

				return
			}

			if !reflect.DeepEqual(res, testCase.expected) {
				t.Fatalf("got %+v, expected %+v", res, testCase.expected)
			}
		})
	}
}
