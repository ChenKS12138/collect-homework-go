package storage

import (
	"mime/multipart"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

// UploadDto upload dto
type UploadDto struct {
	File multipart.File `json:"file"`
	FileHeader *multipart.FileHeader `json:"fileHeader"`
	Secret string `json:"secret"`
	ProjectID string `json:"projectId"`
}

func (u *UploadDto)validate() error {
	err := &validation.Errors{
		"secret":validation.Validate(u.Secret,validation.Required),
		"projectId":validation.Validate(u.ProjectID,validation.Required,is.UUIDv4),
		"file":validation.Validate(u.File,validation.Required),
		"fileHeader":validation.Validate(u.FileHeader,validation.Required),
	}
	return err.Filter()
}