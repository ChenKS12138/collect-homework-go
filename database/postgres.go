package database

import (
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pgext"
	"github.com/spf13/viper"
)

// DBConn db connectim
func DBConn() (*pg.DB,error) {
	viper.AutomaticEnv()
	db :=pg.Connect(&pg.Options{
		Network: viper.GetString("DB_NETWORK"),
		Addr: viper.GetString("DB_ADDR"),
		User: viper.GetString("DB_USER"),
		Password : viper.GetString("DB_PASSWORD"),
		Database: viper.GetString("DB_DATABASE"),
	})

	if err := checkConn(db);err !=nil {
		return nil,err
	}

	if viper.GetBool("DB_DEBUG") {
		db.AddQueryHook(pgext.DebugHook{
				// Print all queries.
				Verbose: true,
		})
	}

	Store.Admin = NewAdminStore(db)
	Store.Project = NewProjectStore(db)
	Store.Submission = NewSubmissionStore(db)
	Store.InvitationCode = NewInvitationCodeStore(db)

	return db,nil
}

func checkConn(db *pg.DB) error {
	var n int 
	_,err := db.Query(pg.Scan(&n),"SELECT 1");
	return err;
}