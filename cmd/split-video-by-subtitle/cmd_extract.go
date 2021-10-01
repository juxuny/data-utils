package main

import "github.com/spf13/cobra"

type extractCmd struct {
	Flag struct {
		globalFlag
		InDir  string // working dir
		OutDir string // extract subtitle and save here
	}
}

func (t *extractCmd) initFlag(cmd *cobra.Command) {
	initGlobalFlag(cmd, &t.Flag.globalFlag)
	cmd.PersistentFlags().StringVar(&t.Flag.InDir, "in", "tmp", "working dir")
	cmd.PersistentFlags().StringVar(&t.Flag.OutDir, "out", "out", "extract subtitle and save as .srt")
}

func (t *extractCmd) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use: "extract",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
	t.initFlag(cmd)
	return cmd
}
