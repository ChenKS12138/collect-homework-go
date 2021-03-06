package service

import (
	"errors"
	"log"

	"github.com/ChenKS12138/collect-homework-go/model"
	"github.com/ChenKS12138/collect-homework-go/testing/request"
	"github.com/ChenKS12138/collect-homework-go/util"
)

// ProjectList project list
func ProjectList(baseURL string) (ok bool,projects *[]model.ProjectWithAdminName,err error){
	listResponse,err := request.ProjectList(baseURL+"/project/")
	if err != nil {
		return false,nil,err
	}
	if !listResponse.Success {
		log.Println(listResponse)
		return false,nil,errors.New(listResponse.ErrorText)
	}
	return true, &listResponse.Data.Project,nil
}

// ProjectInsert project insert
func ProjectInsert(baseURL string,token string,fileNameExample string,pattern string,extensions []string,labels []string) (ok bool,err error){
	newProject :=  &struct {
		Name string `json:"name"`
		FileNamePattern string `json:"fileNamePattern"`
		FileNameExtensions []string `json:"fileNameExtensions"`
		FileNameExample string `json:"fileNameExample"`
		Labels []string `json:"labels"`
	} {
		Name: "project_name_"+util.RandString(10),
		FileNamePattern: pattern,
		FileNameExtensions: extensions,
		FileNameExample: fileNameExample,
		Labels: labels,
	}
	_,projectsBefore,err := ProjectList(baseURL)
	if err != nil {
		log.Panicln("Project Insert Before Fail")
		return false,err
	}
	insertResponse,err := request.ProjectInsert(baseURL+"/project/insert",token,newProject)
	if err != nil {
		return false,errors.New(insertResponse.ErrorText)
	}
	if !insertResponse.Success {
		log.Println("Project Insert Fail")
		return false,err
	}
	_,projectAfter,err := ProjectList(baseURL)
	if err != nil {
		log.Println("Project Insert After Fail")
		return false,err
	}
	if len(*projectAfter) - len(*projectsBefore) != 1 {
		log.Println("Project Insert Count Strange")
		return false,errors.New("Project Insert Count Strange")
	}
	newProjectExist := false
	for _,project := range(*projectAfter) {
		// 跳过fileNameExtensions检查
		if project.Name == newProject.Name &&
			project.FileNamePattern == newProject.FileNamePattern &&
			project.FileNameExample == newProject.FileNameExample {
				newProjectExist = true
			}
	}
	if !newProjectExist {
		log.Println("Project Insert New Project Missing")
		return false,err
	}
	return true,nil
}

// ProjectOwn project own
func ProjectOwn(baseURL string,token string)(ok bool,projects *[]model.ProjectWithAdminName,err error){
	ownResponse,err := request.ProjectOwn(baseURL+"/project/own",token)
	if err != nil {
		return false,nil,err
	}
	if !ownResponse.Success {
		log.Println(ownResponse)
		return false,nil,errors.New(ownResponse.ErrorText)
	}
	return true,&ownResponse.Data.Project,nil
}

// ProjectUpdateName project name
func ProjectUpdateName(baseURL string,token string,projectID string,fileNameExample string)(ok bool,err error){
	_,projects,err := ProjectOwn(baseURL,token)
	if err != nil {
		return false,err
	}
	targetProject := &model.ProjectWithAdminName{};
	found := false
	for _,project := range(*projects){
		if project.ID == projectID {
			found = true
			targetProject.FileNameExample = fileNameExample
			targetProject.FileNameExtensions = project.FileNameExtensions
			targetProject.FileNamePattern = project.FileNamePattern
			targetProject.ID = project.ID
			targetProject.Labels = project.Labels
			// targetProject.Name = projectName
		}
	}
	if !found {
		return false,errors.New("No Such Project ID")
	}
	updateResponse,err := request.ProjectUpdate(baseURL+"/project/update",token,&struct {
		ID string `json:"id"`
		Usable bool `json:"usable"`
		Name string `json:"name"`
		FileNamePattern string `json:"fileNamePattern"`
		FileNameExtensions []string `json:"fileNameExtensions"`
		FileNameExample string `json:"fileNameExample"`
		Labels []string `json:"labels"`
	}{
		ID: targetProject.ID,
		Usable: true,
		Name: targetProject.Name,
		FileNamePattern: targetProject.FileNamePattern,
		FileNameExtensions: targetProject.FileNameExtensions,
		FileNameExample: targetProject.FileNameExample,
		Labels: []string{"label1","label2","label3"},
	})
	if err != nil {
		return false,err
	}
	if !updateResponse.Success {
		log.Println(updateResponse)
		return false,errors.New("Update Response Fail")
	}
	return true,nil
}

// ProjectDelete project delete
func ProjectDelete(baseURL string,token string,projectID string) (ok bool,err error){
	_,projects,err := ProjectOwn(baseURL,token)
	if err != nil {
		return false,err
	}
	targetProject := &model.ProjectWithAdminName{};
	found := false
	for _,project := range(*projects){
		if project.ID == projectID {
			found = true
			targetProject.FileNameExample = project.FileNameExample
			targetProject.FileNameExtensions = project.FileNameExtensions
			targetProject.FileNamePattern = project.FileNamePattern
			targetProject.ID = project.ID
			targetProject.Name = project.Name
		}
	}
	if !found {
		return false,errors.New("No Such Project ID")
	}
	updateResponse,err := request.ProjectUpdate(baseURL+"/project/update",token,&struct {
		ID string `json:"id"`
		Usable bool `json:"usable"`
		Name string `json:"name"`
		FileNamePattern string `json:"fileNamePattern"`
		FileNameExtensions []string `json:"fileNameExtensions"`
		FileNameExample string `json:"fileNameExample"`
		Labels []string `json:"labels"`
	}{
		ID: targetProject.ID,
		Usable: false,
		Name: targetProject.Name,
		FileNamePattern: targetProject.FileNamePattern,
		FileNameExtensions: targetProject.FileNameExtensions,
		FileNameExample: targetProject.FileNameExample,
		Labels: []string{},
	})
	if err != nil {
		return false,err
	}
	if !updateResponse.Success {
		log.Println(updateResponse)
		return false,errors.New("Update Response Fail")
	}
	return true,nil
}

// ProjectFileList project file list
func ProjectFileList(baseURL string,token string,projectID string)(ok bool,filelist []string,err error) {
	response,err := request.StorageFileList(baseURL+"/project/fileList",token,projectID)
	if err != nil {
		return false,nil,err
	}
	if !response.Success {
		return false,nil,errors.New(response.ErrorText)
	}
	return true,response.Data.Files,nil
}