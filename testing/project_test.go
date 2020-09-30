package api_test

import (
	"collect-homework-go/testing/service"
	"collect-homework-go/util"
	"errors"
	"testing"
)

// GET /project/
func TestProjectList(t *testing.T){
	_,_,err := service.ProjectList(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Test Project List Pass")
}

// GET /project/
// POST /project/insert
func TestSuperAdminProjectInsert(t *testing.T){
	_,token,err := service.AdminLogin(ts.URL,superAdmin.Email,superAdmin.Password)
	if err != nil {
		t.Fatal(err)
	}	
	_,err = service.ProjectInsert(ts.URL,token,"superAdminInsertProject")
	if err != nil {
		t.Fatal(err)
	}
	// t.Log("Test Super Admin Project Insert Pass")
}

// GET /project/
// POST /project/insert
func TestCommonAdminProjectInsert(t *testing.T){
	randomText := util.RandString(6)
	_,token,err := service.AdminRegisterAndLogin(ts.URL,"common_admin_project_insert_"+randomText+"@example.com","secret","commonAdminProjectInsert")
	if err !=nil {
		t.Fatal(err)
	}
	_,err = service.ProjectInsert(ts.URL,token,"commonAdminInsertProject")
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
		Email: superAdmin.Email,
		Password: superAdmin.Password,
		Name: superAdmin.Name,
	}
	commonUser := &user{
		Email: "TestProjectOwn@example.com",
		Password: "TestProjectOwn",
		Name: "commonuser_"+util.RandString(6),
	}
	_,superUserToken,err := service.AdminLogin(ts.URL,superUser.Email,superUser.Password)
	if err != nil {
		t.Fatal(err)
	}
	commonUserProjectName := "commonUserProject_"+util.RandString(6)
	superUserProjectName := "superUserProject_"+util.RandString(6)

	_,commonUserToken,err := service.AdminRegisterAndLogin(ts.URL,commonUser.Email,commonUser.Password,commonUser.Name)

	_,err = service.ProjectInsert(ts.URL,commonUserToken,commonUserProjectName)
	if err != nil {
		t.Fatal(err)
	}
	_,err = service.ProjectInsert(ts.URL,superUserToken,superUserProjectName)
	if err != nil {
		t.Fatal(err)
	}
	_,commonUserProjectOwn,err := service.ProjectOwn(ts.URL,commonUserToken)
	if err !=nil {
		t.Fatal(err)
	}
	_,superUserProjectOwn,err := service.ProjectOwn(ts.URL,superUserToken)
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
		Email: superAdmin.Email,
		Password: superAdmin.Password,
		Name: superAdmin.Name,
	}
	commonUser := &user{
		Email: "TestProjectUpdate@example.com",
		Password: "TestProjectUpdate",
		Name: "common_user_"+util.RandString(6),
	}

	_,superUserToken,err := service.AdminLogin(ts.URL,superUser.Email,superUser.Password)
	if err != nil {
		t.Fatal(err)
	}
	_,commonUserToken,err := service.AdminRegisterAndLogin(ts.URL,commonUser.Email,commonUser.Password,commonUser.Name)
	if err != nil {
		t.Fatal(err)
	}

	oldProjectName := "test_project_update_"+util.RandString(6)
	newProjectName := "test_project_update_"+util.RandString(6)

	_,err = service.ProjectInsert(ts.URL,commonUserToken,oldProjectName)
	if err != nil {
		t.Fatal(err)
	}

	_,projects,err := service.ProjectOwn(ts.URL,commonUserToken)
	if err != nil {
		t.Fatal(err)
	}
	if len(*projects) != 1 {
		t.Fatal("Common User Project Insert Fail")
	}
	targetProjectID := (*projects)[0].ID

	// common user update project name
	_,err = service.ProjectUpdateName(ts.URL,commonUserToken,targetProjectID,newProjectName)
	if err != nil {
		t.Fatal(err)
	}
	_,projects,err = service.ProjectOwn(ts.URL,commonUserToken)
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
	_,err = service.ProjectUpdateName(ts.URL,superUserToken,targetProjectID,oldProjectName)
	if err != nil {
		t.Fatal(err)
	}
	_,projects,err = service.ProjectOwn(ts.URL,superUserToken)
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
		Email: superAdmin.Email,
		Password: superAdmin.Password,
		Name: superAdmin.Name,
	}
	commonUser := &user{
		Email: "TestProjectDelete@example.com",
		Password: "TestProjectDelete",
		Name: "common_user_"+util.RandString(6),
	}

	_,superUserToken,err := service.AdminLogin(ts.URL,superUser.Email,superUser.Password)
	if err != nil {
		t.Fatal(err)
	}
	_,commonUserToken,err := service.AdminRegisterAndLogin(ts.URL,commonUser.Email,commonUser.Password,commonUser.Name)
	if err != nil {
		t.Fatal(err)
	}
	
	projectName1 := "common_user_delete_1_"+util.RandString(6)
	projectName2 := "common_user_delete_2_"+util.RandString(6)

	// insert
	_,err = service.ProjectInsert(ts.URL,commonUserToken,projectName1)
	if err != nil {
		t.Fatal(err)
	}
	_,err = service.ProjectInsert(ts.URL,commonUserToken,projectName2)
	if err != nil {
		t.Fatal(err)
	}
	_,projects,err := service.ProjectList(ts.URL)
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
	_,err = service.ProjectDelete(ts.URL,commonUserToken,projectID1)
	if err != nil {
		t.Fatal(err)
	}
	_,err = service.ProjectDelete(ts.URL,superUserToken,projectID2)
	if err != nil {
		t.Fatal(err)
	}
	_,projects,err = service.ProjectList(ts.URL)
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