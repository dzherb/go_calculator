package calc

import (
	"errors"
	"fmt"
	"unicode"
)

type symbolType int

const (
	digit symbolType = iota
	digitSeparator
	operator
	openingBracket
	closingBracket
)

var validSymbols = map[string]symbolType{
	"1": digit,
	"2": digit,
	"3": digit,
	"4": digit,
	"5": digit,
	"6": digit,
	"7": digit,
	"8": digit,
	"9": digit,
	"0": digit,
	".": digitSeparator,
	"(": openingBracket,
	")": closingBracket,
	"+": operator,
	"-": operator,
	"*": operator,
	"/": operator,
}

type TokenType int

const (
	start TokenType = iota
	Number
	Operator
	OpeningBracket
	ClosingBracket
)

type Token struct {
	Value     string
	TokenType TokenType
}

func tokenize(expression string) ([]Token, error) { //nolint:gocognit,funlen
	if len(expression) == 0 {
		return nil, errors.New("expression is empty")
	}

	res := make([]Token, 0)

	lastTokenContainsDigitSeparator := false
	hasWhitespaceAfterLastToken := false

	for i, char := range expression {
		if unicode.IsSpace(char) {
			hasWhitespaceAfterLastToken = true
			continue
		}

		currentSymbol := string(char)

		currSymbolType, ok := validSymbols[currentSymbol]
		if !ok {
			return nil, fmt.Errorf(
				"expression contains invalid Token at position %d: %s",
				i,
				string(char),
			)
		}

		var lastToken *Token

		if i == 0 {
			lastToken = &Token{Value: "", TokenType: start}
		} else {
			lastToken = &res[len(res)-1]
		}

		var currTokenType TokenType

		switch currSymbolType {
		case operator:
			currTokenType = Operator
		case openingBracket:
			currTokenType = OpeningBracket
		case closingBracket:
			currTokenType = ClosingBracket
		case digit, digitSeparator:
			currTokenType = Number

			if hasWhitespaceAfterLastToken && lastToken.TokenType == Number {
				return nil, fmt.Errorf(
					"unexpected whitespace at position %d",
					i,
				)
			}
		}

		hasWhitespaceAfterLastToken = false

		switch currSymbolType {
		case digit:
			if lastToken.TokenType == Number {
				lastToken.Value += currentSymbol

				continue
			}
		case digitSeparator:
			if lastToken.TokenType != Number ||
				lastTokenContainsDigitSeparator {
				return nil, fmt.Errorf("expression contains invalid Number")
			}

			lastTokenContainsDigitSeparator = true
			lastToken.Value += currentSymbol

			continue
		case operator, openingBracket, closingBracket:
			lastTokenContainsDigitSeparator = false
		}

		res = append(res, Token{
			Value:     currentSymbol,
			TokenType: currTokenType,
		})
	}

	return res, nil
}

// Проверка корректности скобочной последовательности.
func validateBrackets(tokens []Token) error {
	balance := 0

	for _, currToken := range tokens {
		if currToken.TokenType == OpeningBracket { //nolint:staticcheck
			balance++
		} else if currToken.TokenType == ClosingBracket {
			balance--

			if balance < 0 {
				return errors.New("unexpected closing bracket")
			}
		}
	}

	if balance != 0 {
		return errors.New("unexpected opening bracket")
	}

	return nil
}

// Проверка корректного расположения операторов и операндов.
func validateTokenSequence(tokens []Token) error { //nolint:gocognit
	previousType := start

	for _, currToken := range tokens {
		switch currToken.TokenType {
		case Number:
			if previousType == Number {
				return errors.New("no operator between numbers")
			}

			previousType = Number
		case OpeningBracket:
			if previousType == Number {
				return errors.New("no operator before opening bracket")
			}

			previousType = OpeningBracket
		case ClosingBracket:
			if previousType == Operator || previousType == OpeningBracket {
				return errors.New("operator before closing bracket")
			}

			previousType = ClosingBracket
		case Operator, start:
			if (previousType == Operator || previousType == start) &&
				currToken.Value != "-" {
				return errors.New("unexpected operator")
			}

			previousType = Operator
		}
	}

	if previousType == Operator {
		return errors.New("expression ends with Operator")
	}

	return nil
}

// Объединяет токен минуса с токеном числа, если это унарный минус.
func preprocessTokens(tokens []Token) []Token {
	var result []Token

	nextNumberIsNegative := false

	for i := 0; i < len(tokens); i++ {
		currToken := tokens[i]

		var isAfterOpenBr bool

		var isAfterOp bool

		if i > 0 {
			isAfterOpenBr = tokens[i-1].TokenType == OpeningBracket
			isAfterOp = tokens[i-1].TokenType == Operator
		}

		// Унарный минус, если он:
		// 1. В начале выражения
		// 2. После открывающей скобки (
		// 3. После оператора (+, -, *, /)
		if currToken.Value == "-" &&
			(i == 0 || isAfterOpenBr || isAfterOp) {
			nextNumberIsNegative = true

			continue
		}

		if currToken.TokenType == Number && nextNumberIsNegative {
			negativeNumberToken := Token{
				TokenType: Number,
				Value:     tokens[i-1].Value + currToken.Value,
			}
			result = append(result, negativeNumberToken)
		} else {
			result = append(result, currToken)
		}

		nextNumberIsNegative = false
	}

	return result
}

func Tokenize(expression string) ([]Token, error) {
	tokens, err := tokenize(expression)
	if err != nil {
		return nil, err
	}

	err = validateBrackets(tokens)
	if err != nil {
		return nil, err
	}

	err = validateTokenSequence(tokens)
	if err != nil {
		return nil, err
	}

	return preprocessTokens(tokens), nil
}
