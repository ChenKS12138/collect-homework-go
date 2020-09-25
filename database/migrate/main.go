package migrate

import (
	"context"
	"log"

	"collect-homework-go/database"

	"github.com/go-pg/migrations/v8"
)

var (
	ctx  = context.Background()
)


// Migrate migrate
func Migrate(args []string){
	db,err := database.DBConn()
	if err != nil {
		log.Fatal(err)
	}

	oldVersion ,newVersion ,err := migrations.Run(db,args...)
	if err != nil {
		log.Fatal(err)
	}
	if newVersion != oldVersion {
		log.Printf("migrated from version %d to %d\n", oldVersion, newVersion)
	} else {
		log.Printf("version is %d\n", oldVersion)
	}

	if err != nil {
		log.Fatal(err)
	}
}

// Init init
func Init(){
	db,err := database.DBConn()
	if err != nil {
		log.Fatal(err)
	}

	oldVersion ,newVersion ,err := migrations.Run(db,"init")
	if err != nil {
		log.Fatal(err)
	}
	if newVersion != oldVersion {
		log.Printf("migrated from version %d to %d\n", oldVersion, newVersion)
	} else {
		log.Printf("version is %d\n", oldVersion)
	}

	if err != nil {
		log.Fatal(err)
	}
}

// Reset runs reverts all migrations to version 0 and then applies all migrations to latest
func Reset() {
	db, err := database.DBConn()
	if err != nil {
		log.Fatal(err)
	}

	version, err := migrations.Version(db)
	if err != nil {
		log.Fatal(err)
	}

	for version != 0 {
		oldVersion, newVersion, err := migrations.Run(db, "down")
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("migrated from version %d to %d\n", oldVersion, newVersion)
		version = newVersion
	}
	if err != nil {
		log.Fatal(err)
	}
}
