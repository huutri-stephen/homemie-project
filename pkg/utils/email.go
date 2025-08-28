package utils

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"net/smtp"
	"time"

	"homemie/config"
	"homemie/models"

	"gorm.io/gorm"
)

func SendVerificationEmail(cfg config.Config, db *gorm.DB, email, name, token string) error {
	var emailTemplate models.EmailTemplate
	if err := db.Where("name = ?", "VERIFY_EMAIL").First(&emailTemplate).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("verification email template not found")
		}
		return err
	}

	verificationURL := fmt.Sprintf("http://%s:%s%s/verify-email?email=%s&token=%s", cfg.Server.Host, cfg.Server.Port, cfg.Server.ApiVersion, email, token)

	data := struct {
		Name        string
		VerifyURL   string
		SupportEmail string
		Year        int
	}{
		Name:        name,
		VerifyURL:   verificationURL,
		SupportEmail: "support@homemie.com",
		Year:        time.Now().Year(),
	}

	t, err := template.New("verificationEmail").Parse(emailTemplate.Body)
	if err != nil {
		return fmt.Errorf("failed to parse email template: %w", err)
	}

	var body bytes.Buffer
	if err := t.Execute(&body, data); err != nil {
		return fmt.Errorf("failed to execute email template: %w", err)
	}

	auth := smtp.PlainAuth("", cfg.Email.SmtpUser, cfg.Email.SmtpPass, cfg.Email.SmtpHost)
	// You can use fmt.Println for debugging or use a proper logger if available
	fmt.Println(cfg.Email.SmtpUser, cfg.Email.SmtpPass, cfg.Email.SmtpHost)
	smtpAddr := fmt.Sprintf("%s:%s", cfg.Email.SmtpHost, cfg.Email.SmtpPort)
	to := []string{email}
	msg := []byte(
		"To: " + email + "\r\n" +
			"Subject: " + emailTemplate.Subject + "\r\n" +
			"MIME-version: 1.0;\r\n" +
			"Content-Type: text/html; charset=\"UTF-8\";\r\n\r\n" +
			body.String(),
	)

	if err := smtp.SendMail(smtpAddr, auth, cfg.Email.SenderEmail, to, msg); err != nil {
		return fmt.Errorf("failed to send verification email: %w", err)
	}

	return nil
}