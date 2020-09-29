package service

import (
	"collect-homework-go/model"
	"collect-homework-go/testing/request"
	"errors"
	"log"
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
func ProjectInsert(baseURL string,token string) (ok bool,err error){
	newProject :=  &struct {
		Name string `json:"name"`
		FileNamePattern string `json:"fileNamePattern"`
		FileNameExtensions []string `json:"fileNameExtensions"`
		FileNameExample string `json:"fileNameExample"`
	} {
		Name: "操作系统实验1",
		FileNamePattern: "\\w+-\\w+",
		FileNameExtensions: []string{"doc","docx"},
		FileNameExample: "B123456-cattchen.doc",
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