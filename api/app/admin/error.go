package admin

import (
	"net/http"

	"github.com/ChenKS12138/collect-homework-go/util"
)

var (
	// ErrAuthorization error authorization
	ErrAuthorization = &util.ErrResponse{
		HTTPStatusCode: http.StatusOK,
		StatusText: "Authorization Fall",
		ErrorText: "Invalid Email or Password",
		Version: util.Version,
	}

	// ErrInvitationCodeFrequently error invitation code too frequently
	ErrInvitationCodeFrequently = &util.ErrResponse{
		HTTPStatusCode: http.StatusOK,
		StatusText: "Invation Code Fail",
		ErrorText: "Request Invitation Code Too Frequently",
		Version: util.Version,
	}

	// ErrInvitationCodeWrong error invation code 
	ErrInvitationCodeWrong = &util.ErrResponse{
		HTTPStatusCode: http.StatusOK,
		StatusText: "Invation Code Fail",
		ErrorText: "Invitation Code Wrong",
		Version: util.Version,
	}
	// ErrEmailUsed error email used
	ErrEmailUsed = &util.ErrResponse{
		HTTPStatusCode: http.StatusOK,
		StatusText: "Register Fail",
		ErrorText: "Email Used",
		Version: util.Version,
	}
	// ErrInsufficientAuthority error insufficient authority
	ErrInsufficientAuthority=&util.ErrResponse{
		HTTPStatusCode: http.StatusOK,
		StatusText: "Insufficient Authority",
		ErrorText: "Insufficient Authority",
		Version: util.Version,
	}
)