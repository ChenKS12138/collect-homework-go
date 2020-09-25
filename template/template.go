package template

import (
	"bytes"
	mailTemplate "collect-homework-go/template/mail"
	"html/template"
	"time"
)

// Registry registry
func Registry(code string,email string, t time.Time) (string,error) {
	tmpl,err := template.New("registry").Parse(mailTemplate.InvitationCodeTemplate)
	if err != nil {
		return "",err
	}
	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf,struct {
		Code string;
		Email string;
		Time string;
	}{
		Code: code,
		Email: email,
		Time: t.Format("Mon Jan 2 15:04:05 -0700 MST 2006"),
	})
	return buf.String(),nil
}

// Submission submission
func Submission(projectName string,status string,fileName string, t time.Time,ip string)(string,error){
	tmpl,err := template.New("submission").Parse(mailTemplate.SubmissionTemplate)
	if err != nil {
		return "",err
	}
	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf,struct {
		ProjectName string;
		Status string;
		FileName string;
		Time string;
		IP string;
	}{
		ProjectName: projectName,
		Status: status,
		FileName: fileName,
		Time: t.Format("Mon Jan 2 15:04:05 -0700 MST 2006"),
		IP: ip,
	})
	return buf.String(),nil
}