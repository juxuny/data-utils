package main

import (
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

// 转换字幕
func convertSrt(inFile, outFile string) error {
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

	return nil
}
