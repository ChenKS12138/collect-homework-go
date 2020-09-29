package project

import (
	"collect-homework-go/util"
	"net/http"
)

var (
	// ErrProjectPermission project permission denied
	ErrProjectPermission = &util.ErrResponse{
		HTTPStatusCode: http.StatusOK,
		StatusText: "Project Error",
		ErrorText: "Project Permission Denied",
	}
	// ErrProjectNotFound project not found
	ErrProjectNotFound = &util.ErrResponse{
		HTTPStatusCode: http.StatusOK,
		StatusText: "Project Error",
		ErrorText: "Project Not Found",
	}
)