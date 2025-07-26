package template

import "Form-Mailly-Go/internal/model"

func BuildContactFormMessage(form *model.ContactForm) string {
	return "<div style=\"font-family:Helvetica,Arial,sans-serif; font-size:16px; margin:0; color:#0b0c0c; background-color:#ffffff\">\n" +
		"<span style=\"display:none;font-size:1px;color:#fff;max-height:0\"></span>\n" +
		"<table role=\"presentation\" width=\"100%\" style=\"border-collapse:collapse;min-width:100%;width:100%!important\" cellpadding=\"0\" cellspacing=\"0\" border=\"0\">\n" +
		"  <tr>\n" +
		"    <td bgcolor=\"#0b0c0c\">\n" +
		"      <table role=\"presentation\" align=\"center\" width=\"100%\" style=\"max-width:580px; border-collapse:collapse;\" cellpadding=\"0\" cellspacing=\"0\">\n" +
		"        <tr>\n" +
		"          <td style=\"padding: 20px 10px;\">\n" +
		"            <span style=\"font-size:28px; font-weight:700; color:#ffffff;\">Contact Form Submission</span>\n" +
		"          </td>\n" +
		"        </tr>\n" +
		"      </table>\n" +
		"    </td>\n" +
		"  </tr>\n" +
		"</table>\n" +
		"<table role=\"presentation\" align=\"center\" cellpadding=\"0\" cellspacing=\"0\" border=\"0\" style=\"max-width:580px; width:100%!important;\">\n" +
		"  <tr>\n" +
		"    <td>\n" +
		"      <table width=\"100%\" style=\"border-collapse:collapse\">\n" +
		"        <tr>\n" +
		"          <td bgcolor=\"#1D70B8\" height=\"10\"></td>\n" +
		"        </tr>\n" +
		"      </table>\n" +
		"    </td>\n" +
		"  </tr>\n" +
		"</table>\n" +
		"<table role=\"presentation\" align=\"center\" cellpadding=\"0\" cellspacing=\"0\" border=\"0\" style=\"max-width:580px; width:100%!important;\">\n" +
		"  <tr><td height=\"30\"></td></tr>\n" +
		"  <tr>\n" +
		"    <td style=\"font-size:19px; line-height:1.4; color:#0b0c0c;\">\n" +
		"      <p><strong>Name:</strong> " + form.Name + "</p>\n" +
		"      <p><strong>Email:</strong> " + form.Email + "</p>\n" +
		"      <p><strong>Reason:</strong> " + form.Subject + "</p>\n" +
		"      <p><strong>Message: </strong>" + form.Message + "</p>\n" +
		"      <br>\n" +
		"      <p>Best regards,<br><strong>" + form.ProductName + "</strong><br>" + form.ProductWebsite + "</p>\n" +
		"    </td>\n" +
		"  </tr>\n" +
		"  <tr><td height=\"30\"></td></tr>\n" +
		"</table>\n" +
		"</div>"
}
