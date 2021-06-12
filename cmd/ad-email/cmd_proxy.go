package main

import (
	"fmt"
	"github.com/juxuny/data-utils/log"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"strings"
)

var proxyFlag = struct {
	User string
	Pass string
	Url  string
}{}

var proxyCmd = &cobra.Command{
	Use: "set-proxy",
	Run: func(cmd *cobra.Command, args []string) {
		if proxyFlag.User == "" {
			log.Fatal("--user cannot empty")
		}
		if proxyFlag.Pass == "" {
			log.Fatal("--pass cannot empty")
		}
		resp, err := http.Get(proxyFlag.Url)
		if err != nil {
			log.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			log.Error(resp.Status)
		}
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Error(err)
		}
		address := strings.Trim(string(data), "\r\n\t ")
		httpProxy := fmt.Sprintf("http://%s:%s@%s", proxyFlag.User, proxyFlag.Pass, address)
		socksProxy := fmt.Sprintf("socks5://%s:%s@%s", proxyFlag.User, proxyFlag.Pass, address)
		fmt.Printf("http_proxy=%s\n", httpProxy)
		fmt.Printf("http_proxy=%s\n", httpProxy)
		fmt.Printf("all_proxy=%s\n", socksProxy)
	},
}

func init() {
	initGlobalFlag(proxyCmd)
	proxyCmd.PersistentFlags().StringVar(&proxyFlag.User, "user", "", "proxy user")
	proxyCmd.PersistentFlags().StringVar(&proxyFlag.Pass, "pass", "", "proxy password")
	proxyCmd.PersistentFlags().StringVar(&proxyFlag.Url, "url", "http://ww2502027.v4.dailiyun.com/query.txt?key=NP57E7DAA6&word=&count=1&rand=false&ltime=0&norepeat=false&detail=false", "proxy address api")
	rootCmd.AddCommand(proxyCmd)
}
