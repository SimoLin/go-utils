package test

import (
	"fmt"
	"testing"

	"github.com/SimoLin/go-utils/mail"
)

func TestMainServer(t *testing.T) {
	server_address := "smtp.qq.com:465"
	auth_user := "777777777@qq.com"
	auth_password := "your_smtp_auth_code"
	sender := "777777777@qq.com"
	sender_username := "your_name"
	reveiver := []string{"777777777@qq.com", "777777777@qq.com", "777777777@qq.com"}
	mail_title := "Test Mail Title"
	mail_content := "Test Mail Content"

	mail_sender := mail.New(
		server_address, auth_user, auth_password,
		mail.WithSender(sender),
		mail.WithSenderUsername(sender_username),
		mail.WithReceiver(reveiver),
	)
	err := mail_sender.SendMail(mail_title, mail_content)
	if err != nil {
		fmt.Println(err)
		t.Failed()
	}
}

func TestDoSendMail(t *testing.T) {
	server_address := "smtp.qq.com:465"
	auth_user := "777777777@qq.com"
	auth_password := "your_smtp_auth_code"
	sender := "777777777@qq.com"
	sender_username := "your_name"
	reveiver := []string{"777777777@qq.com", "777777777@qq.com", "777777777@qq.com"}
	mail_title := "Test Mail Title"
	mail_content := "Test Mail Content"

	err := mail.DoSendMail(
		server_address, auth_user, auth_password, mail_title, mail_content,
		mail.WithSender(sender),
		mail.WithSenderUsername(sender_username),
		mail.WithReceiver(reveiver),
	)
	if err != nil {
		fmt.Println(err)
		t.Failed()
	}
}
