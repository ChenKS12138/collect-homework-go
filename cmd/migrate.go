package cmd

import (
	"collect-homework-go/database/migrate"

	"github.com/spf13/cobra"
)

var (
	initPg bool
)

var migrateCmd = &cobra.Command{
	Use: "migrate",
	Run: func(cmd *cobra.Command,args []string){
		argsMig := args[:0]
		for _, arg := range args {
			switch arg {
			case "migrate", "--db_debug", "--reset","init":
			default:
				argsMig = append(argsMig, arg)
			}
		}
		if initPg {
			migrate.Init()
		} 
		migrate.Migrate(argsMig);
	},
}

func init() {
	RootCmd.AddCommand(migrateCmd)
	migrateCmd.Flags().BoolVar(&initPg,"init",false,"init db")
}