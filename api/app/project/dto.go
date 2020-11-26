package project

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

// InsertDto insert dto
type InsertDto struct {
	Name string `json:"name"`
	FileNamePattern string `json:"fileNamePattern"`
	FileNameExtensions []string `json:"fileNameExtensions"`
	FileNameExample string 	`json:"fileNameExample"`
	Labels []string `json:"labels"`
}

func (i *InsertDto)validate() error {
	err := &validation.Errors{
		"name":validation.Validate(i.Name,validation.Required,validation.NotNil),
		"fileNamePattern": validation.Validate(i.FileNamePattern,validation.NotNil),
		"fileNameExtensions":validation.Validate(i.FileNameExtensions,validation.NotNil),
		"fileNameExample":validation.Validate(i.FileNameExample,validation.NotNil),
		"labels":validation.Validate(i.Labels,validation.NotNil),
	}
	return err.Filter()
}

//UpdateDto update dto
type UpdateDto struct {
	ID string `json:"id"`
	Usable bool `json:"usable"`
	FileNamePattern string `json:"fileNamePattern"`
	FileNameExtensions []string `json:"fileNameExtensions"`
	FileNameExample string 	`json:"fileNameExample"`
	SendEmail bool `json:"sendEmail"`
	Visible bool `json:"visible"`
	Labels []string `json:"labels"`
}

func (u *UpdateDto)validate() error {
	err := &validation.Errors{
		"id":validation.Validate(u.ID,validation.Required,is.UUIDv4),
		"fileNamePattern": validation.Validate(u.FileNamePattern,validation.NotNil),
		"fileNameExtensions":validation.Validate(u.FileNameExtensions,validation.NotNil),
		"fileNameExample":validation.Validate(u.FileNameExample,validation.NotNil),
		"usable":validation.Validate(u.Usable,validation.In(true,false)),
		"labels":validation.Validate(u.Labels,validation.NotNil),
	}
	return err.Filter()
}

// DeleteDto delete dto
type DeleteDto struct {
	ID string `json:"id"`
}

func (d *DeleteDto)validate() error {
	err := &validation.Errors{
		"id":validation.Validate(d.ID,validation.Required,is.UUIDv4),
	}
	return err.Filter()
}

// RestoreDto delete dto
type RestoreDto struct {
	ID string `json:"id"`
}

func (d *RestoreDto)validate() error {
	err := &validation.Errors{
		"id":validation.Validate(d.ID,validation.Required,is.UUIDv4),
	}
	return err.Filter()
}


// FileListDto file list dto
type FileListDto struct {
	ID string `json:"id"`
}

func(f *FileListDto)validate() error {
	err := &validation.Errors{
		"id":validation.Validate(f.ID,validation.Required,is.UUIDv4),
	}
	return err.Filter()
}