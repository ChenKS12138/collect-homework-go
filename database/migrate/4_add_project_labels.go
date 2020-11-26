package migrate

import (
	"fmt"

	"github.com/go-pg/migrations/v8"
)

const addProjectColumnLabels = `
ALTER TABLE IF EXISTS projects ADD COLUMN labels text[];
`
func init(){
	up := []string{
		addProjectColumnLabels,
	}
	down := []string {
		"ALTER TABLE if EXISTS projects DROP COLUMN labels;",
	}
	migrations.MustRegisterTx(func(db migrations.DB) error {
		fmt.Println("add project column labels")
		for _,query := range up {
			_,err := db.Exec(query);
			if err !=nil {
				return err
			}
		}
		return nil
	},func(db migrations.DB) error {
		fmt.Println("dropping project column labels")
		for _,query := range down {
			_,err := db.Exec(query)
			if err != nil {
				return err
			}
		}
		return nil
	})
}