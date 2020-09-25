package project

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

// InsertDto insert dto
type InsertDto struct {
	Name string `json:"name,omitempty"`
	FileNamePattern string `json:"fileNamePattern,omitempty"`
	FileNameExtensions []string `json:"fileNameExtensions,omitempty"`
	FileNameExample string 	`json:"fileNameExample,omitempty"`
}

func (i *InsertDto)validate() error {
	err := &validation.Errors{
		"name":validation.Validate(i.Name,validation.Required),
		"fileNamePattern": validation.Validate(i.FileNamePattern,validation.Required),
		"fileNameExtensions":validation.Validate(i.FileNameExtensions,validation.Required),
		"fileNameExample":validation.Validate(i.FileNameExample,validation.Required),
	}
	return err.Filter()
}

//UpdateDto update dto
type UpdateDto struct {
	ID string `json:"id"`
	Usable bool `json:"usable"`
	InsertDto
}

func (u *UpdateDto)validate() error {
	err := &validation.Errors{
		"id":validation.Validate(u.ID,validation.Required,is.UUIDv4),
		"name":validation.Validate(u.Name,validation.Required),
		"fileNamePattern": validation.Validate(u.FileNamePattern,validation.Required),
		"fileNameExtensions":validation.Validate(u.FileNameExtensions,validation.Required),
		"fileNameExample":validation.Validate(u.FileNameExample,validation.Required),
		"usable":validation.Validate(u.Usable,validation.In(true,false)),
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