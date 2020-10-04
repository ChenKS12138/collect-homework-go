package database

import (
	"collect-homework-go/model"

	"github.com/go-pg/pg/v10"
)

// InvitationCodeStore invitation code store
type InvitationCodeStore struct {
	db *pg.DB
}

// NewInvitationCodeStore new invitation code store
func NewInvitationCodeStore(db *pg.DB) (*InvitationCodeStore){
	return &InvitationCodeStore{
		db:db,
	}
}

// SelectByEmail select invitation code by email
func (s *InvitationCodeStore)SelectByEmail(email string) (*model.InvitationCode,error){
	invitationCode := &model.InvitationCode{}
	err := s.db.Model(invitationCode).
		Where("email = ?",email).
		Order("create_at DESC").
		First();
	if err == pg.ErrNoRows {
		return nil,nil
	}
	return invitationCode,err
}

// Insert invitation code insert
func (s *InvitationCodeStore)Insert(invitaionCode *model.InvitationCode) error {
	_,err := s.db.Model(invitaionCode).
		Insert()
	return err
}