package main

import "github.com/spf13/cobra"

var globalFlag = struct {
	Verbose bool
	DbHost  string
	DbPort  int
	DbUser  string
	DbPwd   string
	DbName  string
}{}

func initGlobalFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVarP(&globalFlag.Verbose, "verbose", "v", false, "display debug output")

	// database
	cmd.PersistentFlags().StringVar(&globalFlag.DbHost, "db-host", "127.0.0.1", "database host")
	cmd.PersistentFlags().IntVar(&globalFlag.DbPort, "db-port", 3306, "database port")
	cmd.PersistentFlags().StringVar(&globalFlag.DbUser, "db-user", "root", "user")
	cmd.PersistentFlags().StringVar(&globalFlag.DbPwd, "db-pwd", "", "password for database")
	cmd.PersistentFlags().StringVar(&globalFlag.DbName, "db-name", "", "database schema name")
}
