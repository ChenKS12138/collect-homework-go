package request

import (
	"github.com/ChenKS12138/collect-homework-go/model"
	"net/http"
)

// ProjectList project list
func ProjectList(url string)(*struct{
	BasicResponse
	Data struct {
		Project []model.ProjectWithAdminName `json:"projects"`
	} `json:"data"`
},error){
	response := &struct {
		BasicResponse
		Data struct {
			Project []model.ProjectWithAdminName `json:"projects"`
		} `json:"data"`
	}{}
	err := GetRequest(url,nil,nil,response)
	if err != nil {
		return nil,err
	}
	return response,nil
}

// ProjectOwn project own
func ProjectOwn(url string,token string)(*struct{
	BasicResponse
	Data struct {
		Project []model.ProjectWithAdminName `json:"projects"`
	} `json:"data"`
},error){
	header := &http.Header{}
	header.Set("Authorization","Bearer "+token)
	response := &struct {
		BasicResponse
		Data struct {
			Project []model.ProjectWithAdminName `json:"projects"`
		} `json:"data"`
	}{}
	err := GetRequest(url,header,nil,response)
	if err != nil {
		return nil,err
	}
	return response,nil
}

// ProjectInsert project insert
func ProjectInsert(url string,token string, insertDto *struct {
	Name string `json:"name"`
	FileNamePattern string `json:"fileNamePattern"`
	FileNameExtensions []string `json:"fileNameExtensions"`
	FileNameExample string `json:"fileNameExample"`
})(*struct {
	BasicResponse
	Data bool `json:"data"`
},error){
	header := &http.Header{}
	header.Set("Authorization","Bearer "+token)
	response := &struct{
		BasicResponse
		Data bool `json:"data"`
	}{}
	err := PostRequest(url,header,insertDto,response)
	if err != nil {
		return nil,err
	}
	return response,nil
}


// ProjectUpdate project insert
func ProjectUpdate(url string,token string, updateDto *struct {
	ID string `json:"id"`
	Usable bool `json:"usable"`
	Name string `json:"name"`
	FileNamePattern string `json:"fileNamePattern"`
	FileNameExtensions []string `json:"fileNameExtensions"`
	FileNameExample string `json:"fileNameExample"`
})(*struct {
	BasicResponse
	Data bool `json:"data"`
},error){
	header := &http.Header{}
	header.Set("Authorization","Bearer "+token)
	response := &struct{
		BasicResponse
		Data bool `json:"data"`
	}{}
	err := PostRequest(url,header,updateDto,response)
	if err != nil {
		return nil,err
	}
	return response,nil
}
