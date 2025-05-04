package calculator

import (
	"reflect"
	"testing"
)

func TestTokenizer(t *testing.T) {
	testCases := []struct {
		expression  string
		expected    []token
		expectError bool
	}{
		{
			expression:  "1",
			expected:    []token{{value: "1", tokenType: number}},
			expectError: false,
		},
		{
			expression: "1+2-3",
			expected: []token{
				{value: "1", tokenType: number},
				{value: "+", tokenType: operator},
				{value: "2", tokenType: number},
				{value: "-", tokenType: operator},
				{value: "3", tokenType: number},
			},
			expectError: false,
		},
		{
			expression: "1+2-3.4",
			expected: []token{
				{value: "1", tokenType: number},
				{value: "+", tokenType: operator},
				{value: "2", tokenType: number},
				{value: "-", tokenType: operator},
				{value: "3.4", tokenType: number},
			},
			expectError: false,
		},
		{
			expression: "1+(242-3.4)/33",
			expected: []token{
				{value: "1", tokenType: number},
				{value: "+", tokenType: operator},
				{value: "(", tokenType: openingBracket},
				{value: "242", tokenType: number},
				{value: "-", tokenType: operator},
				{value: "3.4", tokenType: number},
				{value: ")", tokenType: closingBracket},
				{value: "/", tokenType: operator},
				{value: "33", tokenType: number},
			},
			expectError: false,
		},
		{
			expression: "45*4+(2.42-3.4)/33",
			expected: []token{
				{value: "45", tokenType: number},
				{value: "*", tokenType: operator},
				{value: "4", tokenType: number},
				{value: "+", tokenType: operator},
				{value: "(", tokenType: openingBracket},
				{value: "2.42", tokenType: number},
				{value: "-", tokenType: operator},
				{value: "3.4", tokenType: number},
				{value: ")", tokenType: closingBracket},
				{value: "/", tokenType: operator},
				{value: "33", tokenType: number},
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
			res, err := tokenize(testCase.expression)
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
