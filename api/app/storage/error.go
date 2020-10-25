package storage

import (
	"net/http"

	"github.com/ChenKS12138/collect-homework-go/util"
)

var (
	// ErrFileSize file too large
	ErrFileSize = &util.ErrResponse{
		HTTPStatusCode: http.StatusOK,
		StatusText: "Storage Error",
		ErrorText: "File Too Large",
		Version: util.Version,
	}
	// ErrFileNamePattern file name pattern error
	ErrFileNamePattern = &util.ErrResponse{
		HTTPStatusCode: http.StatusOK,
		StatusText: "Storage Error",
		ErrorText: "File Name Pattern Error",
		Version: util.Version,
	}
	// ErrFileNameExtensions file name extensions error
	ErrFileNameExtensions = &util.ErrResponse{
		HTTPStatusCode: http.StatusOK,
		StatusText: "Storage Error",
		ErrorText: "File Extensions Name Error",
		Version: util.Version,
	}
	// ErrFileSecret file secret wrong
	ErrFileSecret = &util.ErrResponse{
		HTTPStatusCode: http.StatusOK,
		StatusText: "Storage Error",
		ErrorText: "File Secret Wrong",
		Version: util.Version,
	}
	// ErrProjectNotExist project not exist
	ErrProjectNotExist = &util.ErrResponse{
		HTTPStatusCode: http.StatusOK,
		StatusText: "Storage Error",
		ErrorText: "File Project Not Exist",
		Version: util.Version,
	}
	// ErrDownloadForbidden downlaod forbidden
	ErrDownloadForbidden = &util.ErrResponse{
		HTTPStatusCode: http.StatusOK,
		StatusText: "Storage Error",
		ErrorText: "File Download Forbidden",
		Version: util.Version,
	}
	// ErrProjectPremissionDenied project permission denied
	ErrProjectPremissionDenied = &util.ErrResponse{
		HTTPStatusCode: http.StatusOK,
		StatusText: "Storage Error",
		ErrorText: "Project Permission Denied",
		Version: util.Version,
	}
)