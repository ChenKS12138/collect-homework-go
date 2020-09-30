package api_test

import (
	"collect-homework-go/testing/service"
	"collect-homework-go/util"
	"errors"
	"testing"
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
	type user struct {
		Email string `json:"email"`
		Password string `json:"password"`
		Name string `json:"name"`
	}
	superUser := &user{
		Email: SuperAdmin.Email,
		Password: SuperAdmin.Password,
		Name: SuperAdmin.Name,
	}
	commonUser := &user{
		Email: "TestProjectOwn@example.com",
		Password: "TestProjectOwn",
		Name: "commonuser_"+util.RandString(6),
	}
	_,superUserToken,err := service.AdminLogin(Ts.URL,superUser.Email,superUser.Password)
	if err != nil {
		t.Fatal(err)
	}
	commonUserProjectName := "commonUserProject_"+util.RandString(6)
	superUserProjectName := "superUserProject_"+util.RandString(6)

	_,commonUserToken,err := service.AdminRegisterAndLogin(Ts.URL,commonUser.Email,commonUser.Password,commonUser.Name)

	_,err = service.ProjectInsert(Ts.URL,commonUserToken,commonUserProjectName,"\\w",[]string{"doc","docx"})
	if err != nil {
		t.Fatal(err)
	}
	_,err = service.ProjectInsert(Ts.URL,superUserToken,superUserProjectName,"\\w",[]string{"doc","docx"})
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
		if project.Name == superUserProjectName &&
				project.AdminName == superUser.Name {
			testCommonUserOk = false
		}
	}
	if !testCommonUserOk {
		t.Fatal(errors.New("Test Project Own Common User Fail"))
	}

	// test super user
	testSuperUserOk := false
	for _,project := range(*superUserProjectOwn){
		if project.Name == commonUserProjectName &&
			project.AdminName == commonUser.Name {
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
	type user struct {
		Email string `json:"email"`
		Password string `json:"password"`
		Name string `json:"name"`
	}
	superUser := &user{
		Email: SuperAdmin.Email,
		Password: SuperAdmin.Password,
		Name: SuperAdmin.Name,
	}
	commonUser := &user{
		Email: "TestProjectUpdate@example.com",
		Password: "TestProjectUpdate",
		Name: "common_user_"+util.RandString(6),
	}

	_,superUserToken,err := service.AdminLogin(Ts.URL,superUser.Email,superUser.Password)
	if err != nil {
		t.Fatal(err)
	}
	_,commonUserToken,err := service.AdminRegisterAndLogin(Ts.URL,commonUser.Email,commonUser.Password,commonUser.Name)
	if err != nil {
		t.Fatal(err)
	}

	oldProjectName := "test_project_update_"+util.RandString(6)
	newProjectName := "test_project_update_"+util.RandString(6)

	_,err = service.ProjectInsert(Ts.URL,commonUserToken,oldProjectName,"\\w",[]string{"doc","docx"})
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
	_,err = service.ProjectUpdateName(Ts.URL,commonUserToken,targetProjectID,newProjectName)
	if err != nil {
		t.Fatal(err)
	}
	_,projects,err = service.ProjectOwn(Ts.URL,commonUserToken)
	commonUserUpdateProjectNameOk := false
	for _,project := range(*projects){
		if project.ID == targetProjectID &&
			project.Name == newProjectName {
				commonUserUpdateProjectNameOk = true
			}
	}
	if !commonUserUpdateProjectNameOk {
		t.Fatal(errors.New("Common User Update Project Name Fail"))
	}

	// super user update project name
	_,err = service.ProjectUpdateName(Ts.URL,superUserToken,targetProjectID,oldProjectName)
	if err != nil {
		t.Fatal(err)
	}
	_,projects,err = service.ProjectOwn(Ts.URL,superUserToken)
	superUserUpdateProjectNameOk := false
	for _,project := range(*projects){
		if project.ID == targetProjectID &&
			project.Name == oldProjectName {
				superUserUpdateProjectNameOk = true
			}
	}
	if !superUserUpdateProjectNameOk {
		t.Fatal(errors.New("Super User Update Project Name Fail"))
	}
}

func TestProjectDelete(t *testing.T){
	type user struct {
		Email string `json:"email"`
		Password string `json:"password"`
		Name string `json:"name"`
	}
	superUser := &user{
		Email: SuperAdmin.Email,
		Password: SuperAdmin.Password,
		Name: SuperAdmin.Name,
	}
	commonUser := &user{
		Email: "TestProjectDelete@example.com",
		Password: "TestProjectDelete",
		Name: "common_user_"+util.RandString(6),
	}

	_,superUserToken,err := service.AdminLogin(Ts.URL,superUser.Email,superUser.Password)
	if err != nil {
		t.Fatal(err)
	}
	_,commonUserToken,err := service.AdminRegisterAndLogin(Ts.URL,commonUser.Email,commonUser.Password,commonUser.Name)
	if err != nil {
		t.Fatal(err)
	}
	
	projectName1 := "common_user_delete_1_"+util.RandString(6)
	projectName2 := "common_user_delete_2_"+util.RandString(6)

	// insert
	_,err = service.ProjectInsert(Ts.URL,commonUserToken,projectName1,"\\w",[]string{"doc","docx"})
	if err != nil {
		t.Fatal(err)
	}
	_,err = service.ProjectInsert(Ts.URL,commonUserToken,projectName2,"\\w",[]string{"doc","docx"})
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
		switch project.Name {
		case projectName1:
			projectExist1 = true
			projectID1= project.ID
		case projectName2:
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
		case projectName1:
			projectExist1 = true
		case projectName2:
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