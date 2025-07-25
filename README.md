# Form-Mailly-Go

Form-Mailly-Go is a lightweight, server-side contact form API written in Go. It handles simple POST requests from a frontend contact form, validates the incoming data, and sends the form contents via email using SMTP. The project is modular and easy to set up, making it a great starter solution for product/project websites that need a basic inquiry, feedback, or support contact form backend.

## Features

- **RESTful API** with endpoints for health check and contact form submission.
- **Validates** incoming form data before processing.
- **Automatic email notification**: Sends submitted contact form information as an email to a configured address.
- **Configuration via `.env`** file for SMTP/Email credentials.
- Clean and modular structure (using Go packages).
- Simple HTML email formatting for inquiries.
- Easy to integrate with any frontend or static site.

## Endpoints

| Method | Route            | Description                 |
|--------|------------------|-----------------------------|
| GET    | `/api/health`    | Health check endpoint       |
| POST   | `/api/contact`   | Submit contact form (JSON)  |

## Sample `ContactForm` JSON

```json
{
  "name": "Your Name",
  "email": "you@email.com",
  "subject": "Contact Subject",
  "message": "Your message here...",
  "product_name": "ProductName",
  "product_website": "https://yourproductsite.com"
}
```

## How it Works

1. The form POSTs data as JSON to `/api/contact`.
2. The backend validates the required fields: name, email, subject, message.
3. If valid, it formats and sends the contents as an HTML email using the SMTP settings from the `.env` file.
4. If invalid, it responds with an HTTP 400 error.

## Environment Variables (`.env` required)

```
SENDER_EMAIL=your@email.com
SENDER_EMAIL_PASSWORD=yourpassword
RECEIVER_EMAIL=destination@email.com
SMTP_HOST=smtp.example.com
SMTP_PORT=587
```

## Getting Started

1. Fill your SMTP info in `.env`.
2. Install dependencies (Go modules, [github.com/joho/godotenv](https://github.com/joho/godotenv)).
3. Run:

```bash
go run main.go
```

The server will start on port 8080.

## File Structure

- `main.go`: HTTP server, routing, endpoint logic.
- `internal/formdata.go`: ContactForm structure.
- `internal/sendemail.go`: Email sending logic.
- `internal/validator.go`: Basic validation logic.
- `internal/template/contact.go`: Builds a simple HTML email template.
