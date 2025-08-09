package service

import (
	"Form-Mailly-Go/internal/model"
	"regexp"
	"strings"
	"unicode/utf8"
)

type Validator struct {
	Errors []string
}

func NewValidator() *Validator {
	return &Validator{Errors: make([]string, 0, 6)}
}

// Check for Required fields
func (v *Validator) Required(field, value string) {
	if strings.TrimSpace(value) == "" {
		v.Errors = append(v.Errors, field+" is required")
	}
}

// Email validator
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func (v *Validator) Email(field, value string) {
	if !emailRegex.MatchString(value) {
		v.Errors = append(v.Errors, field+" is not a valid email address")
	}
}

// URL validator
func (v *Validator) URL(field, value string) {
	if value == "" {
		return
	}

	// Fast URL validation
	if !strings.HasPrefix(value, "http://") && !strings.HasPrefix(value, "https://") {
		v.Errors = append(v.Errors, field+" must start with http:// or https://")
		return
	}

	if strings.Contains(value, " ") {
		v.Errors = append(v.Errors, field+" contains invalid spaces")
	}
}

// MaxLength validator
func (v *Validator) MaxLength(field, value string, max int) {
	if utf8.RuneCountInString(value) > max {
		v.Errors = append(v.Errors, field+" exceeds maximum length of "+string(rune(max))+" characters")
	}
}

func (v *Validator) IsValid() bool {
	return len(v.Errors) == 0
}

func ValidateFormData(form *model.ContactForm) []string {
	v := NewValidator()

	v.Required("name", form.Name)
	v.Required("email", form.Email)
	v.Required("subject", form.Subject)
	v.Required("message", form.Message)
	//v.Required("product_name", form.ProductName)
	//v.Required("product_website", form.ProductWebsite)

	v.Email("email", form.Email)
	//v.URL("product_website", form.ProductWebsite)

	//v.MaxLength("name", form.Name, 100)
	//v.MaxLength("email", form.Email, 255)
	//v.MaxLength("subject", form.Subject, 200)
	//v.MaxLength("product_name", form.ProductName, 100)

	return v.Errors
}
