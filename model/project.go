package model

import (
	"time"

	"github.com/go-pg/pg/v10/orm"
)

// Project project
type Project struct {
	ID string `json:"id,omitempty" pg:"id,pk,type:uuid"`
	Name string `json:"name,omitempty" pg:"name"`
	AdminID string `json:"adminId,omitempty" pg:"admin_id"`
	FileNamePattern string `json:"fileNamePattern,omitempty" pg:"file_name_pattern"`
	FileNameExtensions []string `json:"fileNameExtensions,omitempty" pg:"file_name_extensions,array"`
	FileNameExample string `json:"fileNameExample,omitempty" pg:"file_name_example"`
	Usable bool `json:"usable" pg:"usable,notnull,use_zero,default:true"`
	CreateAt time.Time `json:"createAt,omitempty" pg:"create_at"`
	UpdateAt time.Time `json:"updateAt,omitempty" pg:"update_at"`
}

// BeforeUpdate before update
func (p *Project) BeforeUpdate(db orm.DB) error {
	p.UpdateAt = time.Now();
	return nil;
}


// ProjectWithAdminName project with admin name
type ProjectWithAdminName struct {
	tableName struct{} `pg:"projects,alias:project"`
	Project
	// Extra
	AdminName string `json:"adminName,omitempty"`
}

// ProjectWithAdminEmail project wiaht email
type ProjectWithAdminEmail struct {
	tableName struct{} `pg:"projects,alias:project"`
	ProjectWithAdminName
	AdminEmail string `json:"adminEmail,omitempty"`
}