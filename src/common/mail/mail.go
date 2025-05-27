package mail

import (
	"TTCS/src/common/log"
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
	"html/template"
	"mime"
	"os"
	"path/filepath"
	"time"
)

type GmailService struct {
	*gmail.Service
}

func NewGmailService() *GmailService {
	ctx := context.Background()
	log.Info(ctx, "📧Connecting to email service... ")
	config := oauth2.Config{
		ClientID:     os.Getenv("MAIL_CLIENT_ID"),
		ClientSecret: os.Getenv("MAIL_CLIENT_SECRET"),
		Endpoint:     google.Endpoint,
		RedirectURL:  "http://localhost",
		Scopes:       []string{"https://www.googleapis.com/auth/gmail.send"},
	}
	token := oauth2.Token{
		AccessToken:  os.Getenv("MAIL_ACCESS_TOKEN"),
		TokenType:    "Bearer",
		RefreshToken: os.Getenv("MAIL_REFRESH_TOKEN"),
		Expiry:       time.Now(),
	}
	var tokenSource = config.TokenSource(ctx, &token)

	srv, err := gmail.NewService(context.Background(), option.WithTokenSource(tokenSource))
	if err != nil {
		log.Error(ctx, "Unable to retrieve Gmail client: %v", err)
		panic("Failed to login to gmail")
	}
	if srv != nil {
		log.Info(ctx, "📧Email service is initialized !")
	}
	return &GmailService{srv}
}

func (g *GmailService) SendEmailOAuth2(title, to string, data interface{}, template string) error {
	emailBody, err := g.parseTemplate(template, data)
	if err != nil {
		return errors.New("unable to parse email template")
	}

	var message gmail.Message

	emailTo := "To: " + to + "\r\n"
	subject := "Subject: " + mime.QEncoding.Encode("UTF-8", title) + "\n"
	mime2 := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"
	msg := []byte(emailTo + subject + mime2 + "\n" + emailBody)

	message.Raw = base64.URLEncoding.EncodeToString(msg)

	// Send the message
	_, err = g.Users.Messages.Send("me", &message).Do()
	if err != nil {
		return err
	}
	return nil
}

func (g *GmailService) parseTemplate(templateFileName string, data interface{}) (string, error) {
	templatePath, err := filepath.Abs(fmt.Sprintf("src/common/mail/template/%s", templateFileName))
	if err != nil {
		return "", err
	}
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return "", err
	}
	body := buf.String()
	return body, nil
}
