package main

import (
	"github.com/juxuny/data-utils/log"
	"github.com/juxuny/data-utils/model"
	"github.com/spf13/cobra"
)

var weiboInitFlag = struct {
}{}

var weiboInitCmd = &cobra.Command{
	Use: "weibo-init",
	Run: func(cmd *cobra.Command, args []string) {
		dbConfig := getDbConfigFromCommandLineArgs()
		db, err := model.Open(dbConfig)
		if err != nil {
			log.Fatal(err)
			return
		}
		defer func() {
			_ = db.Close()
		}()
		if err := db.AutoMigrate(
			&model.WeiboUser{},
			&model.WeiboFans{},
		).Error; err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	initGlobalFlag(weiboInitCmd)
	rootCmd.AddCommand(weiboInitCmd)
}
