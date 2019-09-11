package smtp

import (
	"encoding/base64"
	"fmt"
	"net/smtp"
	"strings"
)

// Service 是一個用 SMTP 寄郵件的物件
type Service struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
}

// Configuration is configuration
type Configuration struct {
	SMTP Service
}

// NewSMTPService 回傳一個 SMTP 服務物件
func NewSMTPService(cfg *Configuration) (Service, error) {
	smtpService := Service{
		Username: cfg.SMTP.Username,
		Password: cfg.SMTP.Password,
		Host:     cfg.SMTP.Host,
		Port:     cfg.SMTP.Port,
	}

	return smtpService, nil
}

// SendMail 使用 SMTP 寄發電子郵件
func (ss *Service) SendMail(sendFrom string, sendTo []string, subject, body, mailtype string) error {
	auth := smtp.PlainAuth("", ss.Username, ss.Password, ss.Host)
	host := fmt.Sprintf("%s:%d", ss.Host, ss.Port)

	subject = fmt.Sprintf("=?UTF-8?B?%s?=", base64.StdEncoding.EncodeToString([]byte(subject)))
	contentType := "Content-Type: text/plain; charset=UTF-8"
	if len(mailtype) > 0 {
		contentType = fmt.Sprintf("Content-Type: text/%s; charset=UTF-8", mailtype)
	}
	emailFormat := fmt.Sprintf("To: %s\r\nFrom: %s >\r\nSubject: %s\r\n%s\r\n\r\n%s", strings.Join(sendTo, ","), sendFrom, subject, contentType, body)
	msg := []byte(emailFormat)

	fmt.Printf("\n\n%+v \n%+v \n%+v \n%+v \n%+v \n", host, auth, ss.Username, sendTo, msg)
	err := smtp.SendMail(host, auth, sendFrom, sendTo, msg)

	return err
}

// Ping 用來檢查 smtp 連線
func (ss *Service) Ping() error {
	host := fmt.Sprintf("%s:%d", ss.Host, ss.Port)
	_, err := smtp.Dial(host)

	return err
}
