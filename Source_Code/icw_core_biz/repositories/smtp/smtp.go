package smtp

import (
	"crypto/tls"
	"errors"
	"fmt"
	"mime"
	"net/smtp"
	"strings"

	"icw_core_biz/configs"
)

// Repository SMTP 简单邮件传输协议服务
type Repository struct {
	SMTPHost            string
	SMTPPort            int
	SMTPPassword        string
	SMTPFromName        string
	SMTPFromEmail       string
	EmailCodeTTLMinutes int
}

func NewRepository(cfg configs.Config) *Repository {
	return &Repository{
		SMTPHost:            cfg.SMTPHost,
		SMTPPort:            cfg.SMTPPort,
		SMTPPassword:        cfg.SMTPPassword,
		SMTPFromName:        cfg.SMTPFromName,
		SMTPFromEmail:       cfg.SMTPFromEmail,
		EmailCodeTTLMinutes: int(cfg.EmailCodeTTL.Minutes()),
	}
}

// Configured 校验 SMTP 配置
func (r *Repository) Configured() bool {
	return r.SMTPHost != "" && r.SMTPPort != 0 && r.SMTPPassword != "" && r.SMTPFromName != "" && r.SMTPFromEmail != "" && r.EmailCodeTTLMinutes != 0
}

// SendEmailCode 发送验证码邮件
func (r *Repository) SendEmailCode(to, scene, code string) error {
	if !r.Configured() {
		return errors.New("SMTP service not configured")
	}
	subject, sceneName, err := emailSceneText(scene)
	if err != nil {
		return err
	}
	body := buildEmailCodeHTML(sceneName, code, r.EmailCodeTTLMinutes)
	return r.send(to, subject, body)
}

// send 发送 HTML 邮件
func (r *Repository) send(to, subject, htmlBody string) error {
	fromName := mime.QEncoding.Encode("UTF-8", r.SMTPFromName)
	encodedSubject := mime.QEncoding.Encode("UTF-8", subject)
	message := strings.Join([]string{
		fmt.Sprintf("From: %s <%s>", fromName, r.SMTPFromEmail),
		fmt.Sprintf("To: %s", to),
		"Subject: " + encodedSubject,
		"MIME-Version: 1.0",
		"Content-Type: text/html; charset=UTF-8",
		"",
		htmlBody,
	}, "\r\n")

	addr := fmt.Sprintf("%s:%d", r.SMTPHost, r.SMTPPort)
	auth := smtp.PlainAuth("", r.SMTPFromEmail, r.SMTPPassword, r.SMTPHost)

	conn, err := tls.Dial("tcp", addr, &tls.Config{ServerName: r.SMTPHost})
	if err != nil {
		return err
	}
	defer func() {
		_ = conn.Close()
	}()

	client, err := smtp.NewClient(conn, r.SMTPHost)
	if err != nil {
		return err
	}
	defer func() {
		_ = client.Quit()
	}()

	if err := client.Auth(auth); err != nil {
		return err
	}
	if err := client.Mail(r.SMTPFromEmail); err != nil {
		return err
	}
	if err := client.Rcpt(to); err != nil {
		return err
	}

	writer, err := client.Data()

	if err != nil {
		return err
	}
	if _, err := writer.Write([]byte(message)); err != nil {
		return err
	}

	return writer.Close()
}
