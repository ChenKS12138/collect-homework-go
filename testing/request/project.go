package request

import (
	"net/http"

	"github.com/ChenKS12138/collect-homework-go/model"
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
	Labels []string `json:"labels"`
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
	Labels []string `json:"labels"`
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

// ProjectFileList project file list
func ProjectFileList(url string,token string,projectID string) (*struct {
	BasicResponse
	Data struct {
		Files []struct{
			Name string `json:"name"`
			Code string `json:"code"`
		} `json:"files"`
	} `json:"data"`
},error) {
	response := &struct {
		BasicResponse
		Data struct {
			Files []struct{
				Name string `json:"name"`
				Code string `json:"code"`
			} `json:"files"`
		} `json:"data"`
	}{}
	header := &http.Header{}
	header.Set("Authorization","Bearer "+token)
	err := GetRequest(url,header,map[string]string{
		"id":projectID,
	},response)
	return response,err
}