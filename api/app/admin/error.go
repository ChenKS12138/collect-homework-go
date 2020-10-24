package admin

import (
	"github.com/ChenKS12138/collect-homework-go/util"
	"net/http"
)

var (
	// ErrAuthorization error authorization
	ErrAuthorization = &util.ErrResponse{
		HTTPStatusCode: http.StatusOK,
		StatusText: "Authorization Fall",
		ErrorText: "Invalid Email or Password",
	}

	// ErrInvitationCodeFrequently error invitation code too frequently
	ErrInvitationCodeFrequently = &util.ErrResponse{
		HTTPStatusCode: http.StatusOK,
		StatusText: "Invation Code Fail",
		ErrorText: "Request Invitation Code Too Frequently",
	}

	// ErrInvitationCodeWrong error invation code 
	ErrInvitationCodeWrong = &util.ErrResponse{
		HTTPStatusCode: http.StatusOK,
		StatusText: "Invation Code Fail",
		ErrorText: "Invitation Code Wrong",
	}
	// ErrEmailUsed error email used
	ErrEmailUsed = &util.ErrResponse{
		HTTPStatusCode: http.StatusOK,
		StatusText: "Register Fail",
		ErrorText: "Email Used",
	}
)