package main

import (
	"github.com/juxuny/data-utils/log"
	"github.com/spf13/cobra"
	"io/ioutil"
	"path"
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
			list, err := ioutil.ReadDir(t.Flag.In)
			if err != nil {
				log.Fatal(err)
			}
			for _, item := range list {
				if !item.IsDir() {
					continue
				}
				subtitleList, err := ioutil.ReadDir(path.Join(t.Flag.In, item.Name()))
				if err != nil {
					log.Fatal(err)
				}
				for _, sub := range subtitleList {
					ext := path.Ext(sub.Name())
					if ext != ".srt" {
						continue
					}
					src := path.Join(t.Flag.In, item.Name(), sub.Name())
					dst := path.Join(t.Flag.In, item.Name(), item.Name()+".subtitle", sub.Name())
					log.Debug(src, " ", dst)
					if err := runCommand("mv", src, dst); err != nil {
						log.Fatal(err)
					}
				}
			}
		},
	}
	t.initFlag(cmd)
	return cmd
}

func init() {
	rootCmd.AddCommand((&unpackCmd{}).Build())
}
