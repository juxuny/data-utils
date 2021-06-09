package main

import (
	"fmt"
	"github.com/juxuny/data-utils/email"
	"github.com/juxuny/data-utils/log"
	"github.com/juxuny/data-utils/model"
	"github.com/spf13/cobra"
	"io/ioutil"
	"strings"
)

var importFlag = struct {
	DataFile string
}{}

var importCmd = &cobra.Command{
	Use: "import",
	Run: func(cmd *cobra.Command, args []string) {
		if importFlag.DataFile == "" {
			log.Fatal("--data-file cannot be empty")
			return
		}
		data, err := ioutil.ReadFile(importFlag.DataFile)
		if err != nil {
			log.Fatal(err)
			return
		}
		dbConfig := getAdEmailDbConfigFromCommandLineArgs()
		var emailAddressList []interface{}
		l := strings.Split(string(data), "\n")
		for _, line := range l {
			emailAddress := strings.Trim(line, " \r\n\t")
			if emailAddress == "" {
				continue
			}
			if err := email.IsValidEmail(emailAddress); err != nil {
				log.Error(err)
				continue
			}
			log.Debug(emailAddress)
			emailAddressList = append(emailAddressList, emailAddress)
		}
		statement := fmt.Sprintf("INSERT IGNORE INTO %s (email) VALUES ", model.AdEmail{}.TableName())
		holders := strings.Trim(strings.Repeat("(?),", len(emailAddressList)), ",")
		db, err := model.Open(dbConfig)
		if err != nil {
			log.Fatal(err)
			return
		}
		defer func() {
			_ = db.Close()
		}()
		if err := db.Exec(statement+holders, emailAddressList...).Error; err != nil {
			log.Error(err)
			return
		}
	},
}

func init() {
	initGlobalFlat(importCmd)
	importCmd.PersistentFlags().StringVar(&importFlag.DataFile, "data-file", "tmp/email.list", "")
	rootCmd.AddCommand(importCmd)
}
