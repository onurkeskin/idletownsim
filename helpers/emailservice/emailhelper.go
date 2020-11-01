package emailhelper

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/mail"
	"net/smtp"
	"strings"
)

var (
	auth smtp.Auth
)

func init() {
	auth = smtp.PlainAuth(
		"",
		"yoursmtp....",
		"yoursmtp....",
		"yoursmtp....",
	)
}

func SendMail(toStr string, titleStr string, msgStr string) {

	from := mail.Address{"", "onurabi@onurabi.com"}
	to := mail.Address{"", toStr}
	title := titleStr

	header := make(map[string]string)
	header["From"] = from.String()
	header["To"] = to.String()
	header["Subject"] = encodeRFC2047(title)
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/plain; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(msgStr))

	err := smtp.SendMail(
		"smtp.mailgun.org:2525",
		auth,
		"deneme@example.org",
		[]string{to.Address},
		[]byte(message),
	)
	if err != nil {
		log.Fatal(err)
	}
}

func encodeRFC2047(String string) string {
	// use mail's rfc2047 to encode any string
	addr := mail.Address{String, ""}
	return strings.Trim(addr.String(), " <>")
}
