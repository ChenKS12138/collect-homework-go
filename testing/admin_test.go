package api_test

import (
	"errors"
	"log"
	"strings"
	"testing"

	"github.com/ChenKS12138/collect-homework-go/api/app/admin"
	"github.com/ChenKS12138/collect-homework-go/database"
	fileBytes "github.com/ChenKS12138/collect-homework-go/testing/bytes"
	"github.com/ChenKS12138/collect-homework-go/testing/service"
	"github.com/ChenKS12138/collect-homework-go/util"
)

// POST /admin/login
func TestSuperAdminAuth(t *testing.T) {
	_,_,err := service.AdminLogin(Ts.URL,SuperAdmin.Email,SuperAdmin.Password)

	if err != nil {
		t.Fatal(err)
	}
	// t.Log("Test Super Admin Auth Pass")
}

// POST /admin/registryCode
// POST /admin/registry
// POST /admin/login
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
	_,_,err := service.AdminRegisterAndLogin(Ts.URL,userInfo.Email,userInfo.Password,userInfo.Name)
	if err != nil {
		t.Fatal(err)
	}
	// t.Log("Test Common Admin Auth Pass")
}

// 不允许使用错误账号密码登陆
func TestAdminLoginWrong(t *testing.T){
	userInfo := struct {
		Email string `json:"email"`
		Password string `json:"password"`
		Name string `json:"name"`
	} {
		Email:  "testWrong@example.com",
		Password: "password2",
		Name: "testWrong",
	}
	_,_,err := service.AdminRegisterAndLogin(Ts.URL,userInfo.Email,userInfo.Password,userInfo.Name)
	if err != nil {
		t.Fatal(err)
	}
	ok,_,err := service.AdminLogin(Ts.URL,userInfo.Email,userInfo.Password+"123")
	
	if ok || !strings.Contains(err.Error(),admin.ErrAuthorization.ErrorText) {
		t.Fatal(errors.New("Test Admin Login Wrong Fail"))
	}
	// t.Log("Test Admin Login Wrong Pass")
}

// 不允许使用已被注册的邮箱注册
func TestAdminRegisterEmailUsed(t *testing.T){
	userInfo := struct {
		Email string `json:"email"`
		Password string `json:"password"`
		Name string `json:"name"`
	} {
		Email:  SuperAdmin.Email,
		Password: SuperAdmin.Password,
		Name: "testEmailUsed",
	}
	ok,_,err := service.AdminRegisterAndLogin(Ts.URL,userInfo.Email,userInfo.Password,userInfo.Name)
	if ok || !strings.Contains(err.Error(),admin.ErrEmailUsed.ErrorText){
		t.Fatal(errors.New("Test Admin Register Email Used Fail"))
	}
}

// 不允许同一个email频繁申请邀请码
func TestAdminInvitationCodeFrequence(t *testing.T) {
	email := "TestAdminInvitationCodeFrequence@example.com"
	_,err := service.AdminInvitationCode(Ts.URL,email)
	if err != nil {
		t.Fatal(err)
	}
	ok,err := service.AdminInvitationCode(Ts.URL,email)

	if ok || !strings.Contains(err.Error(),admin.ErrInvitationCodeFrequently.ErrorText) {
		t.Fatal("Test Admin Invitation Code Frequence Fail")
	}
}

// 不允许使用错误的邀请码注册
func TestAdminRegisterWrongCode(t *testing.T){
	userInfo := struct {
		Email string `json:"email"`
		Password string `json:"password"`
		Name string `json:"name"`
	} {
		Email:  "TestAdminRegisterWrongCode@example.com",
		Password: "password1",
		Name: "TestAdminRegisterWrongCode",
	}
	_,err := service.AdminInvitationCode(Ts.URL,userInfo.Email)
	if err != nil {
		t.Fatal(err)
	}
	
	invitationCode,err := database.Store.InvitationCode.SelectByEmail(userInfo.Email);
	if err != nil {
		t.Fatal(err)
	}
	ok,err := service.AdminRegister(Ts.URL,userInfo.Email,userInfo.Password,userInfo.Name,invitationCode.Code+"123")

	if ok || ! strings.Contains(err.Error(),admin.ErrInvitationCodeWrong.ErrorText){
		t.Fatal("Test Admin Register Wrong Code Fail")
	}
}

// GET /admin/status
func TestAdminStatus(t * testing.T){ }

// POST /admin/subToken
func TestAdminSubToken(t *testing.T){
	ok,token,_ := service.AdminLogin(Ts.URL,SuperAdmin.Email,SuperAdmin.Password)
	if !ok {
		t.Fatal("Test Admin SubToken, SuperAdmin Login Fail")
	}
	_,err := service.ProjectInsert(Ts.URL,token,"test_storage_upload"+util.RandString(6),"^B\\d{8}-.{2,4}-.{2}\\d$",[]string{"doc","docx"})
	if err != nil {
		t.Fatal(err)
	}
	_,projects,err := service.ProjectOwn(Ts.URL,token)
	if err != nil {
		t.Fatal(err)
	}
	if len(*projects) != 1 {
		t.Fatal("Project Insert Abnormal")
	}
	projectID := (*projects)[0].ID
	

	ok,subToken := service.AdminSubToken(Ts.URL,token,10,5)
	if !ok {
		t.Fatal("Test Admin SubToken, Request SubToken Fail")
	}

	ok,_,err = service.ProjectOwn(Ts.URL,subToken)
	if ok {
		log.Println(err)
		t.Fatal("Test Admin SubToken, Should Be Invalid")
	}
	projectNames := []string{ "B11111111-陈陈陈-实验1.doc", "B11111112-陈陈-实验1.doc"}

	_,err = service.StorageUpload(Ts.URL,util.RandString(6),projectID,projectNames[1],fileBytes.Docx)
	if err != nil{
		t.Fatal(err)
	}
	_,err = service.StorageDownload(Ts.URL,subToken,projectID)
	if err != nil {
		t.Fatal(err)
	}
}