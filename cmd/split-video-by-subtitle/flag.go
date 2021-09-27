package main

import "github.com/spf13/cobra"

type globalFlag struct {
	Verbose  bool
	CacheDir string
}

func initGlobalFlag(cmd *cobra.Command, f *globalFlag) {
	cmd.PersistentFlags().BoolVarP(&f.Verbose, "verbose", "v", false, "display debug output")
	cmd.PersistentFlags().StringVar(&f.CacheDir, "cache", "tmp/cache", "cache data directory")
}
