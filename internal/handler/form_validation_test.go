package handler

import (
	"Form-Mailly-Go/internal/model"
	"strings"
	"testing"
)

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
			err := validateContactForm(&tc.form)
			if err != tc.wantError {
				t.Errorf("ValidateContactForm() = %q, want %q", err, tc.wantError)
			}
		})
	}
}

func stringOfLength(n int) string {
	return strings.Repeat("a", n)
}
