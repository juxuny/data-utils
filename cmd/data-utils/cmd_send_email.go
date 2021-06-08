package main

import (
	"github.com/ghodss/yaml"
	"github.com/juxuny/data-utils/email"
	"github.com/juxuny/data-utils/log"
	"github.com/spf13/cobra"
	"io/ioutil"
	"strings"
	"time"
)

var emailFlag struct {
	DataFile    string
	EmailList   []string
	Config      string
	Host        string
	User        string
	Password    string
	DisplayName string
	Ssl         bool
	BatchSize   int
	Delay       int
}

var emailCmd = &cobra.Command{
	Use: "email",
	Run: func(cmd *cobra.Command, args []string) {
		if emailFlag.User == "" || emailFlag.Password == "" {
			log.Fatal("user or password cannot be empty")
		}
		var receiverList []string // 收件者邮箱
		if emailFlag.DataFile != "" {
			data, err := ioutil.ReadFile(emailFlag.DataFile)
			if err != nil {
				log.Fatal(err)
			}
			receiverList = strings.Split(string(data), "\n")
		} else if len(emailFlag.EmailList) > 0 {
			receiverList = emailFlag.EmailList
		}
		for i := range receiverList {
			receiverList[i] = strings.Trim(receiverList[i], " \r")
		}

		if emailFlag.Config == "" {
			log.Fatal("config cannot be empty")
		}
		configData, err := ioutil.ReadFile(emailFlag.Config)
		if err != nil {
			log.Fatal(err)
		}
		var config email.ContentConfig
		if err := yaml.Unmarshal(configData, &config); err != nil {
			log.Fatal(err)
		}
		client := email.NewClient(email.ClientConfig{
			User:        emailFlag.User,
			DisplayName: emailFlag.DisplayName,
			Password:    emailFlag.Password,
			Host:        emailFlag.Host,
			Ssl:         emailFlag.Ssl,
		})
		for i := 0; i < len(receiverList); i += emailFlag.BatchSize {
			var batch = make([]string, 0)
			for j := 0; j < len(receiverList) && j < emailFlag.BatchSize; j++ {
				batch = append(batch, receiverList[i+j])
			}
			config.To = batch
			if err := client.Send(config); err != nil {
				log.Error(err)
			}
			time.Sleep(time.Second * time.Duration(emailFlag.Delay))
		}
	},
}

func init() {
	initGlobalFlat(emailCmd)
	emailCmd.PersistentFlags().StringVar(&emailFlag.DataFile, "data-file", "", "a list of email")
	emailCmd.PersistentFlags().StringSliceVar(&emailFlag.EmailList, "email", []string{}, "email")
	emailCmd.PersistentFlags().StringVar(&emailFlag.Config, "config", "tmp/config.yaml", "email content")
	emailCmd.PersistentFlags().StringVar(&emailFlag.Host, "host", "smtp.163.com:25", "email smtp host name")
	emailCmd.PersistentFlags().StringVar(&emailFlag.User, "user", "", "user name")
	emailCmd.PersistentFlags().StringVar(&emailFlag.DisplayName, "alias", "", "alias")
	emailCmd.PersistentFlags().StringVar(&emailFlag.Password, "password", "", "password")
	emailCmd.PersistentFlags().BoolVar(&emailFlag.Ssl, "ssl", false, "enable ssl")
	emailCmd.PersistentFlags().IntVar(&emailFlag.BatchSize, "batch-size", 10, "batch size")
	emailCmd.PersistentFlags().IntVar(&emailFlag.Delay, "delay", 1, "delay second")
	rootCmd.AddCommand(emailCmd)
}
