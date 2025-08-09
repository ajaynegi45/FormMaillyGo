package validation

import (
	"Form-Mailly-Go/internal/model"
	"strings"
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
		v := New()
		v.ValidateField(Field{
			Name:  "email",
			Value: strPtr("bademail"),
			Rules: []Rule{RequiredRule(), EmailRule()},
		})

		if v.IsValid() {
			t.Error("Expected validator to be invalid")
		}
		if v.Error == "" {
			t.Error("Expected validator to have an error")
		}
	})

	t.Run("Valid input has no errors", func(t *testing.T) {
		v := New()
		v.ValidateField(Field{
			Name:  "email",
			Value: strPtr("user@example.com"),
			Rules: []Rule{RequiredRule(), EmailRule()},
		})

		if !v.IsValid() {
			t.Error("Expected validator to be valid")
		}
	})
}

func TestValidateContactForm(t *testing.T) {
	cases := map[string]struct {
		form      model.ContactForm
		wantError string
	}{
		"All valid": {
			form: model.ContactForm{
				Name:    "John Doe",
				Email:   "john@example.com",
				Subject: "Hello",
				Message: "Valid message",
			},
			wantError: "",
		},
		"All fields empty": {
			form:      model.ContactForm{},
			wantError: "name is required",
		},
		"Invalid email": {
			form: model.ContactForm{
				Name:    "Alice",
				Email:   "invalid-email",
				Subject: "Hi",
				Message: "Short message",
			},
			wantError: "email is not a valid email address",
		},
		"Exceeds max length": {
			form: model.ContactForm{
				Name:    stringOfLength(101),
				Email:   "john@example.com",
				Subject: stringOfLength(301),
				Message: "Valid",
			},
			wantError: "name must be less than or equal to 100 characters",
		},
		"Whitespace + bad email + long subject + empty message": {
			form: model.ContactForm{
				Name:    "   ",
				Email:   "bad",
				Subject: stringOfLength(250),
				Message: "",
			},
			wantError: "name is required", // first error stops validation here
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			err := ValidateContactForm(&tc.form)
			if err != tc.wantError {
				t.Errorf("ValidateContactForm() = %q, want %q", err, tc.wantError)
			}
		})
	}
}

func stringOfLength(n int) string {
	return strings.Repeat("a", n)
}

func strPtr(s string) *string {
	return &s
}
