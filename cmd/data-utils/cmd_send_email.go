package main

import (
	"github.com/ghodss/yaml"
	"github.com/juxuny/data-utils/email"
	"github.com/juxuny/data-utils/log"
	"github.com/spf13/cobra"
	"io/ioutil"
	"math/rand"
	"strings"
	"time"
)

var emailFlag struct {
	SenderData string
	DataFile   string
	EmailList  []string
	Config     string
	Host       string
	User       string
	Password   string
	Ssl        bool
	BatchSize  int
	Delay      int
}

var emailCmd = &cobra.Command{
	Use: "email",
	Run: func(cmd *cobra.Command, args []string) {
		var clientConfigList []email.ClientConfig
		if emailFlag.SenderData != "" {
			data, err := ioutil.ReadFile(emailFlag.SenderData)
			if err != nil {
				log.Fatal(err)
			}
			l := strings.Split(string(data), "\n")
			randIndex := rand.Perm(len(l))
			for _, index := range randIndex {
				item := l[index]
				sender := strings.Split(item, "----")
				if len(sender) == 2 {
					clientConfigList = append(clientConfigList, email.ClientConfig{
						User:        strings.Trim(sender[0], " \n\r\t-"),
						DisplayName: "",
						Password:    strings.Trim(sender[1], " \n\r\t-"),
						Host:        emailFlag.Host,
						Ssl:         emailFlag.Ssl,
					})
				}
			}
		} else if emailFlag.User != "" && emailFlag.Password != "" {
			clientConfigList = append(clientConfigList, email.ClientConfig{
				User:     emailFlag.User,
				Password: emailFlag.Password,
				Host:     emailFlag.Host,
				Ssl:      emailFlag.Ssl,
			})
		}
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

		senderIndex := 0
		if clientConfigList == nil || len(clientConfigList) == 0 {
			log.Fatal("not sender config")
			return
		}
		for i := 0; i < len(receiverList); i += emailFlag.BatchSize {
			log.Info("used:", clientConfigList[senderIndex].User)
			client := email.NewClient(email.ClientConfig{
				User:        clientConfigList[senderIndex].User,
				DisplayName: clientConfigList[senderIndex].DisplayName,
				Password:    clientConfigList[senderIndex].Password,
				Host:        clientConfigList[senderIndex].Host,
				Ssl:         clientConfigList[senderIndex].Ssl,
			})
			var batch = make([]string, 0)
			for j := 0; j < len(receiverList) && j < emailFlag.BatchSize; j++ {
				batch = append(batch, receiverList[i+j])
			}
			config.To = batch
			for _, item := range batch {
				log.Info("sending: ", item)
			}
			if err := client.Send(config); err != nil {
				log.Error(err)
			}
			time.Sleep(time.Second * time.Duration(emailFlag.Delay))
			senderIndex += 1
			senderIndex %= len(clientConfigList)
		}
	},
}

func init() {
	initGlobalFlat(emailCmd)
	emailCmd.PersistentFlags().StringVar(&emailFlag.SenderData, "sender-file", "tmp/sender.list", "sender email address list")
	emailCmd.PersistentFlags().StringVar(&emailFlag.DataFile, "data-file", "", "a list of email")
	emailCmd.PersistentFlags().StringSliceVar(&emailFlag.EmailList, "email", []string{}, "email")
	emailCmd.PersistentFlags().StringVar(&emailFlag.Config, "config", "tmp/config.yaml", "email content")
	emailCmd.PersistentFlags().StringVar(&emailFlag.Host, "host", "smtp.163.com:25", "email smtp host name")
	emailCmd.PersistentFlags().StringVar(&emailFlag.User, "user", "", "user name")
	emailCmd.PersistentFlags().StringVar(&emailFlag.Password, "password", "", "password")
	emailCmd.PersistentFlags().BoolVar(&emailFlag.Ssl, "ssl", false, "enable ssl")
	emailCmd.PersistentFlags().IntVar(&emailFlag.BatchSize, "batch-size", 10, "batch size")
	emailCmd.PersistentFlags().IntVar(&emailFlag.Delay, "delay", 1, "delay second")
	rootCmd.AddCommand(emailCmd)
}
