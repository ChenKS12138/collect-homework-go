package service

import (
	"errors"
	"log"

	"github.com/ChenKS12138/collect-homework-go/database"
	"github.com/ChenKS12138/collect-homework-go/testing/request"
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
		return false,"", errors.New(loginResponse.ErrorText)
	}
	return true,loginResponse.Data,nil
}

// AdminRegisterAndLogin admin register and login
func AdminRegisterAndLogin(baseURL string,email string,password string,name string) (ok bool,token string,err error){
	// request invitation code
	_,err = AdminInvitationCode(baseURL,email)

	if err != nil {
		return false,"",err
	}

	invitationCode,err := database.Store.InvitationCode.SelectByEmail(email);
	if err != nil {
		return false,"",err
	}

	// request register
	_,err = AdminRegister(baseURL,email,password,name,invitationCode.Code)
	if err != nil {
		return false,"",err
	}

	return AdminLogin(baseURL,email,password)
}

// AdminInvitationCode admin invitation code
func AdminInvitationCode(baseURL string,email string)(ok bool,err error){
	invitationResponse,err := request.AdminInvitationCode(baseURL+"/admin/invitationCode",&struct{
		Email string `json:"email"`
	}{
		Email: email,
	})
	if err != nil {
		return false,err
	}
	if !invitationResponse.Success {
		log.Println(invitationResponse)
		return false,errors.New(invitationResponse.ErrorText)
	}
	return true,nil
}

// AdminRegister admin register
func AdminRegister(baseURL string,email string,password string,name string,invitationCode string)(ok bool,err error){
	registerResponse,err := request.AdminRegister(baseURL+"/admin/register",&struct{
		Email string `json:"email"`
		Password string `json:"password"`
		Name string `json:"name"`
		InvitationCode string `json:"invitationCode"`
	}{
		Email: email,
		Password: password,
		Name: name,
		InvitationCode: invitationCode,
	})
	if err != nil {
		return false,err
	}
	if !registerResponse.Success {
		log.Println(registerResponse)
		return false,errors.New(registerResponse.ErrorText)
	}
	return true,nil
}

// AdminStatus admin status
func AdminStatus(baseURL string,token string) (bool,*struct{
	FileCount int `json:"fileCount"`
	ProjectCount int `json:"projectCount"`
	TotalSize int64 `json:"totalSize"`
},error) {
	result,err := request.AdminStatus(baseURL+"/admin/status",token)
	if err != nil {
		return false,nil,err
	}
	return true,result,nil
}

// AdminSubToken admin subToken
func AdminSubToken(baseURL string,token string,expire int64,authCode uint32) (bool,string) {
	result,err := request.AdminSubToken(baseURL+"/admin/subToken",token,&struct{
		Expire int64 `json:"expire"`
		AuthCode uint32 `json:"authCode"`
	}{
		Expire: expire,
		AuthCode: authCode,
	})
	if err != nil {
		log.Println(err)
		return false,""
	}
	return result.Success,result.Data
}