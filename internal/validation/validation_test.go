package validation

import (
	"testing"
)

func TestRequiredRule(t *testing.T) {
	rule := RequiredRule()

	cases := map[string]struct {
		input    string
		expected bool
	}{
		"Empty string":  {"", false},
		"Spaces only":   {"   ", false},
		"Valid string":  {"hello", true},
		"Unicode input": {"こんにちは", true},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			valid, _ := rule("field", strPtr(tc.input))
			if valid != tc.expected {
				t.Errorf("Expected %v, got %v", tc.expected, valid)
			}
		})
	}
}

func TestEmailRule(t *testing.T) {
	rule := EmailRule()

	cases := map[string]struct {
		input    string
		expected bool
	}{
		"Missing @":       {"email.com", false},
		"Missing domain":  {"test@", false},
		"Invalid chars":   {"test@#%.com", false},
		"Valid":           {"test@example.com", true},
		"Valid with +":    {"john.doe+123@gmail.com", true},
		"Unicode address": {"ユーザー@例.jp", false},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			valid, _ := rule("email", strPtr(tc.input))
			if valid != tc.expected {
				t.Errorf("Expected %v, got %v", tc.expected, valid)
			}
		})
	}
}

func TestMaxLengthRule(t *testing.T) {
	rule := MaxLengthRule(5)

	cases := map[string]struct {
		input    string
		expected bool
	}{
		"Under limit":        {"123", true},
		"At limit":           {"12345", true},
		"Over limit":         {"123456", false},
		"Empty":              {"", true},
		"Multibyte over":     {"你好世界abc", false}, // 6 runes
		"Multibyte at limit": {"你好世", true},      // 3 runes
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			valid, _ := rule("field", strPtr(tc.input))
			if valid != tc.expected {
				t.Errorf("Expected %v, got %v", tc.expected, valid)
			}
		})
	}
}

func TestValidator_ValidateFieldAndIsValid(t *testing.T) {
	t.Run("Invalid input adds error", func(t *testing.T) {
		validator := NewValidator()
		validator.ValidateField(Field{
			Name:  "email",
			Value: strPtr("bademail"),
			Rules: []Rule{RequiredRule(), EmailRule()},
		})

		if validator.IsValid() {
			t.Error("Expected validator to be invalid")
		}
		if validator.Error == "" {
			t.Error("Expected validator to have an error")
		}
	})

	t.Run("Valid input has no errors", func(t *testing.T) {
		validator := NewValidator()
		validator.ValidateField(Field{
			Name:  "email",
			Value: strPtr("user@example.com"),
			Rules: []Rule{RequiredRule(), EmailRule()},
		})

		if !validator.IsValid() {
			t.Error("Expected validator to be valid")
		}
	})
}

func strPtr(s string) *string {
	return &s
}
