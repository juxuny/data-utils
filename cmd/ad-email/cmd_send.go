package main

import (
	"fmt"
	"github.com/ghodss/yaml"
	data_utils "github.com/juxuny/data-utils"
	"github.com/juxuny/data-utils/email"
	"github.com/juxuny/data-utils/log"
	"github.com/juxuny/data-utils/model"
	"github.com/juxuny/env"
	"github.com/juxuny/gomail"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.org/x/net/proxy"
	"io/ioutil"
	"math/rand"
	"strings"
	"time"
)

var sendFlag = struct {
	SenderFile string
	ConfigFile string
	BatchSize  int
	Count      int // 发送多少批
	Delay      int
	Email      []string
	CC         []string // 额外要发送的其它邮箱

	Host     string
	Ssl      bool
	ProxyUrl string
}{}

func getSenderClientConfigList(fileName string) (ret []email.ClientConfig, err error) {
	senderList, err := data_utils.GetListFromFile(fileName)
	if err != nil {
		log.Error(err)
		return ret, errors.Wrap(err, "load sender email list failed")
	}
	if len(senderList) == 0 {
		return ret, errors.Wrap(err, "sender list is empty")
	}
	clientConfigList := make([]email.ClientConfig, len(senderList))
	for i, item := range senderList {
		sender := strings.Split(item, "---")
		emailAddress := strings.Trim(sender[0], " \r\t-")
		if err := email.IsValidEmail(emailAddress); err != nil {
			return ret, errors.Wrap(err, "invalid email address:"+emailAddress)
		}
		password := strings.Trim(sender[1], "\r\t- ")
		clientConfigList[i] = email.ClientConfig{
			User:        emailAddress,
			DisplayName: emailAddress,
			Password:    password,
			Host:        sendFlag.Host,
			Ssl:         sendFlag.Ssl,
			CC:          sendFlag.CC,
		}
	}
	return clientConfigList, nil
}

func enqueueEmail(db *model.DB, batchSize int) (ret model.AdEmailList, err error) {
	scope := db.Where("count = 0")
	if len(sendFlag.Email) > 0 {
		scope = scope.Where("email IN (?)", sendFlag.Email)
	}
	if err := scope.Limit(batchSize).Find(&ret).Error; err != nil {
		if !model.IsErrNoDataInDb(err) {
			return ret, errors.Wrap(err, "read table failed: ad_email")
		}
	}
	return
}

func incCount(db *model.DB, ids ...int64) error {
	if len(ids) == 0 {
		return nil
	}
	if err := db.Exec(fmt.Sprintf("UPDATE %s SET count = count + 1, updated_at = CURRENT_TIMESTAMP(), last_success_time = CURRENT_TIMESTAMP() WHERE id IN (?)", model.AdEmail{}.TableName()), ids).Error; err != nil {
		return errors.Wrap(err, "inc count failed")
	}
	return nil
}

func setFailed(db *model.DB, msg string, ids ...int64) error {
	if len(ids) == 0 {
		return nil
	}
	if err := db.Exec(fmt.Sprintf("UPDATE %s SET last_error = ?, updated_at = CURRENT_TIMESTAMP(), last_error_time = CURRENT_TIMESTAMP() WHERE id IN (?)", model.AdEmail{}.TableName()), msg, ids).Error; err != nil {
		return errors.Wrap(err, "inc count failed")
	}
	return nil
}

func parseContentConfig(fileName string) (ret email.ContentConfig, err error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Error(err)
		return ret, errors.Wrap(err, "read file failed:"+fileName)
	}
	if err := yaml.Unmarshal(data, &ret); err != nil {
		return ret, err
	}
	return ret, nil
}

