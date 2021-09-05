package main

import (
	"github.com/juxuny/data-utils/dict"
	"github.com/juxuny/data-utils/lib"
	"github.com/juxuny/data-utils/log"
	"github.com/spf13/cobra"
	"os"
)

var splitFlag = struct {
	DataFileCet4   string
	DataFileCet6   string
	InSrt          string
	OutSrt         string
	HighlightColor string
	FontFace       string
	FontSize       int
	GlobalFontSize int
	ExpandSeconds  int
	InputFile      string
	OutExt         string
	AutoRun        bool
}{}

func checkArgument() {
	if splitFlag.InputFile == "" {
		log.Fatal("missing --input -i argument")
	}
	if stat, err := os.Stat(splitFlag.InputFile); os.IsNotExist(err) {
		log.Fatal("input file not found:", splitFlag.InputFile)
	} else if stat.IsDir() {
		log.Fatal("input file is a directory")
	}
}

var splitCmd = &cobra.Command{
	Use: "split",
	Run: func(cmd *cobra.Command, args []string) {
		checkArgument()
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

		// auto run split script
		if splitFlag.AutoRun {
			if err := runCommand("bash", getSplitScriptFileName()); err != nil {
				log.Fatal(err)
			}
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
	splitCmd.PersistentFlags().IntVar(&splitFlag.FontSize, "size", 14, "font size. 48")
	splitCmd.PersistentFlags().IntVar(&splitFlag.ExpandSeconds, "expand", 10, "expand seconds")
	splitCmd.PersistentFlags().StringVarP(&splitFlag.InputFile, "input", "i", "", "input video file")
	splitCmd.PersistentFlags().StringVar(&splitFlag.OutExt, "ext", "mp4", "output video type")
	splitCmd.PersistentFlags().IntVar(&splitFlag.GlobalFontSize, "global-size", 12, "global font size")
	splitCmd.PersistentFlags().BoolVar(&splitFlag.AutoRun, "auto-run", true, "auto run the split script")

	rootCmd.AddCommand(splitCmd)
}
