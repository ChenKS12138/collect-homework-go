package model

import (
	"time"

	"github.com/go-pg/pg/v10/orm"
)

// Admin admin
type Admin struct {
	ID string `json:"id,omitempty" pg:"id,pk,type:uuid"`
	IsSuperAdmin bool `json:"isSuperAdmin,omitempty" pg:"is_super_admin"`
	Name string `json:"name,omitempty" pg:"name"`
	Email string `json:"email,omitempty" pg:"email"`
	Password string `json:"password,omitempty" pg:"password"`
	CreateAt time.Time `json:"createAt,omitempty" pg:"create_at"`
	UpdateAt time.Time `json:"updateAt,omitempty" pg:"update_at"`
}

// BeforeUpdate before update
func (a *Admin) BeforeUpdate(db orm.DB) error {
	a.UpdateAt = time.Now();
	return nil;
}
