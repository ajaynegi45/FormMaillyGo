package model

type Email struct {
	SentTo      string `json:"sent_to"`
	Subject     string `json:"subject,omitempty"`
	Message     string `json:"message"`
	ProductName string `json:"product_name,omitempty"`
}

type EmailResult struct {
	Email  string `json:"email"`
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}
