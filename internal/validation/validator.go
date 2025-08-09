package validation

type Validator struct {
	Error string
}

func New() *Validator {
	return &Validator{
		Error: "", // no error initially
	}
}

func (v *Validator) ValidateField(field Field) {

	// If already have an error, skip further validations
	if v.Error != "" {
		return
	}

	for _, rule := range field.Rules {
		if ok, msg := rule(field.Name, field.Value); !ok {
			v.Error = msg
			return // fail-fast: stop at first error detected for this field
		}
	}
}

func (v *Validator) IsValid() bool {
	return v.Error == ""
}
