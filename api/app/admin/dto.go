package admin

import (
	"github.com/ChenKS12138/collect-homework-go/auth"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

// LoginDto login dto
type LoginDto struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	Captcha      string `json:"captcha"`
	CaptchaToken string `json:"captchaToken"`
}

func (l *LoginDto) validate() error {
	err := &validation.Errors{
		"email":        validation.Validate(l.Email, validation.Required, is.Email),
		"password":     validation.Validate(l.Password, validation.Required),
		"captcha":      validation.Validate(l.Captcha, validation.Required),
		"captchaToken": validation.Validate(l.CaptchaToken, validation.Required),
	}
	return err.Filter()
}

// InvitationCodeDto invation code dto
type InvitationCodeDto struct {
	Email        string `json:"email"`
	Captcha      string `json:"captcha"`
	CaptchaToken string `json:"captchaToken"`
}

func (i *InvitationCodeDto) validate() error {
	err := &validation.Errors{
		"email":        validation.Validate(i.Email, validation.Required, is.Email),
		"captcha":      validation.Validate(i.Captcha, validation.Required),
		"captchaToken": validation.Validate(i.CaptchaToken, validation.Required),
	}
	return err.Filter()
}

// RegisterDto register dto
type RegisterDto struct {
	LoginDto
	Name           string `json:"name"`
	InvitationCode string `json:"invitationCode"`
	Captcha        string `json:"captcha"`
	CaptchaToken   string `json:"captchaToken"`
}

func (r *RegisterDto) validate() error {
	err := &validation.Errors{
		"email":          validation.Validate(r.Email, validation.Required, is.Email),
		"password":       validation.Validate(r.Password, validation.Required),
		"name":           validation.Validate(r.Name, validation.Required),
		"invitationCode": validation.Validate(r.InvitationCode, validation.Required),
		"captcha":        validation.Validate(r.Captcha, validation.Required),
		"captchaToken":   validation.Validate(r.CaptchaToken, validation.Required),
	}
	return err.Filter()
}

// SubTokenDto subToken dto
type SubTokenDto struct {
	Expire   int64     `json:"expire"` // unit time.Minute
	AuthCode auth.Code `json:"authCode"`
}

func (s *SubTokenDto) validate() error {
	err := &validation.Errors{
		"expire": validation.Validate(s.Expire, validation.Min(1), validation.Required),
		"authCode": validation.Validate(s.AuthCode, validation.In(
			auth.CodeFileR+auth.CodeFileX, // download by subToken
		)),
	}
	return err.Filter()
}
