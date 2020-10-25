package util

import (
	"net/http"

	"github.com/go-chi/render"
	validation "github.com/go-ozzo/ozzo-validation"
)

// DataResponse data response
type DataResponse struct {
	Data interface{} `json:"data"`
	Success bool `json:"success"`           // user-level status message
	StatusText string `json:"status"`
	Version string `json:"version"`
}

// NewDataResponse new data response
func NewDataResponse(d interface{}) (*DataResponse){
	return &DataResponse{
		Success: true,
		Data: d,
		StatusText: "ok",
		Version: Version,
	}
}

// ErrResponse renderer type for handling all sorts of errors.
type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	Success bool `json:"success"`           // user-level status message
	StatusText string `json:"status"`
	AppCode          int64             `json:"code,omitempty"`   // application-specific error code
	ErrorText        string            `json:"error,omitempty"`  // application-level error message, for debugging
	ValidationErrors validation.Errors `json:"errors,omitempty"` // user level model validation errors
	Version string `json:"version"`
}

// Render sets the application-specific error code in AppCode.
func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}


// ErrInvalidRequest returns status 422 Unprocessable Entity including error message.
func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusUnprocessableEntity,
		StatusText:     http.StatusText(http.StatusUnprocessableEntity),
		ErrorText:      err.Error(),
		Success: false,
		Version: Version,
	}
}

// ErrRender returns status 422 Unprocessable Entity rendering response error.
func ErrRender(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusUnprocessableEntity,
		StatusText:     "Error rendering response.",
		ErrorText:      err.Error(),
		Success: false,
		Version: Version,
	}
}

// ErrValidation returns status 422 Unprocessable Entity stating validation errors.
func ErrValidation( err error) render.Renderer {
	return &ErrResponse{
		Err:              err,
		HTTPStatusCode:   http.StatusUnprocessableEntity,
		StatusText:       http.StatusText(http.StatusUnprocessableEntity),
		ErrorText:        err.Error(),
		// ValidationErrors: valErr,
		Success: false,
		Version: Version,
	}
}

var (
	// ErrBadRequest return status 400 Bad Request for malformed request body.
	ErrBadRequest = &ErrResponse{HTTPStatusCode: http.StatusBadRequest, StatusText: http.StatusText(http.StatusBadRequest),Version: Version}

	// ErrNotFound returns status 404 Not Found for invalid resource request.
	ErrNotFound = &ErrResponse{HTTPStatusCode: http.StatusNotFound, StatusText: http.StatusText(http.StatusNotFound),Version:Version}

	// ErrInternalServerError returns status 500 Internal Server Error.
	ErrInternalServerError = &ErrResponse{HTTPStatusCode: http.StatusInternalServerError, StatusText: http.StatusText(http.StatusInternalServerError),Version: Version}
)
