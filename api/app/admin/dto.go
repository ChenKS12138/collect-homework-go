package admin

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

// LoginDto login dto
type LoginDto struct {
	Email string `json:"email"`
	Password string `json:"password"`
}


func (l *LoginDto)validate() error {
	err := &validation.Errors {
		"email":validation.Validate(l.Email,validation.Required, is.Email ),
		"password":validation.Validate(l.Password,validation.Required),
	}
	return err.Filter()
}

// InvitationCodeDto invation code dto
type InvitationCodeDto struct {
	Email string `json:"email"`
}

func (i *InvitationCodeDto)validate() error {
	err := &validation.Errors{
		"email":validation.Validate(i.Email,validation.Required,is.Email),
	}
	return err.Filter()
}

// RegisterDto register dto
type RegisterDto struct {
	LoginDto
	Name string `json:"name"`
	InvitationCode string `json:"invitationCode"`
}

func (r *RegisterDto)validate() error {
	err := &validation.Errors{
		"email":validation.Validate(r.Email,validation.Required, is.Email ),
		"password":validation.Validate(r.Password,validation.Required),
		"name":validation.Validate(r.Name,validation.Required),
		"invitationCode": validation.Validate(r.InvitationCode,validation.Required),
	}
	return err.Filter()
}