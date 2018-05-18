package models

import (
	"strings"
	"encoding/base64"
	"net/smtp"
	"fmt"
	"net/mail"
)

const (
	SMTP_SERVER = "smtp.zoho.com"
)

type HusMail struct {
	User		string
	Password	string
}

func encodeRFC2047(str string) string{
	addr := mail.Address{str, ""}
	return strings.Trim(addr.String(), " <@>")
}

func (sender HusMail) SendMail(dest []string, subject, body_message string) error {
	header := make(map[string]string)
	header["From"] = sender.User
	header["To"] = strings.Join(dest, ",")
	header["Subject"] = encodeRFC2047(subject)
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/html; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte("<html><body>"+body_message+"</body></html>"))

	err := smtp.SendMail(SMTP_SERVER + ":587",
		smtp.PlainAuth("", sender.User, sender.Password, SMTP_SERVER),
		sender.User, dest, []byte(message))
	return err
}
