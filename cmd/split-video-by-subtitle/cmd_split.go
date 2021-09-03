package main

import (
	"github.com/juxuny/data-utils/dict"
	"github.com/juxuny/data-utils/lib"
	"github.com/juxuny/data-utils/log"
	"github.com/spf13/cobra"
)

var splitFlag = struct {
	DataFileCet4   string
	DataFileCet6   string
	InSrt          string
	OutSrt         string
	HighlightColor string
	FontFace       string
	FontSize       int
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
		if err := convertSrt(splitFlag.InSrt, splitFlag.OutSrt, func(content string) (words []dict.Word) {
			splitWords := lib.SplitByCharset([]byte(content), " ?.!<>")
			for _, w := range splitWords {
				if v, b := dictCET4.Data[string(w)]; b {
					words = append(words, v)
				}
			}
			return
		}); err != nil {
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
	splitCmd.PersistentFlags().StringVar(&splitFlag.HighlightColor, "color", "#f7db9f", "highlight color e.g #fff0cf")
	splitCmd.PersistentFlags().StringVar(&splitFlag.FontFace, "face", "Cronos Pro Light", "font face. e.g Cronos Pro Light")
	splitCmd.PersistentFlags().IntVar(&splitFlag.FontSize, "size", 48, "font size. 48")

	rootCmd.AddCommand(splitCmd)
}