var sendCmd = &cobra.Command{
	Use: "send",
	Run: func(cmd *cobra.Command, args []string) {
		// load content config
		contentConfig, err := parseContentConfig(sendFlag.ConfigFile)
		if err != nil {
			log.Fatal(err)
			return
		}
		// load content template
		clientConfigList, err := getSenderClientConfigList(sendFlag.SenderFile)
		if err != nil {
			log.Fatal(err)
			return
		}
		// rand sort
		randIndex := rand.Perm(len(clientConfigList))
		tmp := make([]email.ClientConfig, len(clientConfigList))
		for i, index := range randIndex {
			tmp[i] = clientConfigList[index]
		}
		clientConfigList = tmp
		tmp = nil
		index := 0
		dbConfig := getAdEmailDbConfigFromCommandLineArgs()
		db, err := model.Open(dbConfig)
		if err != nil {
			log.Fatal(err)
			return
		}
		defer func() {
			_ = db.Close()
		}()
		running := true
		for running {
			if sendFlag.Count > 0 {
				sendFlag.Count -= 1
			}
			data_utils.RecoverRun(func() {
				receiverList, err := enqueueEmail(db, sendFlag.BatchSize)
				if err != nil {
					log.Error(err)
					running = false
					return
				}
				if len(receiverList) == 0 {
					running = false
					return
				}
				clientConfig := clientConfigList[index]
				index = (index + 1) % len(clientConfigList)
				emailClient := email.NewClient(clientConfig)
				receiverEmailList := receiverList.GetEmailList()
				if len(receiverList) == 0 {
					log.Info("email list is empty")
					return
				}
				if sendFlag.ProxyUrl != "" {
					proxyInfo, err := email.FetchProxyInfo(sendFlag.ProxyUrl)
					if err != nil {
						log.Fatal(err)
					}
					log.Info(fmt.Sprintf("proxy address: %v:%v", proxyInfo.Ip, proxyInfo.Port))
					gomail.SetSocks5Proxy(proxyInfo.Ip, proxyInfo.Port, proxy.Auth{
						User:     env.GetString("PROXY_USER"),
						Password: env.GetString("PROXY_PASS"),
					})
				}
				contentConfig.To = receiverEmailList
				log.Info("sending, from ", clientConfig.User, " to ", receiverEmailList)
				if err := emailClient.Send(contentConfig); err != nil {
					log.Error(err)
					log.Info("send failed, from ", clientConfig.User, " to ", receiverEmailList)
					if ids := receiverList.GetIdList(); ids != nil && len(ids) > 0 {
						if err := setFailed(db, err.Error(), receiverList.GetIdList()...); err != nil {
							log.Error(err)
						}
					}
					running = false
					return
				}
				ids := receiverList.GetIdList()
				if len(ids) > 0 {
					if err := incCount(db, ids...); err != nil {
						log.Error(err)
						running = false
						return
					}
				}
			})
			if sendFlag.Count == 0 {
				break
			}
			time.Sleep(time.Second * time.Duration(sendFlag.Delay))
		}

	},
}

func init() {
	initGlobalFlat(sendCmd)
	sendCmd.PersistentFlags().StringVar(&sendFlag.SenderFile, "send-file", "tmp/sender.list", "sender email address list")
	sendCmd.PersistentFlags().StringVar(&sendFlag.ConfigFile, "config", "tmp/config.yaml", "email content")
	sendCmd.PersistentFlags().StringVar(&sendFlag.Host, "host", "smtp.163.com:25", "email smtp host name")
	sendCmd.PersistentFlags().BoolVar(&sendFlag.Ssl, "ssl", false, "enable ssl")
	sendCmd.PersistentFlags().IntVar(&sendFlag.BatchSize, "batch-size", 10, "batch size")
	sendCmd.PersistentFlags().IntVar(&sendFlag.Delay, "delay", 1, "delay second")
	sendCmd.PersistentFlags().StringSliceVar(&sendFlag.Email, "email", []string{}, "email")
	sendCmd.PersistentFlags().StringSliceVar(&sendFlag.CC, "cc", []string{}, "extend email")
	sendCmd.PersistentFlags().IntVar(&sendFlag.Count, "count", -1, "how many batch to send")
	sendCmd.PersistentFlags().StringVar(&sendFlag.ProxyUrl, "proxy-url", "", "proxy address api address")
	rootCmd.AddCommand(sendCmd)
}
