## FormMaillyGo – Free, Serverless Contact Form to Email Gateway

**FormMaillyGo** is a lightweight backend that turns any contact form submission into a **real email** — **fast, secure, and completely free (within generous [AWS Lambda limits](https://aws.amazon.com/lambda/pricing/) )**.

### 🧑‍💻 Whom is this for?

* Solo developers, indie hackers, or small startups
* Static websites (like those on Netlify, Vercel, GitHub Pages)
* Anyone tired of paid contact form services (like Formspree, Getform, etc.)

---

## 🚀 Why Use FormMaillyGo?

Most free contact form services only allow **50–100 submissions/month**. That’s *very* limited.

FormMaillyGo runs on **AWS Lambda (or any serverless platform)**, which has a **much higher free tier**. It’s ideal for scaling contact form handling **at zero or very low cost**.   

You only pay if your usage goes beyond free limits — and even then, it’s cheaper than paid form services.

The AWS Lambda free tier includes one million free requests per month and 400,000 GB-seconds of compute time per month

---

## 🛠️ How It Works

1. Your frontend sends a **JSON POST** request to `/api/contact`
2. FormMaillyGo:

   * ✅ Validates input (like name, email, etc.)
   * 📧 Formats a clean HTML email
   * 📤 Sends the message via SMTP (Gmail, SES, Postmark, etc.)
3. You receive the message directly in your inbox.

---

## 🔐 Key Benefits

* **Zero dependency**: No external SDKs or vendor lock-in.
* **Military-grade validation**: Validates names, emails, message content.
* **Fully async-ready**: No delay to your main app.
* **Perfect for serverless**: Cold-start optimized.

---

## 🧪 API Endpoints

| Method | Endpoint       | Description                 |
| ------ | -------------- | --------------------------- |
| GET    | `/api/health`  | Check if the server is live |
| POST   | `/api/contact` | Send contact form data      |

### Example Contact Form Payload:

```json
{
  "name": "Alice",
  "email": "alice@email.com",
  "subject": "Product Feedback",
  "message": "Loved your product!",
  "product_name": "MySite",
  "product_website": "https://mysite.com"
}
```

---

## 🧾 Setup

### 1. Set up your `.env` file:

```
SENDER_EMAIL=your@gmail.com
SENDER_EMAIL_PASSWORD=your-gmail-app-password
RECEIVER_EMAIL=you@example.com
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
```

### 2. Run the server:

```bash
go run main.go
```

Works out-of-the-box on port **8080**.

---

## 📁 File Structure (Simplified)

| File                  | Role                                   |
| --------------------- | -------------------------------------- |
| `main.go`             | Sets up server + routes                |
| `handler.go`          | Handles requests + JSON validation     |
| `email.go`            | Sends emails using SMTP                |
| `contact.go`          | Defines the contact form structure     |
| `validator.go`        | Validates inputs like email, URL, etc. |
| `config.go`           | Loads SMTP config from `.env`          |
| `template/contact.go` | Generates HTML email template          |

---

## 🧘 Simplicity + Power

FormMaillyGo gives you **control without complexity**. You get:

* Unlimited usage (within AWS/GCP free tier)
* Full validation
* Fully customizable
* No monthly fees
