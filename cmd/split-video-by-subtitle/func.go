package main

import (
	"fmt"
	"github.com/juxuny/data-utils/dict"
	"github.com/juxuny/data-utils/log"
	"github.com/juxuny/data-utils/srt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
)

var globalFlag = struct {
	Verbose bool
}{}

func initGlobalFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVarP(&globalFlag.Verbose, "verbose", "v", false, "display debug output")

}

type CetFilter func(content string) (words []dict.Word)

// 转换字幕
func convertSrt(inFile, outFile string, filter CetFilter) error {
	inData, err := ioutil.ReadFile(inFile)
	if err != nil {
		log.Debug(err)
		return errors.Wrap(err, "convert srt failed")
	}
	blocks, err := srt.Parse(inData)
	if err != nil {
		log.Debug(err)
		return errors.Wrap(err, "convert srt failed")
	}

	for _, b := range blocks {
		b.Subtitle.ForEach(func(n *srt.Node) bool {
			if n.Type != srt.NodeTypeText {
				return false
			}
			words := filter(string(n.Content))
			if len(words) == 0 {
				return false
			}
			for _, w := range words {
				log.Info("found word: ", w.Name)
				//n.Content = []byte(strings.ReplaceAll(string(n.Content), w.Name, fmt.Sprintf("<font size=\"48\"><b>%s</b></font>", w.Name)))
				n.Content = dict.ReplaceWords(
					n.Content,
					[]byte(w.Name),
					[]byte(fmt.Sprintf(
						"<font face=\"%s\" size=\"%d\" color=\"%s\" size=\"48\"><b>%s</b></font>",
						splitFlag.FontFace,
						splitFlag.FontSize,
						splitFlag.HighlightColor,
						w.Name,
					)),
				)
			}
			return false
		})
	}

	// write srt data
	out, err := os.Create(outFile)
	if err != nil {
		log.Debug(err)
		return errors.Wrap(err, "create out srt failed")
	}
	defer out.Close()
	for _, b := range blocks {
		if _, err := out.WriteString(b.String() + "\n\n"); err != nil {
			break
		}
	}
	log.Info("saved srt file: ", outFile)
	return nil
}
