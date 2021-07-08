package outputs

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/smtp"
	"strconv"
	"strings"

	"github.com/aquasecurity/postee/formatting"
	"github.com/aquasecurity/postee/layout"
)

var (
	errThereIsNoRecipient = errors.New("there is no recipient")
)

type EmailOutput struct {
	Name       string
	User       string
	Password   string
	Host       string
	Port       int
	Sender     string
	Recipients []string
	UseMX      bool
}

func (email *EmailOutput) GetName() string {
	return email.Name
}

func (email *EmailOutput) Init() error {
	log.Printf("Starting Email output %q...", email.Name)
	if email.Sender == "" {
		email.Sender = email.User
	}
	return nil
}

func (email *EmailOutput) Terminate() error {
	log.Printf("Email output terminated\n")
	return nil
}

func (email *EmailOutput) GetLayoutProvider() layout.LayoutProvider {
	return new(formatting.HtmlProvider)
}

func (email *EmailOutput) Send(content map[string]string) error {
	subject := content["title"]
	body := content["description"]
	recipients := getHandledRecipients(email.Recipients, &content, email.Name)
	if len(recipients) == 0 {
		return errThereIsNoRecipient
	}

	if email.UseMX {
		sendViaMxServers(email.Sender, subject, body, recipients)
		return nil
	}

	msg := fmt.Sprintf(
		"To: %s\r\n"+
			"From: %s\r\n"+
			"Subject: %s\r\n"+
			"Content-Type: text/html; charset=UTF-8\r\n\r\n%s\r\n",
		strings.Join(recipients, ","), email.Sender, subject, body)

	auth := smtp.PlainAuth("", email.User, email.Password, email.Host)
	err := smtp.SendMail(email.Host+":"+strconv.Itoa(email.Port), auth, email.Sender, recipients, []byte(msg))
	if err != nil {
		log.Println("SendMail Error:", err)
		log.Printf("From: %q, to %v via %q", email.Sender, email.Recipients, email.Host)
		return err
	}
	log.Println("Email was sent successfully!")
	return nil
}

func sendViaMxServers(from, subj, msg string, recipients []string) {
	for _, rcpt := range recipients {
		at := strings.LastIndex(rcpt, "@")
		if at < 0 {
			log.Printf("%q isn't email", rcpt)
			continue
		}
		host := rcpt[at+1:]
		mxs, err := net.LookupMX(host)
		if err != nil {
			log.Print(err)
			continue
		}
		for _, mx := range mxs {
			message := fmt.Sprintf(
				"To: %s\r\n"+
					"From: %s\r\n"+
					"Subject: %s\r\n"+
					"Content-Type: text/html; charset=UTF-8\r\n\r\n%s\r\n", rcpt, from, subj, msg)

			if err := smtp.SendMail(mx.Host+":25", nil, from, []string{rcpt}, []byte(message)); err != nil {
				log.Printf("SendMail error to %q via %q", rcpt, mx.Host)
				log.Print(err)
				continue
			}
			log.Printf("The message to %q was sent successful via %q!", rcpt, mx.Host)
			break
		}
	}
}
