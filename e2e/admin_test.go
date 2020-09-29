package api_test

import (
	"testing"

	"collect-homework-go/e2e/service"
)

// POST /admin/login
func TestSuperAdminAuth(t *testing.T) {
	_,_,err := service.AdminLogin(ts.URL,superAdmin.Email,superAdmin.Password)

	if err != nil {
		t.Fatal(err)
	}
	t.Log("Test Super Admin Auth Pass")
}

// POST /admin/registryCode
// POST /admin/registry
// POST /admin/login
// TODO 对邮件模块的验证
func TestCommonAdminAuth(t *testing.T){
	userInfo := struct {
		Email string `json:"email"`
		Password string `json:"password"`
		Name string `json:"name"`
	} {
		Email:  "test@example.com",
		Password: "password",
		Name: "admin",
	}
	_,_,err := service.AdminRegisterAndLogin(ts.URL,userInfo.Email,userInfo.Password,userInfo.Name)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Test Common Admin Auth Pass")
}
