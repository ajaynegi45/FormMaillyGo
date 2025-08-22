package template

import "Form-Mailly-Go/internal/model"

func BuildContactFormMessage2(form *model.ContactForm) string {
	return `<!DOCTYPE html>
<html lang="en">
<head><meta charset="UTF-8"><title>Contact Received</title></head>
<body style="margin:0;padding:0;background-color:#e6ecf0;font-family:Arial,sans-serif;">
  <table width="100%" cellpadding="0" cellspacing="0">
    <tr>
      <td align="center" style="padding:40px 0;">
        <table width="600" cellpadding="0" cellspacing="0" style="background:#ffffff;border-radius:10px;box-shadow:0 4px 12px rgba(0,0,0,0.1);overflow:hidden;">
          
          <!-- Header Section -->
          <tr>
            <td style="background:#393E46;padding:24px 32px;color:#ffffff;text-align:left;">
              <h2 style="margin:0;font-size:22px;">New Contact Request</h2>
              <p style="margin:4px 0 0;font-size:13px;opacity:0.8;">` + GetCurrentFormattedTime() + `</p>
            </td>
          </tr>

          <!-- Details Section -->
          <tr>
            <td style="padding:32px;">
              <table width="100%" cellpadding="0" cellspacing="0" style="font-size:15px;line-height:1.6;color:#333;">
                <tr><td style="padding:8px 0;"><strong>👤 Name:</strong></td><td>` + form.Name + `</td></tr>
                <tr><td style="padding:8px 0;"><strong>📧 Email:</strong></td><td><a href="mailto:` + form.Email + `" style="color:#2563eb;text-decoration:none;">` + form.Email + `</a></td></tr>
                <tr><td style="padding:8px 0;"><strong>🎯 Subject:</strong></td><td>` + form.Subject + `</td></tr>
                <tr><td style="padding:8px 0;" colspan="2"><strong>💬 Message:</strong><br>` + form.Message + `</td></tr>
              </table>
            </td>
          </tr>

          <!-- Footer -->
          <tr>
            <td style="background:#f5f5f5;text-align:center;padding:12px;color:#888;font-size:12px;">
              Sent securely via <strong><a href="` + form.ProductWebsite + `" style="color:#2563eb;text-decoration:none;">` + form.ProductName + `</a></strong>
            </td>
          </tr>

        </table>
      </td>
    </tr>
  </table>
</body>
</html>`
}
