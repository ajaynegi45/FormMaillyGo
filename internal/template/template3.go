package template

import "Form-Mailly-Go/internal/model"

func BuildContactFormMessage3(form *model.ContactForm) string {
	return `
<table width="100%" cellpadding="0" cellspacing="0" style="background-color:#f8fafc; padding:20px;">
  <tr>
    <td align="center">
      <table width="600" style="background-color:#ffffff; border-radius:10px; box-shadow:0 4px 6px -1px rgba(0,0,0,0.1); overflow:hidden;">
        <!-- Header -->
        <tr>
          <td style="background:linear-gradient(135deg, #667eea 0%, #764ba2 100%); padding:30px; text-align:center;">
            <h1 style="color:#ffffff; font-size:24px; margin:0;">📨 New Contact Form Submission</h1>
          </td>
        </tr>

        <!-- Body -->
        <tr>
          <td style="padding:30px; font-family:Segoe UI, sans-serif; color:#0f172a;">
            <p style="font-size:16px;">Hi there 👋🏻,</p>
            <p style="font-size:15px; line-height:1.6; color:#64748b;">You've received a new message from your contact form:</p>
            <table width="100%" cellpadding="8" cellspacing="0" style="margin-top:20px; font-size:15px;">
              <tr><td width="100" style="font-weight:600;">Name:</td><td>` + form.Name + `</td></tr>
              <tr><td style="font-weight:600;">Email:</td><td><a href="mailto:` + form.Email + `" style="color:#2563eb;text-decoration:none;">` + form.Email + `</a></td></tr>
              <tr><td style="font-weight:600;">Subject:</td><td>` + form.Subject + `</td></tr>
              <tr>
                <td style="font-weight:600;">Message:</td>
                <td><div>` + form.Message + `</div></td>
              </tr>
            </table>
          </td>
        </tr>

 		<!-- Footer -->
        <tr>
          <td style="text-align:center; padding:20px; background-color:#f1f5f9; font-size:13px; color:#64748b;">
            <p style="font-size:14px; color:#94a3b8;">This message was submitted via <strong><a href="` + form.ProductWebsite + `" style="color:#2563eb;text-decoration:none;">` + form.ProductName + `</a></strong>.</p>
          </td>
        </tr>
      </table>
    </td>
  </tr>
</table>`
}
