package email

import (
	"bytes"
	"fmt"
	"html/template"
	"github.com/gohex/gohex/internal/application/port"
	"github.com/gohex/gohex/pkg/errors"
	"github.com/gohex/gohex/internal/domain/vo"
	"github.com/gohex/gohex/pkg/tracer"
	"gopkg.in/gomail.v2"
)

type smtpEmailService struct {
	config  SMTPConfig
	logger  Logger
	metrics MetricsReporter
}

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	WebsiteURL string
}

func NewSMTPEmailService(config SMTPConfig, logger Logger, metrics MetricsReporter) port.EmailService {
	return &smtpEmailService{
		config:  config,
		logger:  logger,
		metrics: metrics,
	}
}

func (s *smtpEmailService) loadTemplate(name string) (*template.Template, error) {
	templatePath := fmt.Sprintf("templates/emails/%s", name)
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		s.logger.Error("failed to load email template", 
			"template", name,
			"error", err,
		)
		return nil, fmt.Errorf("failed to load template %s: %w", name, err)
	}
	return tmpl, nil
}

func (s *smtpEmailService) SendWelcomeEmail(email string, name string) error {
	timer := s.metrics.StartTimer("email_send_duration", "type", "welcome")
	defer timer.Stop()

	data := map[string]interface{}{
		"Name": name,
	}

	if err := s.sendEmail(email, "Welcome!", "welcome.html", data); err != nil {
		s.metrics.IncrementCounter("email_send_failure", "type", "welcome")
		return err
	}

	s.metrics.IncrementCounter("email_send_success", "type", "welcome")
	return nil
}

func (s *smtpEmailService) SendPasswordResetEmail(email string, resetToken string) error {
	timer := s.metrics.StartTimer("email_send_duration", "type", "password_reset")
	defer timer.Stop()

	data := map[string]interface{}{
		"ResetToken": resetToken,
		"ResetURL":   fmt.Sprintf("%s/reset-password?token=%s", s.config.WebsiteURL, resetToken),
	}

	if err := s.sendEmail(email, "Reset Your Password", "password_reset.html", data); err != nil {
		s.metrics.IncrementCounter("email_send_failure", "type", "password_reset")
		return err
	}

	s.metrics.IncrementCounter("email_send_success", "type", "password_reset")
	return nil
}

func (s *smtpEmailService) SendPasswordChangedNotification(email string) error {
	timer := s.metrics.StartTimer("email_send_duration", "type", "password_changed")
	defer timer.Stop()

	if err := s.sendEmail(email, "Password Changed", "password_changed.html", nil); err != nil {
		s.metrics.IncrementCounter("email_send_failure", "type", "password_changed")
		return err
	}

	s.metrics.IncrementCounter("email_send_success", "type", "password_changed")
	return nil
}

func (s *smtpEmailService) sendEmail(to, subject, templateName string, data interface{}) error {
	tmpl, err := template.ParseFiles(fmt.Sprintf("templates/emails/%s", templateName))
	if err != nil {
		return err
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", s.config.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body.String())

	d := gomail.NewDialer(s.config.Host, s.config.Port, s.config.Username, s.config.Password)
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

// 其他方法实现... 