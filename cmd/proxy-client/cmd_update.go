package main

import (
	"context"
	"fmt"
	"github.com/juxuny/data-utils/log"
	"github.com/juxuny/data-utils/model"
	px "github.com/juxuny/data-utils/proxy"
	"github.com/juxuny/supervisor/proxy"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"time"
)

var updateFlag = struct {
	Host string
}{}

func getController(host string) (client proxy.ProxyClient, err error) {
	conn, err := grpc.Dial(host, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, errors.Wrap(err, "connect failed")
	}
	client = proxy.NewProxyClient(conn)
	return client, nil
}

var updateCmd = &cobra.Command{
	Use: "update",
	Run: func(cmd *cobra.Command, args []string) {
		proxyController, err := getController(updateFlag.Host)
		if err != nil {
			log.Fatal(err)
			return
		}

		dbConfig := getDbConfigFromCommandLineArgs()
		var db *model.DB
		db, err = model.Open(dbConfig)
		if err != nil {
			log.Fatal(err)
			return
		}
		defer func() {
			_ = db.Close()
		}()
		proxyServerInfo, err := px.RandServer(db)
		if err != nil {
			log.Fatal(err)
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		_, err = proxyController.Update(ctx, &proxy.UpdateReq{
			Status: &proxy.Status{
				Remote: fmt.Sprintf("%s:%d", proxyServerInfo.Ip, proxyServerInfo.Port),
			},
		})
		if err != nil {
			log.Fatal(err)
		}
		log.Info("set proxy ", fmt.Sprintf("%s:%d", proxyServerInfo.Ip, proxyServerInfo.Port))
	},
}

func init() {
	initGlobalFlag(updateCmd)
	updateCmd.PersistentFlags().StringVar(&updateFlag.Host, "host", "127.0.0.1:9999", "proxy server control plane")
	rootCmd.AddCommand(updateCmd)
}
