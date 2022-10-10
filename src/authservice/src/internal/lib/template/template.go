package template

import (
	"bytes"
	"fmt"
	"os"
	"text/template"

	"github.com/aqaurius6666/authservice/src/internal/var/c"
	"github.com/google/uuid"
)

var (
	frontend = os.Getenv("CONFIG_FRONTEND_URL")
)
var (
	TEMPLATE_VARS = map[c.OTP_TYPE]TemplateVar{
		c.OTP_TYPE_FORGOT_PASSWORD: {
			EmailSender:  "no-reply@anygonow.com",
			Text1:        "We got a request to reset your Anygonow password.",
			Text2:        "Please click this button to reset your password. This email will be expired in 5 minutes.",
			CompanyEmail: "support@anygonow.com",
			ButtonText:   "Reset password",
			Subject:      "[AnygoNow] Reset Password Email",
		},
		c.OTP_TYPE_REGISTER: {
			EmailSender:  "no-reply@anygonow.com",
			Text1:        "Congratulations, you have successfully registered",
			Text2:        "Please click this button to verify your email.",
			CompanyEmail: "support@anygonow.com",
			ButtonText:   "Verify",
			Subject:      "[AnygoNow] Verification Email",
		},
		c.OTP_TYPE_CHANGE_MAIL_AND_PASS: {
			EmailSender:  "no-reply@anygonow.com",
			Text1:        "Congratulations, you account is successfully actived by changing default email and password",
			Text2:        "Please click this button to verify your email.",
			CompanyEmail: "support@anygonow.com",
			ButtonText:   "Active",
			Subject:      "[AnygoNow] Active account",
		},
	}
)

type TemplateVar struct {
	Link         string
	Text1        string
	Text2        string
	UserEmail    string
	CompanyEmail string
	ButtonText   string
	EmailSender  string
	Subject      string
	Boundary     string
}

func (t TemplateVar) WithLink(link string) TemplateVar {
	t.Link = link
	return t
}
func (t TemplateVar) WithUserEmail(email string) TemplateVar {
	t.UserEmail = email
	return t
}

func (t TemplateVar) WithBoundary(boundary string) TemplateVar {
	t.Boundary = boundary
	return t
}

func RegisterMailTemplate(to, otp, otpId string) []byte {
	var err error
	link := fmt.Sprintf("%s/verify/active-account?otp=%s&otpId=%s", frontend, otp, otpId)
	temp := template.New("mail")
	temp, err = temp.Parse(mail)
	if err != nil {
		return []byte("error mail")
	}
	body := new(bytes.Buffer)
	err = temp.Execute(body, TEMPLATE_VARS[c.OTP_TYPE_REGISTER].WithLink(link).WithUserEmail(to).WithBoundary(uuid.NewString()))
	if err != nil {
		return []byte("error mail")
	}
	return body.Bytes()
}

func ForgotMailTemplate(to string, otp, otpId string) []byte {
	var err error
	link := fmt.Sprintf("%s/verify/reset-password?otp=%s&otpId=%s", frontend, otp, otpId)
	temp := template.New("mail")
	temp, err = temp.Parse(mail)
	if err != nil {
		return []byte("error mail")
	}
	body := new(bytes.Buffer)
	err = temp.Execute(body, TEMPLATE_VARS[c.OTP_TYPE_FORGOT_PASSWORD].WithLink(link).WithUserEmail(to).WithBoundary(uuid.NewString()))
	if err != nil {
		return []byte("error mail")
	}
	return body.Bytes()
}

func ChangeMailAndPassTemplate(to string, otp, otpId string) []byte {
	var err error
	link := fmt.Sprintf("%s/verify/active-account?otp=%s&otpId=%s", frontend, otp, otpId)
	temp := template.New("mail")
	temp, err = temp.Parse(mail)
	if err != nil {
		return []byte("error mail")
	}
	body := new(bytes.Buffer)
	err = temp.Execute(body, TEMPLATE_VARS[c.OTP_TYPE_CHANGE_MAIL_AND_PASS].WithLink(link).WithUserEmail(to).WithBoundary(uuid.NewString()))
	if err != nil {
		return []byte("error mail")
	}
	return body.Bytes()
}

var mail = `From: Anygonow <{{ .EmailSender }}>
Subject: {{ .Subject }}
To: {{ .UserEmail }}
MIME-Version: 1.0
Content-Type: multipart/alternative; boundary="{{ .Boundary }}"

--{{ .Boundary }}
Content-Type: text/plain

Hi {{ .UserEmail }},

{{ .Text1 }}
{{ .Text2 }}
{{ .Link }}

For further questions, please contact: {{ .CompanyEmail }}

--{{ .Boundary }}
Content-Type: text/html

<!DOCTYPE html>
<html>

<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<title>AnyGoNow</title>
</head>

<body style="background-color: #ffffff">
	<div style="background-color: #f6f6f6; max-width: 600px; margin: 0 auto;">
		<table style="padding-left: 20px; padding-top: 30px">
			<tr>
				<td>
					<img src="https://d26p4pe0nfrg62.cloudfront.net/email/icon70x70.png"> </img>
				</td>
			</tr>
			<tr>
				<td>
					<div>Hi <span href="mailto:{{ .UserEmail }}">{{ .UserEmail }}</span>,</div>
				</td>
			</tr>
		</table>
		<table style="padding-left: 20px; padding-top: 10px">
			<tr>
				<td>
					<div>
            {{ .Text1 }}
					</div>
				</td>
			</tr>
			<tr>
				<td>
					<div>
            {{ .Text2 }}
					</div>
				</td>
			</tr>
			<tr>
				<td style="padding-top: 10px;">
						<div style="background-color: #ff511a;
						width: 100px;
						height: 30px;
						border-radius: 4px;
						text-align: center;" >
							<a style="text-align: center;
							display: inline-block;
							vertical-align: sub;
							text-decoration: none !important;
							color: black;
							margin: 0 auto" href="{{.Link }}">{{ .ButtonText }}</a>
						</div>
					
				</td>
			</tr>
		</table>
		<table style="padding-left: 20px; padding-top: 20px; padding-bottom: 30px;">
			<tr>
				<td>
					<div>
						For further questions, please contact: <span><a href="mailto:{{ .CompanyEmail }}">{{ .CompanyEmail }}</a></span>
					</div>
				</td>
			</tr>
		</table>
	</div>

</body>

</html>
--{{ .Boundary }}--
`
