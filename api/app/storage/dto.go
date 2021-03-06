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

// DownloadDto download dto
type DownloadDto struct {
	ID string `json:"id"`
}

func (d *DownloadDto)validate() error {
	err := &validation.Errors{
		"id":validation.Validate(d.ID,validation.Required,is.UUIDv4),
	}
	return err.Filter()
}

// FileCountDto file count dto
type FileCountDto struct {
	ID string `json:"id"`
}

func (f *FileCountDto)validate() error {
	err := &validation.Errors{
		"id":validation.Validate(f.ID,validation.Required,is.UUIDv4),
	}
	return err.Filter()
}

// FileListDto file list dto
type FileListDto struct {
	ID string `json:"id"`
}

func (f *FileListDto)validate()error {
	err := &validation.Errors{
		"id":validation.Validate(f.ID,validation.Required,is.UUIDv4),
	}
	return err.Filter()
}

// ProjectSizeDto project size dto
type ProjectSizeDto struct {
	ID string `json:"id"`
}

func (p * ProjectSizeDto)validate() error {
	err := &validation.Errors{
		"id":validation.Validate(p.ID,validation.Required,is.UUIDv4),
	}
	return err.Filter()
}

// DownloadSelectivelyDto download dto
type DownloadSelectivelyDto struct {
	ID string `json:"id"`
	Code string `json:"code"`
}

func (d *DownloadSelectivelyDto)validate() error {
	err := &validation.Errors{
		"id":validation.Validate(d.ID,validation.Required,is.UUIDv4),
		"code":validation.Validate(d.Code,validation.Required,is.Hexadecimal,is.UpperCase),
	}
	return err.Filter()
}

// RawFileDto raw file dto
type RawFileDto struct {
	ID string `json:"id"`
	FileName string `json:"id"`
}

func (r *RawFileDto)validate() error {
	err := &validation.Errors{
		"id":validation.Validate(r.ID,validation.Required,is.UUIDv4),
		"filename":validation.Validate(r.FileName,validation.Required),
	}
	return err.Filter()
}