package main

import (
	"github.com/juxuny/data-utils/log"
	"github.com/juxuny/data-utils/model"
	"github.com/juxuny/data-utils/proxy"
	"github.com/spf13/cobra"
)

var fetchFlag = struct {
	Providers []string
	All       bool
}{}

var fetchCmd = &cobra.Command{
	Use: "fetch",
	Run: func(cmd *cobra.Command, args []string) {
		if !fetchFlag.All && len(fetchFlag.Providers) == 0 {
			log.Fatal("you should use --all or --provider")
		}
		dbConfig := getDbConfigFromCommandLineArgs()
		db, err := model.Open(dbConfig)
		if err != nil {
			log.Fatal(err)
			return
		}
		defer func() {
			_ = db.Close()
		}()
		if err := proxy.InitFetcher(proxy.NewSaveServerListHandler(db)); err != nil {
			log.Fatal(err)
			return
		}
		_, err = proxy.Fetch(false)
		if err != nil {
			log.Fatal(err)
			return
		}
	},
}

func init() {
	initGlobalFlag(fetchCmd)
	fetchCmd.PersistentFlags().StringSliceVar(&fetchFlag.Providers, "provider", []string{}, "provider name (66ip)")
	fetchCmd.PersistentFlags().BoolVarP(&fetchFlag.All, "all", "a", false, "use all provider name")
	rootCmd.AddCommand(fetchCmd)
}
