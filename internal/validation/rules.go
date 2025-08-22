package validation

import (
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

// Rule defines the function signature for validation rules.
// Each rule receives the field name and its value, and returns:
// - a boolean indicating if the value is valid
// - an error message (used only if validation fails)
type Rule func(fieldName string, value *string) (bool, string)

// RequiredRule checks that the value is not empty (after trimming spaces).
func RequiredRule() Rule {
	return func(field string, value *string) (bool, string) {
		if strings.TrimSpace(*value) == "" {
			return false, field + " is required"
		}
		return true, ""
	}
}

// Regex to validate email format.
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// EmailRule checks that the value is a valid email address.
// Note: Allows empty values — use RequiredRule in combination to enforce presence.
func EmailRule() Rule {
	return func(field string, value *string) (bool, string) {
		if value == nil || strings.TrimSpace(*value) == "" {
			return true, "" // Considered valid if empty
		}
		if !emailRegex.MatchString(strings.TrimSpace(*value)) {
			return false, field + " is not a valid email address"
		}
		return true, ""
	}
}

// MaxLengthRule ensures the string length does not exceed a maximum number of runes.
func MaxLengthRule(max int) Rule {
	return func(field string, value *string) (bool, string) {
		if utf8.RuneCountInString(*value) > max {
			return false, field + " must be less than or equal to " + strconv.Itoa(max) + " characters"
		}
		return true, ""
	}
}

// ProductNameRule provides a default product name if the user leaves it blank.
func ProductNameRule() Rule {
	return func(field string, value *string) (bool, string) {
		if strings.TrimSpace(*value) == "" {
			*value = "FormMaillyGo" // Default product name
		}
		return true, ""
	}
}

// UrlRule validates that the input is a properly formatted URL.
// If blank, it assigns a default URL.
// It checks for:
// - valid URI format
// - http/https scheme
// - presence of host with a valid domain and TLD

var tldRegex = regexp.MustCompile(`\.(?:[A-Za-z]{2,63}|xn--[A-Za-z0-9-]{2,59})$`)

func UrlRule() Rule {
	return func(field string, value *string) (bool, string) {
		// 1. Assign default if empty
		if value == nil || strings.TrimSpace(*value) == "" {
			*value = "https://ajaynegi45.github.io/FormMaillyGo/"
			return true, ""
		}

		// 2. Parse and validate URL
		parsedURL, err := url.ParseRequestURI(strings.TrimSpace(*value))
		if err != nil {
			return false, field + " must be a valid URL"
		}

		// 3. Check for http(s)
		if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
			return false, field + " must start with http:// or https://"
		}

		// 4. Check host presence
		if parsedURL.Host == "" {
			return false, field + " must contain a valid domain"
		}

		// 5. Ensure host contains a dot (e.g., google.com)
		if !strings.Contains(parsedURL.Host, ".") {
			return false, field + " must contain a valid domain name (like example.com)"
		}

		// 6. check for valid TLD using regex (strict mode with IDN punycode support)
		// e.g., domain.tld where tld = 2 to 63 characters, alphabetic
		if !tldRegex.MatchString(parsedURL.Host) {
			return false, field + " must have a valid top-level domain (e.g., .com, .org)"
		}

		return true, ""
	}
}

// NumericRule checks that the value consists of digits only (0-9).
// It allows empty values — use RequiredRule in combination to enforce presence.
func NumericRule() Rule {
	return func(field string, value *string) (bool, string) {
		if value == nil || *value == "" {
			return true, "" // Considered valid if empty
		}
		// Match only digits
		if !regexp.MustCompile(`^\d+$`).MatchString(*value) {
			return false, field + " must contain only numeric characters"
		}
		return true, ""
	}
}
