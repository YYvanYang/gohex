package output

type EmailService interface {
    SendWelcomeEmail(email string, name string) error
    SendPasswordResetEmail(email string, resetToken string) error
    SendPasswordChangedNotification(email string) error
    SendVerificationEmail(email string, verificationCode string) error
    SendLoginNotification(email string, ip string, userAgent string) error
    SendAccountLockedNotification(email string, reason string) error
}

type EmailTemplateConfig struct {
    TemplatePath string
    Subject      string
    From         string
    ReplyTo      string
    Attachments  []EmailAttachment
}

type EmailAttachment struct {
    Filename string
    Content  []byte
    MimeType string
} 