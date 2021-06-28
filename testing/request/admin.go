package request

import (
	"errors"
	"net/http"
)

// AdminLogin admin login
func AdminLogin(url string, loginDto *struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}) (*struct {
	BasicResponse
	Data string `json:"data"`
}, error) {
	response := &struct {
		BasicResponse
		Data string `json:"data"`
	}{}
	nextLoginDto := &struct {
		Email        string `json:"email"`
		Password     string `json:"password"`
		Captcha      string `json:"captcha"`
		CaptchaToken string `json:"captchaToken"`
	}{
		Email:        loginDto.Email,
		Password:     loginDto.Password,
		Captcha:      "123",
		CaptchaToken: "123",
	}
	err := PostRequest(url, nil, nextLoginDto, &response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// AdminInvitationCode admin invitation code
func AdminInvitationCode(url string, invitationDto *struct {
	Email string `json:"email"`
}) (*struct {
	BasicResponse
}, error) {
	response := &struct {
		BasicResponse
	}{}
	nextInvitationDto := &struct {
		Email        string `json:"email"`
		Captcha      string `json:"captcha"`
		CaptchaToken string `json:"captchaToken"`
	}{
		Email:        invitationDto.Email,
		Captcha:      "123",
		CaptchaToken: "123",
	}
	err := PostRequest(url, nil, nextInvitationDto, &response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// AdminRegister admin register
func AdminRegister(url string, registerDto *struct {
	Email          string `json:"email"`
	Password       string `json:"password"`
	Name           string `json:"name"`
	InvitationCode string `json:"invitationCode"`
}) (*struct {
	BasicResponse
}, error) {
	response := &struct {
		BasicResponse
	}{}
	nextRegisterDto := &struct {
		Email          string `json:"email"`
		Password       string `json:"password"`
		Name           string `json:"name"`
		InvitationCode string `json:"invitationCode"`
		Captcha        string `json:"captcha"`
		CaptchaToken   string `json:"captchaToken"`
	}{
		Email:          registerDto.Email,
		Password:       registerDto.Password,
		Name:           registerDto.Name,
		InvitationCode: registerDto.InvitationCode,
		Captcha:        "123",
		CaptchaToken:   "123",
	}
	err := PostRequest(url, nil, nextRegisterDto, &response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// AdminStatus admin status
func AdminStatus(url string, token string) (*struct {
	FileCount    int   `json:"fileCount"`
	ProjectCount int   `json:"projectCount"`
	TotalSize    int64 `json:"totalSize"`
}, error) {
	response := &struct {
		BasicResponse
		Data struct {
			FileCount    int   `json:"fileCount"`
			ProjectCount int   `json:"projectCount"`
			TotalSize    int64 `json:"totalSize"`
		} `json:"data"`
	}{}
	header := &http.Header{}
	header.Set("Authorization", "Bearer "+token)
	err := GetRequest(url, header, nil, response)
	if err != nil {
		return nil, err
	}
	if !response.Success {
		return nil, errors.New(response.ErrorText)
	}
	return &response.Data, nil
}

// AdminSubToken admin subToken
func AdminSubToken(url string, token string, subTokenDto *struct {
	Expire   int64  `json:"expire"`
	AuthCode uint32 `json:"authCode"`
}) (
	*struct {
		BasicResponse
		Data string `json:"data"`
	}, error) {
	response := &struct {
		BasicResponse
		Data string `json:"data"`
	}{}
	header := &http.Header{}
	header.Set("Authorization", "Bearer "+token)

	err := PostRequest(url, header, subTokenDto, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}
