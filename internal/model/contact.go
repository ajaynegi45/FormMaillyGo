package model

type ContactForm struct {
	Name           string `json:"name"`
	Email          string `json:"email"`
	Subject        string `json:"subject"`
	Message        string `json:"message"`
	ProductName    string `json:"product_name,omitempty"`
	ProductWebsite string `json:"product_website,omitempty"`
}
