package domain

import (
	"fmt"
	"net/smtp"
	"redhat/notification-microservice/config"

	"github.com/jordan-wright/email"
)

// Email - holds email information
type Email struct {
	Recipients []string `json:"recipients"`
	Subject    string   `json:"subject"`
	Body       string   `json:"body"`
}

// EmailServer - interface for holding email server information and sending emails
type EmailServer interface {
	SendEmail(emailData *Email) error
}

// AuthEmailServer - holds Authenticated email server information
type AuthEmailServer struct {
	auth        smtp.Auth
	url         string
	fromAddress string
}

// NewAuthEmailServer - Holds all the email server configuration
func NewAuthEmailServer(config *config.Config) *AuthEmailServer {
	var emailServer AuthEmailServer

	// Set up authentication information.
	emailServer.auth = smtp.PlainAuth(
		"",
		config.SMTPUsername,
		config.SMTPPassword,
		config.SMTPServer,
	)
	emailServer.url = fmt.Sprintf("%s:%d", config.SMTPServer, config.SMTPPort)
	emailServer.fromAddress = config.SMTPUsername

	return &emailServer
}

// SendEmail - sends the actual email with the given details
func (aes *AuthEmailServer) SendEmail(emailData *Email) error {
	e := email.NewEmail()
	e.From = aes.fromAddress
	e.To = emailData.Recipients
	e.Subject = emailData.Subject
	e.Text = []byte(emailData.Body)
	return e.Send(aes.url, aes.auth)
}
