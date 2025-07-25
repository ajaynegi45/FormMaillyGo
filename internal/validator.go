package internal

func ValidateFormData(contactForm *ContactForm) bool {
	if len(contactForm.Name) <= 0 || len(contactForm.Email) <= 0 || len(contactForm.Subject) <= 0 || len(contactForm.Message) <= 0 {
		return false
	}
	return true
}
