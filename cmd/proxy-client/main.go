package main

import (
	"github.com/juxuny/data-utils/model"
	"github.com/juxuny/log-server/log"
	"github.com/spf13/cobra"
	"os"
)

var logger = log.NewLogger("[su]")

var (
	rootCmd = &cobra.Command{
		Use:   "proxy-client",
		Short: "proxy-client",
	}
)

var globalFlag = struct {
	Verbose bool
	DbHost  string
	DbPort  int
	DbUser  string
	DbPwd   string
	DbName  string
}{}

func getDbConfigFromCommandLineArgs() model.Config {
	config := model.Config{
		DbHost:     globalFlag.DbHost,
		DbPort:     globalFlag.DbPort,
		DbUser:     globalFlag.DbUser,
		DbPassword: globalFlag.DbPwd,
		DbName:     globalFlag.DbName,
		DbDebug:    globalFlag.Verbose,
	}
	return config
}

func initGlobalFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVarP(&globalFlag.Verbose, "verbose", "v", false, "display debug output")

	// database
	cmd.PersistentFlags().StringVar(&globalFlag.DbHost, "db-host", "127.0.0.1", "database host")
	cmd.PersistentFlags().IntVar(&globalFlag.DbPort, "db-port", 3306, "database port")
	cmd.PersistentFlags().StringVar(&globalFlag.DbUser, "db-user", "root", "user")
	cmd.PersistentFlags().StringVar(&globalFlag.DbPwd, "db-pwd", "", "password for database")
	cmd.PersistentFlags().StringVar(&globalFlag.DbName, "db-name", "", "database schema name")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		logger.Error(err)
		os.Exit(-1)
	}
}
