package migrate

import (
	"fmt"

	"github.com/go-pg/migrations/v8"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

const deleteUser = `
	TRUNCATE admins CASCADE;
`

func init(){
	viper.AutomaticEnv()
	password,_ := bcrypt.GenerateFromPassword([]byte(viper.GetString("SUPER_USER_PASSWORD")),10);
	bootstrapUsers := `
INSERT INTO admins (is_super_admin,name,email,password)
VALUES 
(TRUE,'`+viper.GetString("SUPER_USER_NAME")+`','`+viper.GetString("SUPER_USER_EMAIL")+`','`+string(password)+`');`

	up := []string{
		bootstrapUsers,
	};
	down := []string{
		deleteUser,
	}
	migrations.MustRegisterTx(func(db migrations.DB) error {
		fmt.Println("add bootstrap admins")
		for _, query := range up {
			_, err := db.Exec(query)
			if err != nil {
				return err
			}
		}
		return nil
	}, func(db migrations.DB) error {
		fmt.Println("truncate admins cascading")
		for _, query := range down {
			_, err := db.Exec(query)
			if err != nil {
				return err
			}
		}
		return nil
	})
}