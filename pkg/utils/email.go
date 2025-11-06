package utils

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"net/smtp"
	"time"

	"homemie/config"
	"homemie/db/models"
)

type EmailTemplates struct {
	cfg *config.Config
	db 	*sql.DB
}

func NewEmailTemplates(cfg *config.Config, db *sql.DB) *EmailTemplates {
	return &EmailTemplates{
		cfg: cfg,
		db:  db,
	}
}

func (e *EmailTemplates) SendVerificationEmail(ctx context.Context, email, name, token string) error {
	emailTemplate, err := models.EmailTemplates(models.EmailTemplateWhere.Name.EQ("VERIFY_EMAIL")).One(ctx, e.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("verification email template not found")
		}
		return err
	}

	verificationURL := fmt.Sprintf("http://%s:%s%s/verify-email?email=%s&token=%s", e.cfg.Server.Host, e.cfg.Server.Port, e.cfg.Server.ApiVersion, email, token)

	data := struct {
		Name         string
		VerifyURL    string
		SupportEmail string
		Year         int
	}{
		Name:         name,
		VerifyURL:    verificationURL,
		SupportEmail: "support@homemie.com",
		Year:         time.Now().Year(),
	}

	t, err := template.New("verificationEmail").Parse(emailTemplate.Body)
	if err != nil {
		return fmt.Errorf("failed to parse email template: %w", err)
	}

	var body bytes.Buffer
	if err := t.Execute(&body, data); err != nil {
		return fmt.Errorf("failed to execute email template: %w", err)
	}

	auth := smtp.PlainAuth("", e.cfg.Email.SmtpUser, e.cfg.Email.SmtpPass, e.cfg.Email.SmtpHost)
	smtpAddr := fmt.Sprintf("%s:%s", e.cfg.Email.SmtpHost, e.cfg.Email.SmtpPort)
	to := []string{email}
	msg := []byte(
		"To: " + email + "\r\n" +
			"Subject: " + emailTemplate.Subject + "\r\n" +
			"MIME-version: 1.0;\r\n" +
			"Content-Type: text/html; charset=\"UTF-8\";\r\n\r\n" +
			body.String(),
	)

	if err := smtp.SendMail(smtpAddr, auth, e.cfg.Email.SenderEmail, to, msg); err != nil {
		return fmt.Errorf("failed to send verification email: %w", err)
	}

	return nil
}

func (e *EmailTemplates) SendPasswordResetEmail(ctx context.Context, email, name, token string) error {
	emailTemplate, err := models.EmailTemplates(models.EmailTemplateWhere.Name.EQ("RESET_PASSWORD")).One(ctx, e.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("password reset email template not found")
		}
		return err
	}

	resetURL := fmt.Sprintf("http://%s:%s%s/reset-password?email=%s&token=%s", e.cfg.Server.Host, e.cfg.Server.Port, e.cfg.Server.ApiVersion, email, token)

	data := struct {
		Name         string
		ResetURL     string
		SupportEmail string
		Year         int
	}{
		Name:         name,
		ResetURL:     resetURL,
		SupportEmail: "support@homemie.com",
		Year:         time.Now().Year(),
	}

	t, err := template.New("passwordResetEmail").Parse(emailTemplate.Body)
	if err != nil {
		return fmt.Errorf("failed to parse email template: %w", err)
	}

	var body bytes.Buffer
	if err := t.Execute(&body, data); err != nil {
		return fmt.Errorf("failed to execute email template: %w", err)
	}

	auth := smtp.PlainAuth("", e.cfg.Email.SmtpUser, e.cfg.Email.SmtpPass, e.cfg.Email.SmtpHost)
	smtpAddr := fmt.Sprintf("%s:%s", e.cfg.Email.SmtpHost, e.cfg.Email.SmtpPort)
	to := []string{email}
	msg := []byte(
		"To: " + email + "\r\n" +
			"Subject: " + emailTemplate.Subject + "\r\n" +
			"MIME-version: 1.0;\r\n" +
			"Content-Type: text/html; charset=\"UTF-8\";\r\n\r\n" +
			body.String(),
	)

	if err := smtp.SendMail(smtpAddr, auth, e.cfg.Email.SenderEmail, to, msg); err != nil {
		return fmt.Errorf("failed to send password reset email: %w", err)
	}

	return nil
}
