package main

import (
	data_utils "github.com/juxuny/data-utils"
	"github.com/juxuny/data-utils/log"
	"github.com/juxuny/data-utils/model"
	"github.com/juxuny/data-utils/proxy"
	"github.com/juxuny/data-utils/proxy/dt"
	"github.com/spf13/cobra"
	"net/url"
	"time"
)

var checkFlag = struct {
	All                 bool
	ProxyAddress        []string
	Concurrent          int // the number of goroutine
	Times               int // test times
	IgnoreConnectFailed bool
}{}

var checkCmd = &cobra.Command{
	Use: "check",
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
		var iterator proxy.Iterator
		if len(checkFlag.ProxyAddress) > 0 {
			iterator = proxy.NewMemoryProxyAddressIterator()
			iterator.Init(checkFlag.ProxyAddress...)
		} else if checkFlag.All {
			iterator = proxy.NewDatabaseProxyAddressIterator(db, proxy.IteratorOption{IgnoreConnectFailed: checkFlag.IgnoreConnectFailed})
		} else {
			log.Fatal("--all and --address are empty")
			return
		}

		ch := make(chan string, 1000)
		resultChan := make(chan *proxy.TestResult, 100)
		for i := 0; i < checkFlag.Concurrent; i++ {
			go func() {
				data_utils.RecoverRun(func() {
					for addr := range ch {
						req, err := url.Parse(addr)
						if err != nil {
							log.Error(err)
							continue
						}
						ret, err := proxy.Test(addr, checkFlag.Times)
						if err != nil {
							log.Error(err)
							continue
						}
						if ret.FailedNum > 0 {
							log.Infof("check %s, failed ratio: %d/%d", addr, ret.FailedNum, ret.SuccessNum+ret.FailedNum)
							for err, count := range ret.ErrSummary {
								log.Infof("%v, count(%d)", err, count)
							}
							db.Table(model.Proxy{}.TableName()).Where("ip = ? AND port = ?", req.Hostname(), req.Port()).Updates(map[string]interface{}{
								"latency":    -1,
								"updated_at": time.Now(),
							})
						} else {
							log.Info("address: ", addr, " latency: ", ret.TotalLatency)
							// update latency
							db.Table(model.Proxy{}.TableName()).Where("ip = ? AND port = ?", req.Hostname(), req.Port()).Updates(map[string]interface{}{
								"latency":    ret.TotalLatency.Seconds(),
								"updated_at": time.Now(),
							})
						}
						resultChan <- ret
					}
				})
			}()
		}

		total := iterator.Len()
		go func() {
			for v, err := iterator.Next(); err != dt.ErrEOF; v, err = iterator.Next() {
				ch <- v
			}
		}()
		var progress = 0
		for i := 0; i < total; i++ {
			<-resultChan
			progress++
			log.Infof("progress (%d/%d)", progress, total)
		}
	},
}

func init() {
	initGlobalFlag(checkCmd)
	checkCmd.PersistentFlags().StringSliceVar(&checkFlag.ProxyAddress, "address", []string{}, "proxy address, e.g socks5://127.0.0.1:1080")
	checkCmd.PersistentFlags().BoolVarP(&checkFlag.All, "all", "a", true, "all proxy server in database")
	checkCmd.PersistentFlags().IntVar(&checkFlag.Concurrent, "concurrent", 5, "the number of goroutine")
	checkCmd.PersistentFlags().IntVar(&checkFlag.Times, "times", 5, "test times")
	checkCmd.PersistentFlags().BoolVar(&checkFlag.IgnoreConnectFailed, "ignore-connect-failed", false, "ignore connect failed")
	rootCmd.AddCommand(checkCmd)
}
