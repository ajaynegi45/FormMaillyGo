package template

import "Form-Mailly-Go/internal/model"

func BuildContactFormMessage1(form *model.ContactForm) string {
	return `<div style=font-family:Helvetica,Arial,sans-serif; font-size:16px; margin:0; color:#0b0c0c; background-color:#ffffff> 
		<span style=display:none;font-size:1px;color:#fff;max-height:0></span> 
		<table role=presentation width=100% style=border-collapse:collapse;min-width:100%;width:100%!important cellpadding=0 cellspacing=0 border=0> 
		  <tr> 
		    <td bgcolor=#0b0c0c> 
		      <table role=presentation align=center width=100% style=max-width:580px; border-collapse:collapse; cellpadding=0 cellspacing=0> 
		        <tr> 
		          <td style=padding: 20px 10px;> 
		            <span style=font-size:28px; font-weight:700; color:#ffffff;>Contact Form Submission</span> 
		          </td> 
		        </tr> 
		      </table> 
		    </td> 
		  </tr> 
		</table> 
		<table role=presentation align=center cellpadding=0 cellspacing=0 border=0 style=max-width:580px; width:100%!important;> 
		  <tr> 
		    <td> 
		      <table width=100% style=border-collapse:collapse> 
		        <tr> 
		          <td bgcolor=#1D70B8 height=10></td> 
		        </tr> 
		      </table> 
		    </td> 
		  </tr> 
		</table> 
		<table role=presentation align=center cellpadding=0 cellspacing=0 border=0 style=max-width:580px; width:100%!important;> 
		  <tr><td height=30></td></tr> 
		  <tr> 
		    <td style=font-size:19px; line-height:1.4; color:#0b0c0c;> 
		      <p><strong>Name:</strong>` + form.Name + `</p> 
		      <p><strong>Email:</strong>` + form.Email + `</p> 
		      <p><strong>Reason:</strong>` + form.Subject + `</p> 
		      <p><strong>Message: </strong>` + form.Message + `</p> 
		      <br> 
		      <p>Best regards,<br><strong>` + form.ProductName + `</strong><br>` + form.ProductWebsite + `</p> 
		    </td> 
		  </tr> 
		  <tr><td height=30></td></tr> 
		</table> 
		</div>`
}
