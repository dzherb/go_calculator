package calculator

import (
	"errors"
	"fmt"
	"unicode"
)

type symbolMeta struct {
	isDigit          bool
	isDigitSeparator bool
	isOperator       bool
	isOpeningBracket bool
	isClosingBracket bool
}

var validTokens = map[string]symbolMeta{
	"1": {isDigit: true},
	"2": {isDigit: true},
	"3": {isDigit: true},
	"4": {isDigit: true},
	"5": {isDigit: true},
	"6": {isDigit: true},
	"7": {isDigit: true},
	"8": {isDigit: true},
	"9": {isDigit: true},
	"0": {isDigit: true},
	".": {isDigitSeparator: true},
	"(": {isOpeningBracket: true},
	")": {isClosingBracket: true},
	"+": {isOperator: true},
	"-": {isOperator: true},
	"*": {isOperator: true},
	"/": {isOperator: true},
}

type tokenType int

const (
	start tokenType = iota
	number
	operator
	openingBracket
	closingBracket
)

type token struct {
	value     string
	tokenType tokenType
}

func tokenize(expression string) ([]token, error) {
	if len(expression) == 0 {
		return nil, errors.New("expression is empty")
	}

	res := make([]token, 0)

	lastTokenContainsDigitSeparator := false
	for i, char := range expression {
		if unicode.IsSpace(char) {
			continue
		}
		currentSymbol := string(char)

		currentSymbolMeta, ok := validTokens[currentSymbol]
		if !ok {
			return nil, fmt.Errorf("expression contains invalid token at position %d: %s", i+1, string(char))
		}

		currentTokenType := number

		if currentSymbolMeta.isOperator {
			currentTokenType = operator
		} else if currentSymbolMeta.isOpeningBracket {
			currentTokenType = openingBracket
		} else if currentSymbolMeta.isClosingBracket {
			currentTokenType = closingBracket
		}

		var lastToken *token
		if i == 0 {
			lastToken = &token{value: "", tokenType: start}
		} else {
			lastToken = &res[len(res)-1]
		}

		if currentSymbolMeta.isDigit {
			if lastToken.tokenType == number {
				lastToken.value = lastToken.value + currentSymbol
				continue
			}
		} else if currentSymbolMeta.isDigitSeparator {
			if !(lastToken.tokenType == number) || lastTokenContainsDigitSeparator {
				return nil, fmt.Errorf("expression contains invalid number")
			}
			lastTokenContainsDigitSeparator = true
			lastToken.value = lastToken.value + currentSymbol
			continue
		} else {
			lastTokenContainsDigitSeparator = false
		}

		res = append(res, token{
			value:     currentSymbol,
			tokenType: currentTokenType,
		})
	}
	return res, nil
}

// Проверка корректности скобочной последовательности
func validateBrackets(tokens []token) error {
	balance := 0
	for _, currToken := range tokens {
		if currToken.tokenType == openingBracket {
			balance++
		} else if currToken.tokenType == closingBracket {
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

// Проверка корректного расположения операторов и операндов
func validateTokenSequence(tokens []token) error {
	previousType := start

	for _, currToken := range tokens {
		if currToken.tokenType == number {
			if previousType == number {
				return errors.New("no operator between numbers")
			}
			previousType = number
		} else if currToken.tokenType == openingBracket {
			if previousType == number {
				return errors.New("no operator before opening bracket")
			}
			previousType = openingBracket
		} else if currToken.tokenType == closingBracket {
			if previousType == operator || previousType == openingBracket {
				return errors.New("operator before closing bracket")
			}
			previousType = closingBracket
		} else {
			if (previousType == operator || previousType == start) && currToken.value != "-" {
				return errors.New("unexpected operator")
			}
			previousType = operator
		}
	}

	if previousType == operator {
		return errors.New("expression ends with operator")
	}
	return nil
}

// Приводит выражения вида -12/-2 к (0-12)/(0-2)
func preprocessTokens(tokens []token) []token {
	var result []token

	for i := 0; i < len(tokens); i++ {
		currToken := tokens[i]

		// Унарный минус, если он:
		// 1. В начале выражения
		// 2. После открывающей скобки (
		// 3. После оператора (+, -, *, /)
		if currToken.value == "-" && (i == 0 || tokens[i-1].tokenType == openingBracket || tokens[i-1].tokenType == operator) {
			result = append(result, token{value: "(", tokenType: openingBracket})
			result = append(result, token{value: "0", tokenType: number})
			result = append(result, token{value: "-", tokenType: operator})
			continue
		}

		result = append(result, currToken)

		// Закрываем скобку, если это число после унарного минуса
		if len(result) > 2 && result[len(result)-2].value == "-" && currToken.tokenType == number {
			result = append(result, token{value: ")", tokenType: closingBracket})
		}
	}

	return result
}

func Tokenize(expression string) ([]token, error) {
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

	return preprocessTokens(tokens)
}
