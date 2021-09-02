package main

import (
	"github.com/juxuny/data-utils/dict"
	"github.com/juxuny/data-utils/log"
	"github.com/spf13/cobra"
)

var splitFlag = struct {
	DataFileCet4 string
	DataFileCet6 string
	InSrt        string
	OutSrt       string
}{}

var splitCmd = &cobra.Command{
	Use: "split",
	Run: func(cmd *cobra.Command, args []string) {
		dictCET4, err := dict.LoadDict(splitFlag.DataFileCet4)
		if err != nil {
			log.Fatal(err)
		}

		dictCET6, err := dict.LoadDict(splitFlag.DataFileCet6)
		if err != nil {
			log.Fatal(err)
		}
		//log.Info(dictCET6)
		//log.Info(dictCET6)
		log.Info("load CET4 words: ", len(dictCET4.Data))
		log.Info("load CET6 words: ", len(dictCET6.Data))

		// convert srt subtitle
		if err := convertSrt(splitFlag.InSrt, splitFlag.OutSrt); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	initGlobalFlag(splitCmd)
	splitCmd.PersistentFlags().StringVar(&splitFlag.DataFileCet4, "cet4", "tmp/dict/CET4_edited.txt", "CET4 words data")
	splitCmd.PersistentFlags().StringVar(&splitFlag.DataFileCet6, "cet6", "tmp/dict/CET6_edited.txt", "CET6 words data")
	splitCmd.PersistentFlags().StringVar(&splitFlag.InSrt, "in-srt", "tmp/eng.srt", "input subtitle file .srt")
	splitCmd.PersistentFlags().StringVar(&splitFlag.OutSrt, "out-srt", "tmp/eng.converted.srt", "output subtitle file .srt")

	rootCmd.AddCommand(splitCmd)
}
