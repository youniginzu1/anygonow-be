package template

import "fmt"

func GetMail(from string, to string, subject string, message string) []byte {
	mail := fmt.Sprintf(`
	To: Vu Nguyen <%s>
	Subject: %s
	%s`, to, subject, message)
	fmt.Printf("mail: %v\n", mail)
	return []byte(mail)
}
