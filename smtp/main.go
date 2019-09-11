package main

import (
	"fmt"
	"projects/testing-go/smtp/smtp"
)

func main() {
	fmt.Println("== start ==")

	cfg := smtp.Configuration{
		SMTP: smtp.Service{
			Username: "azure_011388ace0c90824a9d0ea7d695816ed@azure.com",
			Password: "Silkrode50832747",
			Host:     "smtp.sendgrid.net",
			Port:     587,
		},
	}
	smtpSrv, err := smtp.NewSMTPService(&cfg)
	if err != nil {
		fmt.Printf("failed to create smtp srv: %+v\n", err)
		return
	}

	smtpSrv.Ping()
	sendFrom := "noreply@ServiceAnonymousVPN.com"
	sendTo := []string{"xiao.xiao@tenoz.tw", "renhao.xiao@silkrode.com.tw"}
	err = smtpSrv.SendMail(sendFrom, sendTo, "Subject", "body", "html")
	if err != nil {
		fmt.Printf("failed to send mail: %+v\n", err)
	}

	fmt.Println("==done==")
}
