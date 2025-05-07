package auth

import (
	"fmt"
	"regexp"
)

var DefaultUsernameValidator = CompositeValidator(
	MinLengthValidator(4), //nolint:mnd
	ContainsOnlyValidator(`a-zA-Z0-9`),
)

var DefaultPasswordValidator = CompositeValidator(
	MinLengthValidator(8), //nolint:mnd
	ContainsOnlyValidator(`a-zA-Z0-9!@#$%^&*()_`),
)

type Validator interface {
	Validate(value string) error
}

type ValidatorFunc func(value string) error

func (f ValidatorFunc) Validate(value string) error {
	return f(value)
}

func CompositeValidator(
	validators ...Validator,
) Validator {
	return ValidatorFunc(func(value string) error {
		for _, pv := range validators {
			if err := pv.Validate(value); err != nil {
				return err
			}
		}

		return nil
	})
}

var MinLengthValidator = func(length int) Validator {
	return ValidatorFunc(func(value string) error {
		if len(value) < length {
			return fmt.Errorf("must be at least %d characters", length)
		}

		return nil
	})
}

var ContainsOnlyValidator = func(allowed string) Validator {
	return ValidatorFunc(func(value string) error {
		re := regexp.MustCompile(`^[` + allowed + `]+$`)
		if !re.MatchString(value) {
			return fmt.Errorf("must contain only: %s", allowed)
		}
		return nil
	})
}
