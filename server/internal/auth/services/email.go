package service

import (
	"fmt"
	"net/url"

	"github.com/resendlabs/resend-go"
	"github.com/tdmdh/fit-up-server/shared/config"
)

func SendVerificationEmail(toEmail string, token string) error {
	client := resend.NewClient(config.NewConfig().ResendAPIKey)

	params := &resend.SendEmailRequest{
		From:    "noreply@lornian.com",
		To:      []string{toEmail},
		Subject: "Verify your email address",
		Html:    generateVerificationEmailHTML(token),
	}

	_, err := client.Emails.Send(params)
	return err
}

func generateVerificationEmailHTML(token string) string {
	cfg := config.NewConfig()
	frontendURL := cfg.FrontendURL
	if frontendURL == "" {
		// Default to localhost for development
		frontendURL = "http://localhost:3000"
	}
	verifyURL := fmt.Sprintf("%s/verify?token=%s", frontendURL, url.QueryEscape(token))
	return `
	<!DOCTYPE html>
		<html>
		<head>
			<style>
				body {
					font-family: Arial, sans-serif;
					background-color: #F8FAFC;
					color: #334155;
					margin: 0;
					padding: 0;
				}
				.container {
					max-width: 600px;
					margin: 40px auto;
					background: #FFFFFF;
					padding: 32px;
					border-radius: 8px;
					box-shadow: 0 4px 12px rgba(0,0,0,0.05);
				}
				h1 {
					color: #60A5FA;
					margin-bottom: 16px;
				}
				p {
					line-height: 1.6;
				}
				a.button {
					display: inline-block;
					margin-top: 24px;
					padding: 12px 24px;
					background-color: #60A5FA;
					color: #FFFFFF;
					text-decoration: none;
					border-radius: 4px;
					font-weight: bold;
				}
				a.button:hover {
					background-color: #3B82F6;
				}
				.footer {
					margin-top: 32px;
					font-size: 14px;
					color: #64748B;
				}
			</style>
		</head>
		<body>
			<div class="container">
				<h1>Verify Your Email Address</h1>
				<p>Thank you for signing up! Please click the button below to verify your email address and complete your registration.</p>
				<a href="` + verifyURL + `" class="button">Verify Email Address</a>
				<p>If the button doesn't work, copy and paste this link into your browser:</p>
				<p style="word-break: break-all; color: #60A5FA;">` + verifyURL + `</p>
				<div class="footer">
					<p>If you didn't create an account, you can safely ignore this email.</p>
					<p>This verification link will expire in 24 hours.</p>
				</div>
			</div>
		</body>
		</html>            
	`
}

func SendPasswordResetEmail(toEmail string, token string) error {
	client := resend.NewClient(config.NewConfig().ResendAPIKey)

	resetURL := fmt.Sprintf("https://app.lornian.com/reset-password?token=%s", token)

	params := &resend.SendEmailRequest{
		From:    "noreply@lornian.com",
		To:      []string{toEmail},
		Subject: "Reset your password",
		Html:    generatePasswordResetEmailHTML(resetURL, token),
	}

	_, err := client.Emails.Send(params)
	return err
}

func generatePasswordResetEmailHTML(resetURL, token string) string {
	return fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<style>
				body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
				.container { max-width: 600px; margin: 0 auto; padding: 20px; }
				.header { background-color: #f8f9fa; padding: 20px; text-align: center; }
				.content { padding: 20px; }
				.button { 
					display: inline-block; 
					padding: 12px 24px; 
					background-color: #007bff; 
					color: white; 
					text-decoration: none; 
					border-radius: 4px; 
					margin: 20px 0;
				}
				.footer { font-size: 12px; color: #666; margin-top: 20px; }
			</style>
		</head>
		<body>
			<div class="container">
				<div class="header">
					<h1>Password Reset Request</h1>
				</div>
				<div class="content">
					<p>Hello,</p>
					<p>We received a request to reset your password. If you made this request, click the button below to reset your password:</p>
					<a href="%s" class="button">Reset Password</a>
					<p>If the button doesn't work, you can copy and paste this link into your browser:</p>
					<p><a href="%s">%s</a></p>
					<p>This link will expire in 1 hour for security reasons.</p>
					<p>If you didn't request a password reset, please ignore this email. Your password will remain unchanged.</p>
				</div>
				<div class="footer">
					<p>This is an automated message, please do not reply to this email.</p>
				</div>
			</div>
		</body>
		</html>
	`, resetURL, resetURL, resetURL)
}
