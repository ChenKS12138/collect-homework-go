package database

import (
	"github.com/ChenKS12138/collect-homework-go/model"

	"github.com/go-pg/pg/v10"
)

// SubmissionStore submission store
type SubmissionStore struct {
	db *pg.DB
}

// NewSubmissionStore new submission store
func NewSubmissionStore(db *pg.DB) (*SubmissionStore){
	return &SubmissionStore{
		db:db,
	}
}

// SelectByProjectID select by project id
func (s *SubmissionStore)SelectByProjectID(projectID string) (*model.Submission,error) {
	submission := &model.Submission{};
	err := s.db.Model(submission).
		Where("project_id = ?",projectID).
		First();
	if err == pg.ErrNoRows {
		return nil,nil
	}
	return submission,err
}

// SelectCountByProjectID select count by project id
func (s *SubmissionStore)SelectCountByProjectID(projectID string)(int,error){
	submissions := &[]model.Submission{}
	count,err := s.db.Model(submissions).
		Where("project_id = ?",projectID).
		DistinctOn(`submission."file_name"`).
		Count()
	if err == pg.ErrNoRows {
		return 0,nil
	}
	return count,err
}

// SelectByProjectIDAndName select by project id and name
func (s *SubmissionStore)SelectByProjectIDAndName(projectID string,name string)(*model.Submission,error) {
	submission:= &model.Submission{}
	err := s.db.Model(submission).
		Where("project_id = ?",projectID).
		Where("file_name = ?",name).
		Order("create_at DESC").
		First()
	if err == pg.ErrNoRows {
		return nil,nil
	}
	return submission,err
}

// Insert insert
func (s *SubmissionStore)Insert(submission *model.Submission) error {
	_,err := s.db.Model(submission).Insert();
	return err
}

// SelectCountByAdminID count by admin id
func (s *SubmissionStore)SelectCountByAdminID(adminID string) (int ,error){
	count,err := s.db.Model(&model.Submission{}).
		Join("LEFT JOIN projects").
		JoinOn(`submission."project_id" = projects."id"`).
		DistinctOn(`submission."file_name"`).
		Where(`projects."admin_id" = ?`,adminID).
		Where(`projects."usable" = TRUE`).
		Count()
	if err == pg.ErrNoRows {
		return 0,nil
	}
	return count,err
}