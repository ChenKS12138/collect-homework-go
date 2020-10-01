package request

import (
	"errors"
	"net/http"
)

// AdminLogin admin login
func AdminLogin(url string,loginDto *struct {
	Email string `json:"email"`
	Password string `json:"password"`
}) (*struct{
	BasicResponse
	Data string `json:"data"`
},error) {
	response := &struct{
		BasicResponse
		Data string `json:"data"`
	}{}
	err := PostRequest(url,nil,loginDto,&response)
	if err != nil {
		return nil,err
	}
	return response,nil
}

// AdminInvitationCode admin invitation code
func AdminInvitationCode(url string ,invitationDto *struct {
	Email string `json:"email"`
}) (*struct {
	BasicResponse
},error) {
	response := &struct {
		BasicResponse
	}{}
	err := PostRequest(url,nil,invitationDto,&response)
	if err != nil {
		return nil,err
	}
	return response,nil
}

// AdminRegister admin register
func AdminRegister(url string,registerDto *struct {
	Email string `json:"email"`
	Password string `json:"password"`
	Name string `json:"name"`
	InvitationCode string `json:"invitationCode"`
}) (*struct {
	BasicResponse
},error){
	response := &struct {
		BasicResponse
	}{}
	err := PostRequest(url,nil,registerDto,&response)
	if err != nil {
		return nil,err
	}
	return response,nil
}

// AdminStatus admin status
func AdminStatus(url string,token string) ( *struct{
	FileCount int `json:"fileCount"`
	ProjectCount int `json:"projectCount"`
	TotalSize int64 `json:"totalSize"`
},error) {
	response := & struct {
		BasicResponse
		Data struct{
			FileCount int `json:"fileCount"`
			ProjectCount int `json:"projectCount"`
			TotalSize int64 `json:"totalSize"`
		} `json:"data"`
	}{}
	header := &http.Header{}
	header.Set("Authorization","Bearer "+token)
	err := GetRequest(url,header,nil,response)
	if err != nil {
		return nil,err
	}
	if !response.Success {
		return nil,errors.New(response.ErrorText)
	}
	return &response.Data,nil
}