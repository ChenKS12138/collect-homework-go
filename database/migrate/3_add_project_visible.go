package migrate

import (
	"fmt"

	"github.com/go-pg/migrations/v8"
)

const addProjectColumnVisible = `
ALTER TABLE IF EXISTS projects ADD COLUMN visible boolean NOT NULL DEFAULT TRUE;
`
func init(){
	up := []string{
		addProjectColumnVisible,
	}
	down := []string {
		"ALTER TABLE if EXISTS projects DROP COLUMN visible;",
	}
	migrations.MustRegisterTx(func(db migrations.DB) error {
		fmt.Println("add project column visible")
		for _,query := range up {
		_,err := db.Exec(query);
			if err !=nil {
				return err
			}
		}
		return nil
	},func(db migrations.DB) error {
		fmt.Println("dropping project column visible")
		for _,query := range down {
			_,err := db.Exec(query)
			if err != nil {
				return err
			}
		}
		return nil
	})
}