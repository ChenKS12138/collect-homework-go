package api_test

import (
	"collect-homework-go/api/app/storage"
	fileBytes "collect-homework-go/testing/bytes"
	"collect-homework-go/testing/service"
	"collect-homework-go/util"
	"errors"
	"strings"
	"testing"
)

// POST /storage/upload
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

func TestStorageUploadWrongExtensions(t *testing.T){
	userInfo := struct {
		Email string `json:"email"`
		Password string `json:"password"`
		Name string `json:"name"`
	} {
		Email:  "testStorageUploadWrongExtensions@example.com",
		Password: "testStorageUploadWrongExtensions",
		Name: "testStorageUploadWrongExtensions",
	}
	_,token,err := service.AdminRegisterAndLogin(Ts.URL,userInfo.Email,userInfo.Password,userInfo.Name)
	if err != nil {
		t.Fatal(err)
	}
	_,err = service.ProjectInsert(Ts.URL,token,"test_storage_upload_wrong_extensions_"+util.RandString(6),"^B\\d{8}-.{2,4}-.{2}\\d$",[]string{"zip","rar"})
	if err != nil {
		t.Fatal(err)
	}
	_,projects,err := service.ProjectOwn(Ts.URL,token)
	if err != nil {
		t.Fatal(err)
	}
	if len(*projects) != 1 {
		t.Fatal(errors.New("Project Insert Abnormal"))
	}
	projectID := (*projects)[0].ID
	ok,err := service.StorageUplaod(Ts.URL,util.RandString(6),projectID,"B11111111-陈陈陈-实验1.doc",fileBytes.Docx)
	if ok || ! strings.Contains(err.Error(),storage.ErrFileNameExtensions.ErrorText){
		t.Fatal(errors.New("Test Storage Upload Wrong Wrong Extensions Fail"))
	}
}

func TestStorageUploadWrongFileName(t *testing.T){
	userInfo := struct {
		Email string `json:"email"`
		Password string `json:"password"`
		Name string `json:"name"`
	} {
		Email:  "testStorageUploadWrongFileName@example.com",
		Password: "testStorageUploadWrongFileName",
		Name: "testStorageUploadWrongFileName",
	}
	_,token,err := service.AdminRegisterAndLogin(Ts.URL,userInfo.Email,userInfo.Password,userInfo.Name)
	if err != nil {
		t.Fatal(err)
	}
	_,err = service.ProjectInsert(Ts.URL,token,"test_storage_upload_wrong_filename_"+util.RandString(6),"^B\\d{8}-.{2,4}-.{2}\\d$",[]string{"doc"})
	if err != nil {
		t.Fatal(err)
	}
	_,projects,err := service.ProjectOwn(Ts.URL,token)
	if err != nil {
		t.Fatal(err)
	}
	if len(*projects) != 1 {
		t.Fatal(errors.New("Project Insert Abnormal"))
	}
	projectID := (*projects)[0].ID
	ok,err := service.StorageUplaod(Ts.URL,util.RandString(6),projectID,"B11111111.doc",fileBytes.Docx)
	if ok || ! strings.Contains(err.Error(),storage.ErrFileNamePattern.ErrorText){
		t.Fatal(errors.New("Test Storage Upload Wrong Extensions Fail"))
	}
}

func TestStorageUploadWrongSecret(t *testing.T){
	userInfo := struct {
		Email string `json:"email"`
		Password string `json:"password"`
		Name string `json:"name"`
	} {
		Email:  "testStorageUploadWrongSecret@example.com",
		Password: "testStorageUploadWrongSecret",
		Name: "testStorageUploadWrongSecret",
	}
	_,token,err := service.AdminRegisterAndLogin(Ts.URL,userInfo.Email,userInfo.Password,userInfo.Name)
	if err != nil {
		t.Fatal(err)
	}
	_,err = service.ProjectInsert(Ts.URL,token,"test_storage_upload_wrong_secret_"+util.RandString(6),"^B\\d{8}-.{2,4}-.{2}\\d$",[]string{"doc"})
	if err != nil {
		t.Fatal(err)
	}
	_,projects,err := service.ProjectOwn(Ts.URL,token)
	if err != nil {
		t.Fatal(err)
	}
	if len(*projects) != 1 {
		t.Fatal(errors.New("Project Insert Abnormal"))
	}
	projectID := (*projects)[0].ID
	fileName := "B11111111-陈陈陈-实验1.doc"
	secret := util.RandString(6)
	_,err = service.StorageUplaod(Ts.URL,secret,projectID,fileName,fileBytes.Docx)
	if err != nil{
		t.Fatal(err)
	}
	
	ok,err := service.StorageUplaod(Ts.URL,secret+"_wrong_suffix",projectID,fileName,fileBytes.Docx)
	if ok || ! strings.Contains(err.Error(),storage.ErrFileSecret.ErrorText){
		t.Fatal("Test Storage Upload Fail (Expect ErrFileSecret)")
	}
	_,err = service.StorageUplaod(Ts.URL,secret,projectID,fileName,fileBytes.Docx)
	if err != nil{
		t.Fatal(err)
		t.Fatal("Test Storage Upload Fail (Expect Overwriting Success)")
	}
}