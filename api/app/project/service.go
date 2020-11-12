package project

import (
	"github.com/ChenKS12138/collect-homework-go/auth"
	"github.com/ChenKS12138/collect-homework-go/database"
	"github.com/ChenKS12138/collect-homework-go/model"
	"github.com/ChenKS12138/collect-homework-go/util"
)

// serviceOwn own
func serviceOwn(claim *auth.Claim)(dataResponse *util.DataResponse,errResponse *util.ErrResponse) {
	var err error
	projects := &[]model.ProjectWithAdminName{}

	// 超级管理员允许查看所有
	if claim.IsSuperAdmin {
		projects,err = database.Store.Project.SelectAllWithName()
	} else {
		projects,err = database.Store.Project.SelectByAdminID(claim.ID);
	}

	if err !=nil {
		return nil,util.ErrRender(err)
	}
	return util.NewDataResponse(struct {
		Projects []model.ProjectWithAdminName `json:"projects"`;
	}{
		Projects: *projects,
	}),nil
}

// serviceInsert insert
func serviceInsert(insertDto *InsertDto,claim *auth.Claim) (dataResponse *util.DataResponse,errResponse *util.ErrResponse) {
	project := &model.Project{
		AdminID: claim.ID,
		FileNameExample: insertDto.FileNameExample,
		FileNameExtensions: insertDto.FileNameExtensions,
		FileNamePattern: insertDto.FileNamePattern,
		Name: insertDto.Name,
	}
	
	err := database.Store.Project.Insert(project)
	if err != nil {
		return nil,util.ErrRender(err)
	}
	return util.NewDataResponse(true),nil
}

// serviceUpdate update
func serviceUpdate(updateDto *UpdateDto,claim *auth.Claim) (dataResponse *util.DataResponse,errResponse *util.ErrResponse)  {
	lastProject,err := database.Store.Project.SelectByID(updateDto.ID)

	if lastProject == nil {
		return nil,ErrProjectNotFound
	}

	if err !=nil {
		return nil,util.ErrRender(err)
	}

	if !claim.IsSuperAdmin && lastProject.AdminID != claim.ID {
		return nil,ErrProjectPermission
	}
	lastProject.FileNamePattern = updateDto.FileNamePattern
	lastProject.FileNameExtensions = updateDto.FileNameExtensions
	lastProject.FileNameExample = updateDto.FileNameExample
	lastProject.Usable = updateDto.Usable
	lastProject.SendEmail = updateDto.SendEmail
	lastProject.Visible = updateDto.Visible

	if err = database.Store.Project.Update(lastProject);err != nil {
		return nil,util.ErrRender(err)
	}
	return util.NewDataResponse(true),nil
}

// serviceDelete delete
func serviceDelete(deleteDto *DeleteDto,claim *auth.Claim) (dataResponse *util.DataResponse,errResponse *util.ErrResponse) {
	lastProject,err := database.Store.Project.SelectByID(deleteDto.ID)

	if lastProject == nil {
		return nil,ErrProjectNotFound
	}
	if err != nil {
		return nil,util.ErrRender(err)
	}
	if !claim.IsSuperAdmin && lastProject.AdminID != claim.ID {
		return nil,ErrProjectPermission
	}
	lastProject.Usable = false
	if err = database.Store.Project.Update(lastProject);err != nil {
		return nil,util.ErrRender(err)
	}
	return util.NewDataResponse(true),nil
}

// serviceRestore restore
func serviceRestore(restoreDto *RestoreDto,claim *auth.Claim) (dataResponse *util.DataResponse,errResponse *util.ErrResponse) {
	lastProject,err := database.Store.Project.SelectByID(restoreDto.ID)

	if lastProject == nil {
		return nil,ErrProjectNotFound
	}
	if err != nil {
		return nil,util.ErrRender(err)
	}
	if !claim.IsSuperAdmin {
		return nil,ErrProjectPermission
	}
	lastProject.Usable = true
	if err = database.Store.Project.Update(lastProject);err != nil {
		return nil,util.ErrRender(err)
	}
	return util.NewDataResponse(true),nil
}

func serviceList()(dataResponse *util.DataResponse,errResponse *util.ErrResponse) {
	projects,err := database.Store.Project.SelectAllUsable()
	if err != nil {
		return nil, util.ErrRender(err)
	}
	return util.NewDataResponse(struct {
		Projects []model.ProjectWithAdminName `json:"projects"`;
	}{
		Projects: *projects,
	}),nil
}

func serviceFileList(fileListDto *FileListDto,claim *auth.Claim) (dataResponse *util.DataResponse,errResponse *util.ErrResponse) {
	type StorageFile struct {
		Name string `json:"name"`
		Seq int `json:"seq"`
	}
	if !claim.IsSuperAdmin {
		project,err := database.Store.Project.SelectByAdminIDAndID(claim.ID,fileListDto.ID)
		if err != nil {
			return nil,util.ErrRender(err)
		}
		if project == nil {
			return nil,ErrProjectPermission
		}
	}
	
	files,err := database.Store.Submission.SelectAllFile(fileListDto.ID)
	if err != nil {
		return nil,util.ErrRender(err);
	}
	storageFiles := []StorageFile{}
	for index,file := range(*files){
		storageFiles = append(storageFiles,StorageFile{
			Name: file.FileName,
			Seq: index,
		})
	}

	return util.NewDataResponse(struct{
		Files []StorageFile `json:"files"`
	}{
		Files: storageFiles,
	}),nil
}