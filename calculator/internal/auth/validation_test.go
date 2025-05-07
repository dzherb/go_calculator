package auth_test

import (
	"testing"

	"github.com/dzherb/go_calculator/internal/auth"
)

func TestMinLengthValidator(t *testing.T) {
	v := auth.MinLengthValidator(5)

	tests := []struct {
		input  string
		expect bool
	}{
		{"abcde", true},
		{"abcd", false},
		{"", false},
	}

	for _, tt := range tests {
		err := v.Validate(tt.input)
		if (err == nil) != tt.expect {
			t.Errorf(
				"MinLengthValidator.Validate(%q) = %v, expected valid: %v",
				tt.input,
				err,
				tt.expect,
			)
		}
	}
}

func TestContainsOnlyValidator(t *testing.T) {
	v := auth.ContainsOnlyValidator("a-zA-Z")

	tests := []struct {
		input  string
		expect bool
	}{
		{"Hello", true},
		{"Hello123", false},
		{"", false},
		{"abcXYZ", true},
		{"abc_xyz", false},
	}

	for _, tt := range tests {
		err := v.Validate(tt.input)
		if (err == nil) != tt.expect {
			t.Errorf(
				"ContainsOnlyValidator.Validate(%q) = %v, expected valid: %v",
				tt.input,
				err,
				tt.expect,
			)
		}
	}
}

func TestCompositeValidator(t *testing.T) {
	v := auth.CompositeValidator(
		auth.MinLengthValidator(3),
		auth.ContainsOnlyValidator("a-z"),
	)

	tests := []struct {
		input  string
		expect bool
	}{
		{"abc", true},
		{"ab", false},    // too short
		{"abcd1", false}, // contains invalid character
		{"", false},
	}

	for _, tt := range tests {
		err := v.Validate(tt.input)
		if (err == nil) != tt.expect {
			t.Errorf(
				"CompositeValidator.Validate(%q) = %v, expected valid: %v",
				tt.input,
				err,
				tt.expect,
			)
		}
	}
}

func TestDefaultUsernameValidator(t *testing.T) {
	tests := []struct {
		username string
		valid    bool
	}{
		{"user1", true},
		{"u1", false},     // too short
		{"user_1", false}, // invalid character
		{"ValidUser123", true},
	}

	for _, tt := range tests {
		err := auth.DefaultUsernameValidator.Validate(tt.username)
		if (err == nil) != tt.valid {
			t.Errorf(
				"DefaultUsernameValidator.Validate(%q) = %v, expected valid: %v",
				tt.username,
				err,
				tt.valid,
			)
		}
	}
}

func TestDefaultPasswordValidator(t *testing.T) {
	tests := []struct {
		password string
		valid    bool
	}{
		{"password123!", true},
		{"short1!", false},      // too short
		{"invalid/pass", false}, // contains `/`
		{"GoodPassword_1", true},
	}

	for _, tt := range tests {
		err := auth.DefaultPasswordValidator.Validate(tt.password)
		if (err == nil) != tt.valid {
			t.Errorf(
				"DefaultPasswordValidator.Validate(%q) = %v, expected valid: %v",
				tt.password,
				err,
				tt.valid,
			)
		}
	}
}
