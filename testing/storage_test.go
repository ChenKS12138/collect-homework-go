package api_test

import (
	fileBytes "collect-homework-go/testing/bytes"
	"collect-homework-go/testing/service"
	"collect-homework-go/util"
	"testing"
)

func TestStorageUplaod(t *testing.T){
	userInfo := struct {
		Email string `json:"email"`
		Password string `json:"password"`
		Name string `json:"name"`
	} {
		Email:  "testStorageUpload@example.com",
		Password: "testStorageUpload",
		Name: "testStorageUpload",
	}
	_,token,err := service.AdminRegisterAndLogin(Ts.URL,userInfo.Email,userInfo.Password,userInfo.Name)
	if err != nil {
		t.Fatal(err)
	}
	_,err = service.ProjectInsert(Ts.URL,token,"test_storage_upload"+util.RandString(6),"^B\\d{8}-.{2,4}-.{2}\\d$",[]string{"doc","docx"})
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
	_,err = service.StorageUplaod(Ts.URL,util.RandString(6),projectID,"B11111111-陈陈陈-实验1.doc",fileBytes.Docx)
	if err != nil {
		t.Fatal(err)
	}
}