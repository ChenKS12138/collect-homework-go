package database

import (
	"github.com/ChenKS12138/collect-homework-go/model"

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
	if err == pg.ErrNoRows {
		return nil,nil
	}
	return project,err
}

// SelectAdminEmailByID select admin email by id
func (p *ProjectStore)SelectAdminEmailByID(id string)(*model.ProjectWithAdminEmail,error){
	project := &model.ProjectWithAdminEmail{}
	err := p.db.Model(project).
		Join("LEFT JOIN admins admin").
		JoinOn(`project."admin_id" = admin.id`).
		ColumnExpr(`project.*`).
		ColumnExpr(`admin."name" as admin_name,admin."email" as admin_email`).
		Where("project.\"id\" = ?",id).
		First()
	if err == pg.ErrNoRows {
		return nil,nil
	}
	return project,err
}

//SelectAllWithName select all with admin name
func (p *ProjectStore)SelectAllWithName()(*[]model.ProjectWithAdminName,error) {
	projects := &[]model.ProjectWithAdminName{}
	err := p.db.Model(projects).
		Join("LEFT JOIN admins admin").
		JoinOn(`project."admin_id" = admin."id"`).
		ColumnExpr(`project."name",project."id",project."file_name_pattern",project."file_name_extensions",project."file_name_example",project."create_at",project."update_at",project."usable",project."visible",project."send_email"`).
		ColumnExpr(`admin."name" AS admin_name`).
		Order("create_at DESC").
		Select()
	if err == pg.ErrNoRows {
		return nil,nil
	}
	return projects,err
}

// SelectByAdminID select by admin id
func (p *ProjectStore)SelectByAdminID(adminID string) (*[]model.ProjectWithAdminName,error) {
	projects := &[]model.ProjectWithAdminName{}
	err := p.db.Model(projects).
		Join("LEFT JOIN admins admin").
		JoinOn(`project."admin_id" = admin."id"`).
		Where("admin_id = ?",adminID).
		Where("usable = TRUE").
		ColumnExpr(`project."name",project."id",project."file_name_pattern",project."file_name_extensions",project."file_name_example",project."create_at",project."update_at",project."usable",project."visible",project."send_email"`).
		ColumnExpr(`admin."name" AS admin_name`).
		Order("create_at DESC").
		Select();
	if err == pg.ErrNoRows {
		return nil,nil
	}
	return projects,err
}

// SelectByAdminIDAndID select by admin id
func (p *ProjectStore)SelectByAdminIDAndID(adminID string,id string) (*model.ProjectWithAdminName,error) {
	projects := &model.ProjectWithAdminName{}
	err := p.db.Model(projects).
		Where("admin_id = ?",adminID).
		Where("id = ?",id).
		Column("id","name","file_name_pattern","file_name_extensions","file_name_example","create_at","update_at").
		First()
	if err == pg.ErrNoRows {
		return nil,nil
	}
	return projects,err
}

// SelectAllUsable select all usable
func (p *ProjectStore)SelectAllUsable() (*[]model.ProjectWithAdminName,error) {
	projects := &[]model.ProjectWithAdminName{}
	err := p.db.Model(projects).
		Where("usable = ?",true).
		Where("visible = ?",true).
		Join("LEFT JOIN admins admin").
		JoinOn(`project."admin_id" = admin."id"`).
		ColumnExpr(`project."name",project."id",project."file_name_pattern",project."file_name_extensions",project."file_name_example",project."create_at",project."update_at"`).
		ColumnExpr(`admin."name" AS admin_name`).
		Order("create_at DESC").
		Select()
	if err == pg.ErrNoRows {
		return nil,nil
	}
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