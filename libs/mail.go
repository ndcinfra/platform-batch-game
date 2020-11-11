package libs

import (
	"os"
	"strconv"

	"github.com/astaxie/beego/logs"
	"github.com/joho/godotenv"
	"gopkg.in/mail.v2"
	gomail "gopkg.in/mail.v2"
)

func SendEmail(counts, start, end, elasped string) {
	err := godotenv.Load()
	if err != nil {
		logs.Error("Error loading .env file")
	}

	SMTP := os.Getenv("SMTP")
	SMTPPORT, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	SMTPID := os.Getenv("SMTPID")
	SMTPPASS := os.Getenv("SMTPPASS")
	ENV := os.Getenv("ENV")
	COUNTRY := os.Getenv("COUNTRY")

	subject := "[" + COUNTRY + "-CRONJOB] " + ENV + " Daily Get Game Data Cronjob Result"

	m := gomail.NewMessage()
	m.SetHeader("From", "no-reply@naddicjapan.com")
	m.SetHeader("To", "youngtip@gmail.com", "youngtip@naddic.com", "sjwon@naddic.com")
	m.SetHeader("Subject", subject)

	body := "<br/>--------------------------------" +
		"<br/>Result Counts: " + counts + ", success.<br/><br/>" +
		"<br/>Start Time: " + start + " <br/>" +
		"<br/>End Time: " + end + " <br/>" +
		"<br/>Elasped Time: " + elasped + " <br/>" +
		"<br/>--------------------------------"

	m.SetBody("text/html", body)

	d := gomail.NewDialer(SMTP, SMTPPORT, SMTPID, SMTPPASS)
	d.StartTLSPolicy = mail.MandatoryStartTLS

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		logs.Error("send email error: ", err)
	} else {
		logs.Info("success send email")
	}

}
