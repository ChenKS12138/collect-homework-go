package model

import (
	"time"

	"github.com/go-pg/pg/v10/orm"
)

// Submission submission
type Submission struct {
	ID string `json:"id,omitempty" pg:"id,pk,type:uuid"`
	ProjectID string `json:"projectId,omitempty" pg:"project_id"`
	Secret string `json:"secret,omitempty" pg:"secret"`
	FileName string `json:"fileName,omitempty" pg:"file_name"`
	FilePath string `json:"filePath,omitempty" pg:"file_path"`
	CreateAt time.Time `json:"createAt,omitempty" pg:"create_at"`
	UpdateAt time.Time `json:"updateAt,omitempty" pg:"update_at"`
	IP string `json:"ip,omitempty" pg:"ip"`
}

// BeforeUpdate before update
func (s *Submission) BeforeUpdate(db orm.DB) error {
	s.UpdateAt = time.Now();
	return nil;
}
