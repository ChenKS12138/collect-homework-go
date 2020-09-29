package request

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