package mail

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"strconv"
	"strings"
)

type MailSender struct {
	server_address  string   // 服务端地址，支持带端口格式(smtp.qq.com:465)
	server_port     uint     // 服务端口，默认为465
	auth_user       string   // 用户名
	auth_password   string   // 密码
	sender          string   // 发件人，默认为 auth_user
	sender_username string   // 发件人名称，默认为 auth_user 按 @ 字符切片的前半部分
	content_type    string   // 内容类型格式，"text/plain; charset=UTF-8" | "text/html; charset=UTF-8"
	receiver        []string // 收件人
}

type OptionFunc func(*MailSender)

func initOptions(options ...OptionFunc) *MailSender {
	mail_sender := &MailSender{
		server_address:  "",
		server_port:     0,
		auth_user:       "",
		auth_password:   "",
		sender:          "",
		sender_username: "",
		content_type:    "text/plain; charset=UTF-8",
		receiver:        []string{},
	}
	for _, option_func := range options {
		option_func(mail_sender)
	}
	return mail_sender
}

func WithOptions(options MailSender) OptionFunc {
	return func(mail_sender *MailSender) {
		*mail_sender = options
	}
}

// 发件人为空时，使用 smtp_user 作为发件人
func WithServerPort(i uint) OptionFunc {
	return func(mail_sender *MailSender) {
		mail_sender.server_port = i
	}
}

// 发件人为空时，使用 smtp_user 作为发件人
func WithSender(s string) OptionFunc {
	return func(mail_sender *MailSender) {
		mail_sender.sender = s
	}
}

// 收件人为空时，使用 smtp_user 作为收件人
func WithReceiver(receiver []string) OptionFunc {
	return func(mail_sender *MailSender) {
		mail_sender.receiver = receiver
	}
}

// 发件人名称为空时，使用 smtp_user 中 @ 之前的部分作为发件人名称
func WithSenderUsername(s string) OptionFunc {
	return func(mail_sender *MailSender) {
		mail_sender.sender_username = s
	}
}

// 发件人名称为空时，使用 smtp_user 中 @ 之前的部分作为发件人名称
func WithContentType(s string) OptionFunc {
	return func(mail_sender *MailSender) {
		mail_sender.content_type = s
	}
}

func New(server_address string, auth_user string, auth_password string, options ...OptionFunc) *MailSender {
	mail_sender := initOptions(options...)
	mail_sender.server_address = server_address
	mail_sender.auth_user = auth_user
	mail_sender.auth_password = auth_password
	// 发件人为空时，使用 auth_user 作为发件人
	if mail_sender.sender == "" {
		mail_sender.sender = mail_sender.auth_user
	}
	// 发件人名称为空时，使用 auth_user 中 @ 之前的部分作为发件人名称
	if mail_sender.sender_username == "" {
		mail_sender.sender_username = strings.Split(mail_sender.sender, "@")[0]
	}
	// 收件人为空时，使用 auth_user 作为收件人
	if len(mail_sender.receiver) == 0 {
		mail_sender.receiver = []string{mail_sender.auth_user}
	}
	// 切片获取服务端口
	if mail_sender.server_port == 0 {
		temp_list := strings.Split(mail_sender.server_address, ":")
		if len(temp_list) == 2 {
			mail_sender.server_address = temp_list[0]
			server_port, _ := strconv.Atoi(temp_list[1])
			mail_sender.server_port = uint(server_port)
		} else {
			mail_sender.server_port = 465
		}
	}
	return mail_sender
}

func (mail_sender *MailSender) SendMail(mail_title string, mail_content string) (err error) {
	header := make(map[string]string)
	header["From"] = mail_sender.sender_username + "<" + mail_sender.sender + ">"
	header["To"] = strings.Join(mail_sender.receiver, ",")
	header["Subject"] = mail_title
	header["Content-Type"] = mail_sender.content_type
	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + mail_content
	auth := smtp.PlainAuth(
		"",
		mail_sender.auth_user,
		mail_sender.auth_password,
		mail_sender.server_address,
	)
	err = send_mail_using_tls(
		fmt.Sprintf("%s:%d", mail_sender.server_address, mail_sender.server_port),
		auth,
		mail_sender.sender,
		mail_sender.receiver,
		[]byte(message),
	)
	if err != nil {
		return err
	}
	return
}

// 参考 net/smtp 的func SendMail()
// 使用 net.Dial 连接 tls（SSL） 端口时，smtp.NewClient()会卡住且不提示err
// len(to)>1时，to[1]开始提示是密送
func send_mail_using_tls(addr string, auth smtp.Auth, from string, to []string, msg []byte) (err error) {
	c, err := smtp_dial(addr)
	if err != nil {
		return err
	}
	defer c.Close()
	if auth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(auth); err != nil {
				return err
			}
		}
	}
	if err = c.Mail(from); err != nil {
		return err
	}
	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}

func smtp_dial(addr string) (*smtp.Client, error) {
	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		return nil, err
	}
	host, _, _ := net.SplitHostPort(addr)
	return smtp.NewClient(conn, host)
}

// 单次调用，发送邮件
func DoSendMail(server_address string, auth_user string, auth_password string, mail_title string, mail_content string, options ...OptionFunc) (err error) {
	mail_sender := New(server_address, auth_user, auth_password, options...)
	err = mail_sender.SendMail(mail_title, mail_content)
	return
}
