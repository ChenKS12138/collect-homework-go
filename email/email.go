package email

import (
	"crypto/tls"

	"github.com/spf13/viper"
	"gopkg.in/mail.v2"
)

// SendMail send mail
func SendMail(to string,subject string,body string) error {
	if viper.GetBool("EMAIL_PREVENT") {
		return nil
	}
	user := viper.GetString("EMAIL_USER")
	password := viper.GetString("EMAIL_PASSWORD")
	dialer := mail.NewDialer("smtp.qq.com",587,user,password)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	m := mail.NewMessage()
	m.SetHeader("From",user)
	m.SetHeader("To",to)
	m.SetHeader("Subject",subject)
	m.SetAddressHeader("To",to,to)
	m.SetBody("text/html",body)
	if err := dialer.DialAndSend(m); err != nil {
		return err
	}
	return nil
}