package storage

import (
	"collect-homework-go/util"
	"net/http"
)

var (
	// ErrFileSize file too large
	ErrFileSize = &util.ErrResponse{
		HTTPStatusCode: http.StatusOK,
		StatusText: "Storage Error",
		ErrorText: "File Too Large",
	}
	// ErrFileNamePattern file name pattern error
	ErrFileNamePattern = &util.ErrResponse{
		HTTPStatusCode: http.StatusOK,
		StatusText: "Storage Error",
		ErrorText: "File Name Pattern Error",
	}
	// ErrFileNameExtensions file name extensions error
	ErrFileNameExtensions = &util.ErrResponse{
		HTTPStatusCode: http.StatusOK,
		StatusText: "Storage Error",
		ErrorText: "File Extensions Name Error",
	}
	// ErrFileSecret file secret wrong
	ErrFileSecret = &util.ErrResponse{
		HTTPStatusCode: http.StatusOK,
		StatusText: "Storage Error",
		ErrorText: "File Secret Wrong",
	}
	// ErrProjectNotExist project not exist
	ErrProjectNotExist = &util.ErrResponse{
		HTTPStatusCode: http.StatusOK,
		StatusText: "Storage Error",
		ErrorText: "File Project Not Exist",
	}
	// ErrDownloadForbidden downlaod forbidden
	ErrDownloadForbidden = &util.ErrResponse{
		HTTPStatusCode: http.StatusOK,
		StatusText: "Storage Error",
		ErrorText: "File Download Forbidden",
	}
	// ErrProjectPremissionDenied project permission denied
	ErrProjectPremissionDenied = &util.ErrResponse{
		HTTPStatusCode: http.StatusOK,
		StatusText: "Storage Error",
		ErrorText: "Project Permission Denied",
	}
)