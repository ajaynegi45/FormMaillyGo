package validation

// Field represents a form field to validate
type Field struct {
	Name  string
	Value *string
	Rules []Rule
}
