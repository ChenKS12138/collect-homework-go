package database

import (
	"collect-homework-go/model"

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
	if err != nil {
		return nil,err;
	}
	return submission,nil;
}

// SelectByProjectIDAndName select by project id and name
func (s *SubmissionStore)SelectByProjectIDAndName(projectID string,name string)(*model.Submission,error) {
	submission:= &model.Submission{}
	err := s.db.Model(submission).
		Where("project_id = ?",projectID).
		Where("name = ?",name).
		First()
	if err != nil {
		return nil,err
	}
	return submission,nil
}

// Insert insert
func (s *SubmissionStore)Insert(submission *model.Submission) error {
	_,err := s.db.Model(submission).Insert();
	if err != nil {
		return err
	}
	return nil
}