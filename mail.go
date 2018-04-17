package main

import (
	"bytes"
	"html/template"
	"net/smtp"
	"os"
)

type MailData struct {
	*Payment
	*Order
}

func mail(to string, subject string, body string) {
	if os.Getenv("ENV") == "prod" {
		print("Sending mail to ", to, subject)
		from := "jairo@suncork.net"
		pass := "*************"
		msg := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n" + body

		e := smtp.SendMail("smtp.gmail.com:587",
			smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
			from, []string{to}, []byte(msg))

		if !err(e) {
			print("Email sent")
		}
	}
}

func mailTemplatePayment(payment *Payment, order *Order) string {
	print("Loading payment email template")
	tmpl, e := template.ParseFiles("views/mail/paymentData.html")

	md := MailData{
		payment,
		order,
	}

	if err(e) {
		return ""
	}

	var tpl bytes.Buffer
	e = tmpl.Execute(&tpl, md)
	if err(e) {
		return ""
	}
	print("Email template loaded succesfully")
	return tpl.String()
}
