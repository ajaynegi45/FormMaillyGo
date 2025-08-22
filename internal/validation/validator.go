package validation

// Validator is used to track validation state and error message.
// Only the first encountered validation error is stored (fail-fast).
type Validator struct {
	Error string
}

// NewValidator returns a new Validator with no errors.
func NewValidator() *Validator {
	return &Validator{
		Error: "",
	}
}

// ValidateField runs all validation rules for a given field.
// It stops at the first failed rule and stores the error message.
func (v *Validator) ValidateField(field Field) {
	// Skip validation if an error already exists
	if v.Error != "" {
		return
	}

	for _, rule := range field.Rules {
		if ok, msg := rule(field.Name, field.Value); !ok {
			v.Error = msg // Record the first error encountered
			return
		}
	}
}

// IsValid returns true if no validation errors were encountered.
func (v *Validator) IsValid() bool {
	return v.Error == ""
}
