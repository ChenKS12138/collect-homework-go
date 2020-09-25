package database

import (
	"collect-homework-go/model"

	"github.com/go-pg/pg/v10"
)

// ProjectStore project store
type ProjectStore struct {
	db *pg.DB
}

// NewProjectStore new project store
func NewProjectStore(db *pg.DB) (*ProjectStore) {
	return &ProjectStore{
		db:db,
	}
}

// SelectByID select by id
func (p *ProjectStore)SelectByID(id string) (*model.Project,error){
	project := &model.Project{}
	err := p.db.Model(project).
		Where("id = ?",id).
		First()
	return project,err
}

// SelectAdminEmailByID select admin email by id
func (p *ProjectStore)SelectAdminEmailByID(id string)(*model.ProjectWithAdminEmail,error){
	project := &model.ProjectWithAdminEmail{}
	err := p.db.Model(project).
		Join("LEFT JOIN admins as admin").
		JoinOn(`admin."name" as admin_name`).
		Where("id = ?",id).
		First()
	return project,err;
}


// SelectByAdminID select by admin id
func (p *ProjectStore)SelectByAdminID(adminID string) (*[]model.ProjectWithAdminName,error) {
	projects := &[]model.ProjectWithAdminName{}
	err := p.db.Model(projects).
		Where("admin_id = ?",adminID).
		Column("id","name","file_name_pattern","file_name_extensions","file_name_example","create_at","update_at").
		Select();
	return projects,err
}

// SelectAllUsable select all usable
func (p *ProjectStore)SelectAllUsable() (*[]model.ProjectWithAdminName,error) {
	projects := &[]model.ProjectWithAdminName{}
	err := p.db.Model(projects).
		Where("usable = ?",true).
		Join("LEFT JOIN admins as admin").
		ColumnExpr(`project."id",project."file_name_pattern",project."file_name_extensions",project."file_name_example",project."create_at",project."update_at"`).
		ColumnExpr(`admin."name" as admin_name`).
		JoinOn(`admin."id" = project."admin_id"`).
		Select()
	return projects,err
}

// Insert insert
func (p *ProjectStore)Insert(project *model.Project) error {
	_,err := p.db.Model(project).
		Insert()
	return err
}

// Update update
func (p *ProjectStore)Update(project *model.Project) error {
	_,err := p.db.Model(project).
		WherePK().
		Update()
	return err
}

//Delete delete
func (p *ProjectStore)Delete(id string,adminID string) error {
	_,err := p.db.Model((*model.Project)(nil)).
		WherePK().
		Where("admin_id = ?",adminID).
		Delete()
	return err
}