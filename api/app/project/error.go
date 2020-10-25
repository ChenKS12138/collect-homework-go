package project

import (
	"net/http"

	"github.com/ChenKS12138/collect-homework-go/util"
)

var (
	// ErrProjectPermission project permission denied
	ErrProjectPermission = &util.ErrResponse{
		HTTPStatusCode: http.StatusOK,
		StatusText: "Project Error",
		ErrorText: "Project Permission Denied",
		Version: util.Version,
	}
	// ErrProjectNotFound project not found
	ErrProjectNotFound = &util.ErrResponse{
		HTTPStatusCode: http.StatusOK,
		StatusText: "Project Error",
		ErrorText: "Project Not Found",
		Version: util.Version,
	}
)