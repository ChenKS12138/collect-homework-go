package api_test

import (
	"collect-homework-go/testing/service"
	"collect-homework-go/util"
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
	_,err = service.ProjectInsert(ts.URL,token)
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
	_,err = service.ProjectInsert(ts.URL,token)
	if err != nil {
		t.Fatal(err)
	}
	// t.Log("Test Common Admin Project Insert Pass")
}

