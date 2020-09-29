package service

import (
	"collect-homework-go/database"
	"collect-homework-go/e2e/request"
	"errors"
	"log"
)

// AdminLogin admin login
func AdminLogin(baseURL string,email string,password string) (ok bool,token string,err error){
	loginResponse,err := request.AdminLogin(baseURL+"/admin/login",&struct{
		Email string `json:"email"`
		Password string `json:"password"`
	}{
		Email: email,
		Password: password,
	});
	if err != nil {
		return false,"",err
	}
	if !loginResponse.Success || len(loginResponse.Data)==0 {
		log.Println(loginResponse)
		return false,"",errors.New("Login Response Fail ,With Token: "+loginResponse.Data)
	}
	return true,loginResponse.Data,nil
}

// AdminRegisterAndLogin admin register and login
func AdminRegisterAndLogin(baseURL string,email string,password string,name string) (ok bool,token string,err error){
	// request invitation code
	invitationResponse,err := request.AdminInvitationCode(baseURL+"/admin/invitationCode",&struct{
		Email string `json:"email"`
	}{
		Email: email,
	})
	if err != nil {
		return false,"",err
	}
	if !invitationResponse.Success {
		log.Println(invitationResponse)
		return false,"",errors.New("Invitation Response Fail")
	}

	invitationCode,err := database.Store.InvitationCode.SelectByEmail(email);
	if err != nil {
		return false,"",err
	}

	// request register
	registerResponse,err := request.AdminRegister(baseURL+"/admin/register",&struct{
		Email string `json:"email"`
		Password string `json:"password"`
		Name string `json:"name"`
		InvitationCode string `json:"invitationCode"`
	}{
		Email: email,
		Password: password,
		Name: name,
		InvitationCode: invitationCode.Code,
	})
	if err != nil {
		return false,"",err
	}
	if !registerResponse.Success {
		log.Println(registerResponse)
		return false,"",errors.New("Register Response Fail")
	}

	return AdminLogin(baseURL,email,password)
}