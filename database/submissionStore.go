package database

import (
	"errors"
	"time"

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

// SelectAllFile select all file
func (s *SubmissionStore)SelectAllFile(projectID string) (*[]model.Submission,error) {
	submissions := &[]model.Submission{}
	err := s.db.Model(submissions).
		Where(`"project_id" = ?`,projectID).
		ColumnExpr(`"file_name","min"("create_at") AS "create_at" `).
		Group(`file_name`).
		Order(`create_at DESC`).
		Select()
	if err == pg.ErrNoRows {
		return nil,nil
	}
	return submissions,err
}

// SelectFilePathMap select file path map
func (s *SubmissionStore)SelectFilePathMap(projectID string) (*map[string]string ,error) {
	submissions := &[]model.Submission{}
	err := s.db.Model(submissions).
		Where(`project_id = ?`,projectID).
		DistinctOn(`"file_name","file_path"`).
		ColumnExpr(`"file_name","file_path"`).
		Select()
	if err == pg.ErrNoRows {
		return nil,nil
	}
	if err != nil {
		return nil,err
	}
	filePathMap :=make(map[string]string);
	for _,submission := range(*submissions) {
		filePathMap[submission.FileName] = submission.FilePath
	}
	return &filePathMap,nil
}

// SelectFileLastModifyTimeMap select file last modify time map
func (s*SubmissionStore)SelectFileLastModifyTimeMap(projectID string)(*map[string](time.Time) ,error) {
	submissions := &[]model.Submission{}
	err := s.db.Model(submissions).
		Where(`"project_id" = ?`,projectID).
		ColumnExpr(`"file_name","max"("update_at") AS "update_at" `).
		Group(`file_name`).
		Select()
	if err == pg.ErrNoRows {
		return nil,nil
	}
	if err != nil {
		return nil,err
	}
	filePathMap :=make(map[string](time.Time));
	for _,submission := range(*submissions) {
		filePathMap[submission.FileName] = submission.UpdateAt
	}
	if &filePathMap == nil {
		return nil,errors.New("Not File Path Map")
	}
	return &filePathMap,nil
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