package main

import (
	data_utils "github.com/juxuny/data-utils"
	"github.com/juxuny/data-utils/cmd/data-utils/weibo"
	"github.com/juxuny/data-utils/log"
	"github.com/juxuny/data-utils/model"
	"github.com/spf13/cobra"
	"time"
)

var weiboFlag = struct {
	StartUrl  string
	UseEnv    bool
	Driver    string
	BatchSize int
	RandDelay int
}{}

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
			config.DbConfig = getDbConfigFromCommandLineArgs()
		}
		config.RandDelay = time.Duration(weiboFlag.RandDelay) * time.Second
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
	initGlobalFlag(weiboCmd)
	weiboCmd.PersistentFlags().StringVar(&weiboFlag.StartUrl, "start-url", "", "the first link for crawl")
	weiboCmd.PersistentFlags().BoolVarP(&weiboFlag.UseEnv, "use-env", "e", false, "get arguments from environment variables")
	weiboCmd.PersistentFlags().StringVar(&weiboFlag.Driver, "driver", data_utils.QueueDriverTypeMysql.ToString(), "queue driver")
	weiboCmd.PersistentFlags().IntVar(&weiboFlag.BatchSize, "batch-size", 1, "batch size")
	weiboCmd.PersistentFlags().IntVar(&weiboFlag.RandDelay, "rand-delay", 5, "random delay in seconds")
	rootCmd.AddCommand(weiboCmd)
}
