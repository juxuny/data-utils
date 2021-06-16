package main

import "github.com/juxuny/data-utils/model"

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
