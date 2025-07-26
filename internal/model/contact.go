package model

type ContactForm struct {
	Name           string `json:"name" validate:"required,min=1"`
	Email          string `json:"email" validate:"required,email"`
	Subject        string `json:"subject" validate:"required,min=1"`
	Message        string `json:"message" validate:"required,min=1"`
	ProductName    string `json:"product_name" validate:"min=1"`
	ProductWebsite string `json:"product_website" validate:"url"`
}
