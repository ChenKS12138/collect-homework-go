package email

import (
	"github.com/spf13/viper"
	"gopkg.in/mail.v2"
)

// Dialer mail dialer
var dialer *mail.Dialer
var User string


func init() {
	viper.AutomaticEnv()
	User = viper.GetString("EMAIL_USER")
	password := viper.GetString("EMAIL_PASSWORD")
	dialer = mail.NewDialer("smtp.qq.com",587,User,password)
	dialer.StartTLSPolicy = mail.MandatoryStartTLS
}

// SendMail send mail
func SendMail(to string,subject string,body string,userName string) error {
	viper.AutomaticEnv()
	m := mail.NewMessage()
	m.SetHeader("From",User)
	m.SetHeader("To",to)
	m.SetHeader("Subject",subject)
	m.SetAddressHeader("To",to,userName)
	m.SetBody("text/html",body)
	if err := dialer.DialAndSend(m); err != nil {
		return err
	}
	return nil
}