package api_test

import (
	"errors"
	"testing"

	"github.com/spf13/viper"
)

// POST /admin/login
func TestSuperAdminAuth(t *testing.T) {
	response := &struct{
		BasicResponse
		Data string `json:"data"`
	}{}
	err := postRequest(ts.URL+"/admin/login",&struct{
		Email string `json:"email"`;
		Password string `json:"password"`;
	}  {
		Email:viper.GetString("SUPER_USER_EMAIL"),
		Password: viper.GetString("SUPER_USER_PASSWORD"),
	},&response)

	if err != nil {
		t.Fatal(err)
	}

	if len(response.Data) == 0 {
		t.Fatal(errors.New("Wrong JWT Token"))
	}
	t.Log("Test Super Admin Auth Pass")
}

// POST /admin/registryCode
// POST /admin/registry
// POST /admin/login
// TODO 对邮件模块的验证
func TestCommonAdminAuth(t *testing.T){
	invitaionReponse := &struct {
		BasicResponse
	}{}

	err := postRequest(ts.URL+"/admin/register",&struct {
		Email string `json:"email"`;
	} {
		Email: "test@example.com",
	},&invitaionReponse)

	if err != nil {
		t.Fatal(err)
	}

}
