package service

import (
	"Form-Mailly-Go/internal/model"
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
func (v *Validator) Email(field, value string) {

	// Check for whitespace characters
	if strings.ContainsAny(value, " \t\r\n") {
		v.Errors = append(v.Errors, field+" cannot contain whitespace")
		return
	}

	// Must contain exactly one '@'
	if atCount := strings.Count(value, "@"); atCount != 1 {
		v.Errors = append(v.Errors, field+" must contain exactly one '@' character")
		return
	}

	// Split into local and domain parts
	parts := strings.Split(value, "@")
	local, domain := parts[0], parts[1]

	// Validate local part (before @)
	if len(local) == 0 {
		v.Errors = append(v.Errors, field+" local part cannot be empty")
		return
	}
	if len(local) > 64 {
		v.Errors = append(v.Errors, field+" local part exceeds 64 characters")
		return
	}
	if local[0] == '.' || local[len(local)-1] == '.' {
		v.Errors = append(v.Errors, field+" local part cannot start or end with a dot")
		return
	}
	if strings.Contains(local, "..") {
		v.Errors = append(v.Errors, field+" local part cannot contain consecutive dots")
		return
	}

	// Validate domain part (after @)
	if len(domain) == 0 {
		v.Errors = append(v.Errors, field+" domain part cannot be empty")
		return
	}
	if len(domain) > 255 {
		v.Errors = append(v.Errors, field+" domain part exceeds 255 characters")
		return
	}
	if domain[0] == '.' || domain[len(domain)-1] == '.' {
		v.Errors = append(v.Errors, field+" domain cannot start or end with a dot")
		return
	}
	if strings.Contains(domain, "..") {
		v.Errors = append(v.Errors, field+" domain cannot contain consecutive dots")
		return
	}

	// Domain must contain at least one dot
	if !strings.Contains(domain, ".") {
		v.Errors = append(v.Errors, field+" domain must contain a dot")
		return
	}

	// Validate TLD (last part after dot)
	tldParts := strings.Split(domain, ".")
	tld := tldParts[len(tldParts)-1]
	if len(tld) < 2 {
		v.Errors = append(v.Errors, field+" TLD must be at least 2 characters")
		return
	}
	if strings.ContainsAny(tld, "0123456789") {
		v.Errors = append(v.Errors, field+" TLD cannot contain numbers")
		return
	}

	// Check for invalid characters in domain
	for _, char := range domain {
		if !isValidDomainChar(char) {
			v.Errors = append(v.Errors, field+" contains invalid domain character")
			return
		}
	}
}

func isValidDomainChar(c rune) bool {
	// Allow letters, digits, hyphens, and dots
	return ('a' <= c && c <= 'z') ||
		('A' <= c && c <= 'Z') ||
		('0' <= c && c <= '9') ||
		c == '-' || c == '.'
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
