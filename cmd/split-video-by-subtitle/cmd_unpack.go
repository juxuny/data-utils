package main

import (
	"github.com/juxuny/data-utils/log"
	"github.com/spf13/cobra"
	"io/fs"
	"path"
	"path/filepath"
)

type unpackCmd struct {
	Flag struct {
		globalFlag
		In string
	}
}

func (t *unpackCmd) initFlag(cmd *cobra.Command) {
	initGlobalFlag(cmd, &t.Flag.globalFlag)
	cmd.PersistentFlags().StringVar(&t.Flag.In, "in", ".", "data directory")
}

func (t *unpackCmd) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use: "unpack",
		Run: func(cmd *cobra.Command, args []string) {
			filepath.Walk(t.Flag.In, func(fullFilePath string, info fs.FileInfo, err error) error {
				ext := path.Ext(fullFilePath)
				if ext != ".zip" {
					return nil
				}
				log.Debug(fullFilePath)
				if err := runCommand("unzip", fullFilePath); err != nil {
					log.Error(err)
				}
				return nil
			})
		},
	}
	t.initFlag(cmd)
	return cmd
}

func init() {
	rootCmd.AddCommand((&unpackCmd{}).Build())
}
