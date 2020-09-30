package api_test

import (
	"collect-homework-go/api/app/storage"
	fileBytes "collect-homework-go/testing/bytes"
	"collect-homework-go/testing/service"
	"collect-homework-go/util"
	"errors"
	"strings"
	"testing"

	"github.com/chenhg5/collection"
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

func TestStorageFileCount(t *testing.T){
	type user struct {
		Email string `json:"email"`
		Password string `json:"password"`
		Name string `json:"name"`
	}
	superUserInfo := &user{
		Email: SuperAdmin.Email,
		Password: SuperAdmin.Password,
		Name: SuperAdmin.Name,
	}
	commonUserInfo := &user{
		Email:  "testStorageFileCount@example.com",
		Password: "testStorageFileCount",
		Name: "testStorageFileCount",
	}

	_,superUserToken,err := service.AdminLogin(Ts.URL,superUserInfo.Email,superUserInfo.Password)
	if err != nil {
		t.Fatal(err)
	}
	_,commonUserToken,err := service.AdminRegisterAndLogin(Ts.URL,commonUserInfo.Email,commonUserInfo.Password,commonUserInfo.Name)
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
	_,err = service.StorageUplaod(Ts.URL,util.RandString(6),projectID,"B11111111-陈陈陈-实验1.doc",fileBytes.Docx)
	if err != nil{
		t.Fatal(err)
	}
	_,err = service.StorageUplaod(Ts.URL,util.RandString(6),projectID,"B11111112-陈陈-实验1.doc",fileBytes.Docx)
	if err != nil{
		t.Fatal(err)
	}

	_,count,err := service.StorageFileCount(Ts.URL,commonUserToken,projectID)
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Fatal(errors.New("Test Stroage File Count Fail"))
	}

	_,count,err = service.StorageFileCount(Ts.URL,superUserToken,projectID)
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Fatal(errors.New("Test Storage File Count Fail (Super User)"))
	}

}

func TestStorageFileList(t *testing.T){
	type user struct {
		Email string `json:"email"`
		Password string `json:"password"`
		Name string `json:"name"`
	}
	superUserInfo := &user{
		Email: SuperAdmin.Email,
		Password: SuperAdmin.Password,
		Name: SuperAdmin.Name,
	}
	commonUserInfo := &user{
		Email:  "testStorageFileList@example.com",
		Password: "testStorageFileList",
		Name: "testStorageFileList",
	}
	commonUserInfo2 := &user{
		Email:  "testStorageFileList2@example.com",
		Password: "testStorageFileList2",
		Name: "testStorageFileList2",
	}
	projectNames := []string{ "B11111111-陈陈陈-实验1.doc", "B11111112-陈陈-实验1.doc"}
	_,superUserToken,err := service.AdminLogin(Ts.URL,superUserInfo.Email,superUserInfo.Password)
	if err != nil {
		t.Fatal(err)
	}
	_,commonUserToken,err := service.AdminRegisterAndLogin(Ts.URL,commonUserInfo.Email,commonUserInfo.Password,commonUserInfo.Name)
	if err != nil {
		t.Fatal(err)
	}
	_,commonUserToken2,err := service.AdminRegisterAndLogin(Ts.URL,commonUserInfo2.Email,commonUserInfo2.Password,commonUserInfo2.Name)

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
	_,err = service.StorageUplaod(Ts.URL,util.RandString(6),projectID,projectNames[0],fileBytes.Docx)
	if err != nil{
		t.Fatal(err)
	}
	_,err = service.StorageUplaod(Ts.URL,util.RandString(6),projectID,projectNames[1],fileBytes.Docx)
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