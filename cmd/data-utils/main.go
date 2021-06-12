package main

import (
	"github.com/juxuny/log-server/log"
	"github.com/spf13/cobra"
	"os"
)

var logger = log.NewLogger("[su]")

var (
	rootCmd = &cobra.Command{
		Use:   "data-utils",
		Short: "data-utils",
	}
)

var globalFlag = struct {
	Verbose bool
}{}

func initGlobalFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVarP(&globalFlag.Verbose, "verbose", "v", false, "display debug output")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		logger.Error(err)
		os.Exit(-1)
	}
}
