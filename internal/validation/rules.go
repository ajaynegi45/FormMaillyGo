package validation

import (
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

// Rule is a function that validates a field value.
// Returns (isValid, errorMessage)
type Rule func(fieldName string, value *string) (bool, string)

func RequiredRule() Rule {
	return func(field string, value *string) (bool, string) {
		if strings.TrimSpace(*value) == "" {
			return false, field + " is required"
		}
		return true, ""
	}
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func EmailRule() Rule {
	return func(field string, value *string) (bool, string) {
		if !emailRegex.MatchString(*value) {
			return false, field + " is not a valid email address"
		}
		return true, ""
	}
}

func MaxLengthRule(max int) Rule {
	return func(field string, value *string) (bool, string) {
		if utf8.RuneCountInString(*value) > max {
			return false, field + " must be less than or equal to " + strconv.Itoa(max) + " characters"
		}
		return true, ""
	}
}

func ProductNameRule() Rule {
	return func(field string, value *string) (bool, string) {
		if strings.TrimSpace(*value) == "" {
			*value = "FormMaillyGo"
			return true, "" // This means user doesn't give url
		}
		return true, ""
	}
}

func UrlRule() Rule {
	return func(field string, value *string) (bool, string) {
		// 1. Default value
		if strings.TrimSpace(*value) == "" {
			*value = "https://ajaynegi45.github.io/FormMaillyGo/"
			return true, ""
		}

		// 2. Parse and validate URL
		parsedURL, err := url.ParseRequestURI(*value)
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

		// 6. check for valid TLD using regex (robot strict mode)
		// e.g., domain.tld where tld = 2 to 63 characters, alphabetic
		tldRegex := regexp.MustCompile(`\.[a-zA-Z]{2,63}$`)
		if !tldRegex.MatchString(parsedURL.Host) {
			return false, field + " must have a valid top-level domain (e.g., .com, .org)"
		}

		return true, ""
	}
}
