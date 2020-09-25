package database

import (
	"collect-homework-go/model"

	"github.com/go-pg/pg/v10"
)

// AdminStore admin store
type AdminStore struct {
	db *pg.DB
}

// NewAdminStore new admin store
func NewAdminStore(db *pg.DB) (*AdminStore) {
	return &AdminStore{
		db: db,
	}
}

// SelectByName get admin by name
func (s *AdminStore)SelectByName(name string) (*model.Admin,error){
	admin := &model.Admin{}
	err := s.db.Model(admin).Where("name = ?",name).First();
	if err != nil {
		return nil,err
	}
	return admin,nil
}

// SelectByEmail get admin by email
func (s *AdminStore)SelectByEmail(email string) (*model.Admin,error){
	admin:= &model.Admin{};
	err := s.db.Model(admin).Where("email = ?",email).First();
	if err != nil {
		return nil,err
	}
	return admin,err
}

// Insert insert
func (s *AdminStore)Insert(admin *model.Admin) error {
	_,err := s.db.Model(admin).Insert();
	if err != nil {
		return err
	}
	return nil
}

// SelectByID select by id
func (s *AdminStore)SelectByID(id string) (*model.Admin,error){
	admin:=&model.Admin{}
	err := s.db.Model(admin).Where("id = ?",id).First();
	if err != nil {
		return nil,err
	}
	return admin,err
}