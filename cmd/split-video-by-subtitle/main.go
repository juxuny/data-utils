package main

import (
	"github.com/juxuny/log-server/log"
	"github.com/spf13/cobra"
	"os"
)

var logger = log.NewLogger("[split]")

var (
	rootCmd = &cobra.Command{
		Use:   "split",
		Short: "split",
	}
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		logger.Error(err)
		os.Exit(-1)
	}
}
