package main

import (
	data_utils "github.com/juxuny/data-utils"
	"github.com/juxuny/data-utils/cmd/data-utils/weibo"
	"github.com/juxuny/data-utils/log"
	"github.com/juxuny/data-utils/model"
	"github.com/spf13/cobra"
)

var weiboFlag = struct {
	StartUrl  string
	UseEnv    bool
	DbHost    string
	DbPort    int
	DbName    string
	DbUser    string
	DbPwd     string
	Driver    string
	BatchSize int
}{}

func getWeiboDbConfigFromCommandLineArgs() model.Config {
	config := model.Config{
		DbHost:     weiboFlag.DbHost,
		DbPort:     weiboFlag.DbPort,
		DbUser:     weiboFlag.DbUser,
		DbPassword: weiboFlag.DbPwd,
		DbName:     weiboFlag.DbName,
		DbDebug:    globalFlag.Verbose,
	}
	return config
}

var weiboCmd = &cobra.Command{
	Use:   "weibo",
	Short: "weibo",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		var config data_utils.QueueConfig
		if weiboFlag.UseEnv {
			config.DbConfig, err = model.GetEnvConfig()
			if err != nil {
				log.Fatal(err)
			}
		} else {
			config.DbConfig = getWeiboDbConfigFromCommandLineArgs()
		}
		config.BatchSize = weiboFlag.BatchSize
		config.DriverType = data_utils.QueueDriverType(weiboFlag.Driver)
		queue := data_utils.NewQueue(config)
		queue.RegisterHandler(weibo.NewHandlerBuilder(config))
		if weiboFlag.StartUrl != "" {
			job := data_utils.Job{
				JobType:  model.JobTypeWeibo,
				MetaData: weibo.MetaData{Url: weiboFlag.StartUrl}.ToJson(),
			}
			if err := queue.Enqueue(job); err != nil {
				log.Fatal(err)
			}
		}
		if err := queue.StartDaemon(); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	initGlobalFlat(weiboCmd)
	weiboCmd.PersistentFlags().StringVar(&weiboFlag.StartUrl, "start-url", "", "the first link for crawl")
	weiboCmd.PersistentFlags().BoolVarP(&weiboFlag.UseEnv, "use-env", "e", false, "get arguments from environment variables")
	weiboCmd.PersistentFlags().StringVar(&weiboFlag.Driver, "driver", data_utils.QueueDriverTypeMysql.ToString(), "queue driver")
	weiboCmd.PersistentFlags().IntVar(&weiboFlag.BatchSize, "batch-size", 1, "batch size")

	// database
	weiboCmd.PersistentFlags().StringVar(&weiboFlag.DbHost, "db-host", "127.0.0.1", "database host")
	weiboCmd.PersistentFlags().IntVar(&weiboFlag.DbPort, "db-port", 3306, "database port")
	weiboCmd.PersistentFlags().StringVar(&weiboFlag.DbUser, "db-user", "root", "user")
	weiboCmd.PersistentFlags().StringVar(&weiboFlag.DbPwd, "db-pwd", "", "password for database")
	weiboCmd.PersistentFlags().StringVar(&weiboFlag.DbName, "db-name", "", "database schema name")
	rootCmd.AddCommand(weiboCmd)
}
