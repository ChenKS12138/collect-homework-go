package model

import (
	"time"

	"github.com/go-pg/pg/v10/orm"
)

// InvitationCode invitation code
type InvitationCode struct {
	ID string `json:"id,omitempty" pg:"id,pk,type:uuid"`
	Email string `json:"email,omitempty" pg:"email"`
	Code string `json:"code,omitempty" pg:"code"`
	CreateAt time.Time `json:"createAt,omitempty" pg:"create_at"`
	UpdateAt time.Time `json:"updateAt,omitempty" pg:"update_at"`
}


// BeforeUpdate before update
func (i *InvitationCode) BeforeUpdate(db orm.DB) error {
	i.UpdateAt = time.Now();
	return nil
}
