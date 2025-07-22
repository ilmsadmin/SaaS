package services

import (
	"fmt"
	"log"
	"net/smtp"
	"strings"
)

type SMTPEmailService struct {
	host     string
	port     string
	username string
	password string
	from     string
}

func NewSMTPEmailService(host, port, username, password, from string) *SMTPEmailService {
	return &SMTPEmailService{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
	}
}

// SendInvitationEmail sends user invitation email
func (s *SMTPEmailService) SendInvitationEmail(email, token, inviterName, tenantName string) error {
	subject := fmt.Sprintf("You're invited to join %s", tenantName)

	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Team Invitation</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
        <h1 style="color: #2563eb;">You're Invited!</h1>
        
        <p>Hi there,</p>
        
        <p><strong>%s</strong> has invited you to join <strong>%s</strong>.</p>
        
        <p>Click the button below to accept your invitation and set up your account:</p>
        
        <div style="text-align: center; margin: 30px 0;">
            <a href="http://localhost:3000/auth/accept-invitation?token=%s" 
               style="background-color: #2563eb; color: white; padding: 12px 24px; text-decoration: none; border-radius: 6px; display: inline-block;">
                Accept Invitation
            </a>
        </div>
        
        <p>This invitation will expire in 7 days.</p>
        
        <p>If you didn't expect this invitation, you can safely ignore this email.</p>
        
        <hr style="margin: 30px 0;">
        <p style="font-size: 12px; color: #666;">
            If the button doesn't work, you can copy and paste this link into your browser:<br>
            http://localhost:3000/auth/accept-invitation?token=%s
        </p>
    </div>
</body>
</html>`, inviterName, tenantName, token, token)

	return s.sendEmail(email, subject, body)
}

// SendPasswordResetEmail sends password reset email
func (s *SMTPEmailService) SendPasswordResetEmail(email, token, userName string) error {
	subject := "Reset Your Password"

	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Password Reset</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
        <h1 style="color: #dc2626;">Password Reset Request</h1>
        
        <p>Hi %s,</p>
        
        <p>We received a request to reset your password. If you made this request, click the button below to reset your password:</p>
        
        <div style="text-align: center; margin: 30px 0;">
            <a href="http://localhost:3000/auth/reset-password?token=%s" 
               style="background-color: #dc2626; color: white; padding: 12px 24px; text-decoration: none; border-radius: 6px; display: inline-block;">
                Reset Password
            </a>
        </div>
        
        <p>This link will expire in 1 hour for security reasons.</p>
        
        <p>If you didn't request a password reset, you can safely ignore this email. Your password will not be changed.</p>
        
        <hr style="margin: 30px 0;">
        <p style="font-size: 12px; color: #666;">
            If the button doesn't work, you can copy and paste this link into your browser:<br>
            http://localhost:3000/auth/reset-password?token=%s
        </p>
    </div>
</body>
</html>`, userName, token, token)

	return s.sendEmail(email, subject, body)
}

// SendVerificationEmail sends email verification email
func (s *SMTPEmailService) SendVerificationEmail(email, token, userName string) error {
	subject := "Verify Your Email Address"

	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Email Verification</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
        <h1 style="color: #059669;">Verify Your Email</h1>
        
        <p>Hi %s,</p>
        
        <p>Welcome! Please verify your email address to complete your account setup.</p>
        
        <div style="text-align: center; margin: 30px 0;">
            <a href="http://localhost:3000/auth/verify-email?token=%s" 
               style="background-color: #059669; color: white; padding: 12px 24px; text-decoration: none; border-radius: 6px; display: inline-block;">
                Verify Email
            </a>
        </div>
        
        <p>This verification link will expire in 24 hours.</p>
        
        <p>If you didn't create an account, you can safely ignore this email.</p>
        
        <hr style="margin: 30px 0;">
        <p style="font-size: 12px; color: #666;">
            If the button doesn't work, you can copy and paste this link into your browser:<br>
            http://localhost:3000/auth/verify-email?token=%s
        </p>
    </div>
</body>
</html>`, userName, token, token)

	return s.sendEmail(email, subject, body)
}

// sendEmail sends an email using SMTP
func (s *SMTPEmailService) sendEmail(to, subject, body string) error {
	// For development, just log the email
	log.Printf("SEND EMAIL TO: %s", to)
	log.Printf("SUBJECT: %s", subject)
	log.Printf("BODY: %s", body)

	// If SMTP is not configured, just return without error
	if s.host == "" || s.username == "" {
		log.Printf("SMTP not configured, email not sent")
		return nil
	}

	// Build email message
	message := s.buildMessage(s.from, to, subject, body)

	// SMTP auth
	auth := smtp.PlainAuth("", s.username, s.password, s.host)

	// Send email
	addr := fmt.Sprintf("%s:%s", s.host, s.port)
	err := smtp.SendMail(addr, auth, s.from, []string{to}, []byte(message))
	if err != nil {
		log.Printf("Failed to send email: %v", err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Printf("Email sent successfully to %s", to)
	return nil
}

// buildMessage builds the email message
func (s *SMTPEmailService) buildMessage(from, to, subject, body string) string {
	var msg strings.Builder

	msg.WriteString(fmt.Sprintf("From: %s\r\n", from))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", to))
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
	msg.WriteString("\r\n")
	msg.WriteString(body)

	return msg.String()
}

// MockEmailService for testing and development
type MockEmailService struct{}

func NewMockEmailService() *MockEmailService {
	return &MockEmailService{}
}

func (s *MockEmailService) SendInvitationEmail(email, token, inviterName, tenantName string) error {
	log.Printf("MOCK EMAIL - Invitation sent to %s from %s for tenant %s (token: %s)",
		email, inviterName, tenantName, token)
	return nil
}

func (s *MockEmailService) SendPasswordResetEmail(email, token, userName string) error {
	log.Printf("MOCK EMAIL - Password reset sent to %s for user %s (token: %s)",
		email, userName, token)
	return nil
}

func (s *MockEmailService) SendVerificationEmail(email, token, userName string) error {
	log.Printf("MOCK EMAIL - Verification sent to %s for user %s (token: %s)",
		email, userName, token)
	return nil
}
