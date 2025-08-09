package validation

import (
	"Form-Mailly-Go/internal/model"
)

func ValidateContactForm(form *model.ContactForm) string {
	v := New()

	fields := []Field{
		{
			Name:  "name",
			Value: &form.Name,
			Rules: []Rule{
				RequiredRule(),
				MaxLengthRule(100),
			},
		},
		{
			Name:  "email",
			Value: &form.Email,
			Rules: []Rule{
				RequiredRule(),
				EmailRule(),
				MaxLengthRule(255),
			},
		},
		{
			Name:  "subject",
			Value: &form.Subject,
			Rules: []Rule{
				RequiredRule(),
				MaxLengthRule(300),
			},
		},
		{
			Name:  "message",
			Value: &form.Message,
			Rules: []Rule{
				RequiredRule(),
			},
		},
		{
			Name:  "product_name",
			Value: &form.ProductName,
			Rules: []Rule{
				ProductNameRule(),
			},
		},
		{
			Name:  "product_website",
			Value: &form.ProductWebsite,
			Rules: []Rule{
				UrlRule(),
			},
		},
	}

	for _, field := range fields {
		v.ValidateField(field)
		if !v.IsValid() {
			// Return immediately once an error occurs
			return v.Error
		}
	}

	return "" // no error found, valid form
}
