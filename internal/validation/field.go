package validation

// Field represents a single form field to be validated.
// It includes:
// - Name: the identifier of the field (used in error messages).
// - Value: the actual data provided for the field (can be nil).
// - Rules: a slice of Rule functions that will be applied to validate the field.
type Field struct {
	Name  string
	Value *string
	Rules []Rule // Rules to apply for validation
}
