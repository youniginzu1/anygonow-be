package template

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
)

var (
	frontend = os.Getenv("CONFIG_FRONTEND_URL")
)

func RegisterMailTemplate(otp, otpId string) []byte {
	var err error
	subject := "Register"
	link := fmt.Sprintf("%s/verify/active-account?otp=%s&otpId=%s", frontend, otp, otpId)
	temp := template.New("mail")
	temp, err = temp.Parse(mail)
	if err != nil {
		return []byte("error mail")
	}
	body := new(bytes.Buffer)
	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: %s \n%s\n\n", subject, mimeHeaders)))
	temp.Execute(body, struct {
		Link string
	}{
		Link: link,
	})
	return body.Bytes()
}

func ForgotMailTemplate(otp, otpId string) []byte {
	var err error
	subject := "Reset Password"
	link := fmt.Sprintf("%s/verify/reset-password?otp=%s&otpId=%s", frontend, otp, otpId)
	temp := template.New("mail")
	temp, err = temp.Parse(mail)
	if err != nil {
		return []byte("error mail")
	}
	body := new(bytes.Buffer)
	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: %s \n%s\n\n", subject, mimeHeaders)))
	temp.Execute(body, struct {
		Link string
	}{
		Link: link,
	})
	return body.Bytes()

}

var mail = `
<!DOCTYPE html
    PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">

<head>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
    <meta name="viewport" content="width=320, initial-scale=1" />
    <title>Airmail Confirm</title>
    <style type="text/css">
        .title {
            font-family: Arial;
            font-style: normal;
            font-size: 34px;
            line-height: 39px;
            margin-bottom: 50px;
            margin-top: 50px;
        }
    </style>
</head>

<body>
    <div class="title">[Anygonow]</div>
    <div style="padding:20px; margin:0; display:block; background:#F5F5F5; -webkit-text-size-adjust:none">
        <div style="width: 1153px; margin-left: 34px; margin-top: 20px; background: #F5F5F5; font-size: 18px">
            <a href={{ .Link }}>Click here</a></br>
        </div>
    </div>
</body>

</html>
`
