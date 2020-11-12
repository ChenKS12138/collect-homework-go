package api_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/ChenKS12138/collect-homework-go/api/app/storage"
	fileBytes "github.com/ChenKS12138/collect-homework-go/testing/bytes"
	"github.com/ChenKS12138/collect-homework-go/testing/service"
	"github.com/ChenKS12138/collect-homework-go/util"
	"github.com/chenhg5/collection"
)

// GET /project/
func TestProjectList(t *testing.T){
	_,_,err := service.ProjectList(Ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Test Project List Pass")
}

// GET /project/
// POST /project/insert
func TestSuperAdminProjectInsert(t *testing.T){
	_,token,err := service.AdminLogin(Ts.URL,SuperAdmin.Email,SuperAdmin.Password)
	if err != nil {
		t.Fatal(err)
	}	
	_,err = service.ProjectInsert(Ts.URL,token,"superAdminInsertProject","\\w",[]string{"doc","docx"})
	if err != nil {
		t.Fatal(err)
	}
	// t.Log("Test Super Admin Project Insert Pass")
}

// GET /project/
// POST /project/insert
func TestCommonAdminProjectInsert(t *testing.T){
	randomText := util.RandString(6)
	_,token,err := service.AdminRegisterAndLogin(Ts.URL,"common_admin_project_insert_"+randomText+"@example.com","secret","commonAdminProjectInsert")
	if err !=nil {
		t.Fatal(err)
	}
	_,err = service.ProjectInsert(Ts.URL,token,"commonAdminInsertProject","\\w",[]string{"doc","docx"})
	if err != nil {
		t.Fatal(err)
	}
	// t.Log("Test Common Admin Project Insert Pass")
}

// GET /own
// 仅超级管理员可查看所有项目
func TestProjectOwn(t *testing.T){
	_,superUserToken,err := service.AdminLogin(Ts.URL,SuperAdmin.Email,SuperAdmin.Password)
	if err != nil {
		t.Fatal(err)
	}
	commonUserProjectFileNameExample := "commonUserProject_filename_example"+util.RandString(10)
	superUserProjectFileNameExample := "superUserProject_filename_example"+util.RandString(10)

	commonUserToken,err := generateAdmin()

	_,err = service.ProjectInsert(Ts.URL,commonUserToken,commonUserProjectFileNameExample,"\\w",[]string{"doc","docx"})
	if err != nil {
		t.Fatal(err)
	}
	_,err = service.ProjectInsert(Ts.URL,superUserToken,superUserProjectFileNameExample,"\\w",[]string{"doc","docx"})
	if err != nil {
		t.Fatal(err)
	}
	_,commonUserProjectOwn,err := service.ProjectOwn(Ts.URL,commonUserToken)
	if err !=nil {
		t.Fatal(err)
	}
	_,superUserProjectOwn,err := service.ProjectOwn(Ts.URL,superUserToken)
	if err != nil {
		t.Fatal(err)
	}

	// test common user
	testCommonUserOk := true
	for _,project := range(*commonUserProjectOwn) {
		if project.FileNameExample == superUserProjectFileNameExample  {
			testCommonUserOk = false
		}
	}
	if !testCommonUserOk {
		t.Fatal(errors.New("Test Project Own Common User Fail"))
	}

	// test super user
	testSuperUserOk := false
	for _,project := range(*superUserProjectOwn){
		if project.FileNameExample == commonUserProjectFileNameExample  {
			testSuperUserOk = true
		}
	}
	if !testSuperUserOk {
		t.Fatal(errors.New("Test Project Own Super User Fail"))
	}
}

// POST /project/update
// 超级管理员允许修改其他其他用户创建的project
func TestProjectUpdate(t *testing.T){
	

	_,superUserToken,err := service.AdminLogin(Ts.URL,SuperAdmin.Email,SuperAdmin.Password)
	if err != nil {
		t.Fatal(err)
	}
	commonUserToken,err := generateAdmin()
	if err != nil {
		t.Fatal(err)
	}

	newFileNameExample := "test_project_update_filename_example"+util.RandString(6)
	oldFileNameExample := "test_project_update_filename_example"+util.RandString(6)

	_,err = service.ProjectInsert(Ts.URL,commonUserToken,oldFileNameExample,"\\w",[]string{"doc","docx"})
	if err != nil {
		t.Fatal(err)
	}

	_,projects,err := service.ProjectOwn(Ts.URL,commonUserToken)
	if err != nil {
		t.Fatal(err)
	}
	if len(*projects) != 1 {
		t.Fatal("Common User Project Insert Fail")
	}
	targetProjectID := (*projects)[0].ID

	// common user update project name
	_,err = service.ProjectUpdateName(Ts.URL,commonUserToken,targetProjectID,newFileNameExample)
	if err != nil {
		t.Fatal(err)
	}
	_,projects,err = service.ProjectOwn(Ts.URL,commonUserToken)
	commonUserUpdateProjectNameOk := false
	for _,project := range(*projects){
		if project.ID == targetProjectID &&
			project.FileNameExample == newFileNameExample {
				commonUserUpdateProjectNameOk = true
			}
	}
	if !commonUserUpdateProjectNameOk {
		t.Fatal(errors.New("Common User Update Project Name Fail"))
	}

	// super user update project name
	_,err = service.ProjectUpdateName(Ts.URL,superUserToken,targetProjectID,oldFileNameExample)
	if err != nil {
		t.Fatal(err)
	}
	_,projects,err = service.ProjectOwn(Ts.URL,superUserToken)
	superUserUpdateProjectNameOk := false
	for _,project := range(*projects){
		if project.ID == targetProjectID &&
			project.FileNameExample == oldFileNameExample {
				superUserUpdateProjectNameOk = true
			}
	}
	if !superUserUpdateProjectNameOk {
		t.Fatal(errors.New("Super User Update Project Name Fail"))
	}
}

func TestProjectDelete(t *testing.T){
	_,superUserToken,err := service.AdminLogin(Ts.URL,SuperAdmin.Email,SuperAdmin.Password)
	if err != nil {
		t.Fatal(err)
	}
	commonUserToken,err := generateAdmin()
	if err != nil {
		t.Fatal(err)
	}
	
	fileNameExample1 := "common_user_delete_1_filename_example"+util.RandString(10)
	fileNameExample2 := "common_user_delete_2_filename_example"+util.RandString(10)

	// insert
	_,err = service.ProjectInsert(Ts.URL,commonUserToken,fileNameExample1,"\\w",[]string{"doc","docx"})
	if err != nil {
		t.Fatal(err)
	}
	_,err = service.ProjectInsert(Ts.URL,commonUserToken,fileNameExample2,"\\w",[]string{"doc","docx"})
	if err != nil {
		t.Fatal(err)
	}
	_,projects,err := service.ProjectList(Ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	projectExist1 := false
	projectExist2 := false
	projectID1 := ""
	projectID2 := ""
	for _,project := range(*projects){
		switch project.FileNameExample {
		case fileNameExample1:
			projectExist1 = true
			projectID1= project.ID
		case fileNameExample2:
			projectExist2 = true
			projectID2 = project.ID
		}
	}

	if !projectExist1 || !projectExist2 {
		t.Fatal("User Insert Abnormal")
	}

	// delete
	_,err = service.ProjectDelete(Ts.URL,commonUserToken,projectID1)
	if err != nil {
		t.Fatal(err)
	}
	_,err = service.ProjectDelete(Ts.URL,superUserToken,projectID2)
	if err != nil {
		t.Fatal(err)
	}
	_,projects,err = service.ProjectList(Ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	projectExist1 = false
	projectExist2 = false
	for _,project := range(*projects){
		switch project.Name {
		case fileNameExample1:
			projectExist1 = true
		case fileNameExample2:
			projectExist2 = true
		}
	}
	if projectExist1 {
		t.Fatal(errors.New("User Delete Abnormal"))
	}

	if projectExist2 {
		t.Fatal(errors.New("User Delete Abnormal (Super User)"))
	}
}


func TestProjectFileList(t *testing.T){
	projectNames := []string{ "B11111111-陈陈陈-实验1.doc", "B11111112-陈陈-实验1.doc"}
	_,superUserToken,err := service.AdminLogin(Ts.URL,SuperAdmin.Email,SuperAdmin.Password)
	if err != nil {
		t.Fatal(err)
	}
	commonUserToken,err := generateAdmin()
	if err != nil {
		t.Fatal(err)
	}
	commonUserToken2,err := generateAdmin()
	if err != nil {
		t.Fatal(err)
	}

	_,err = service.ProjectInsert(Ts.URL,commonUserToken,"test_storage_upload_wrong_secret_"+util.RandString(6),"^B\\d{8}-.{2,4}-.{2}\\d$",[]string{"doc"})
	if err != nil {
		t.Fatal(err)
	}
	_,projects,err := service.ProjectOwn(Ts.URL,commonUserToken)
	if err != nil {
		t.Fatal(err)
	}
	if len(*projects) != 1 {
		t.Fatal(errors.New("Project Insert Abnormal"))
	}
	projectID := (*projects)[0].ID
	_,err = service.StorageUpload(Ts.URL,util.RandString(6),projectID,projectNames[0],fileBytes.Docx)
	if err != nil{
		t.Fatal(err)
	}
	_,err = service.StorageUpload(Ts.URL,util.RandString(6),projectID,projectNames[1],fileBytes.Docx)
	if err != nil{
		t.Fatal(err)
	}

	// common user
	_,filelist,err := service.StorageFileList(Ts.URL,commonUserToken,projectID)
	if filelist== nil || 
		!collection.Collect(filelist).Contains(projectNames[0]) ||
		!collection.Collect(filelist).Contains(projectNames[1]) || 
		collection.Collect(filelist).Count()!=2 {
		t.Fatal(errors.New("Test Storage File List Fail (Common User)"))
	}

	// super user
	_,filelist,err = service.StorageFileList(Ts.URL,superUserToken,projectID)
	if filelist== nil || 
	!collection.Collect(filelist).Contains(projectNames[0]) ||
	!collection.Collect(filelist).Contains(projectNames[1]) || 
	collection.Collect(filelist).Count()!=2 {
		t.Fatal(errors.New("Test Storage File List Fail (Super User)"))
	}

	// common user2
	ok,_,err := service.StorageFileList(Ts.URL,commonUserToken2,projectID)
	if ok || !strings.Contains(err.Error(),storage.ErrProjectPremissionDenied.ErrorText) {
		t.Fatal(errors.New("Test Storage File List Fail (Super User2)"))
	}
}