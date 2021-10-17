package main

import (
	"fmt"
	"github.com/juxuny/data-utils/ffmpeg"
	"github.com/juxuny/data-utils/lib"
	"github.com/juxuny/data-utils/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type extractCmd struct {
	Flag struct {
		globalFlag
		InDir string // working dir
		//OutDir string // extract subtitle and save here
	}
}

func (t *extractCmd) initFlag(cmd *cobra.Command) {
	initGlobalFlag(cmd, &t.Flag.globalFlag)
	cmd.PersistentFlags().StringVar(&t.Flag.InDir, "in", "", "working dir")
	//cmd.PersistentFlags().StringVar(&t.Flag.OutDir, "out", "", "extract subtitle and save as .srt")
}

func (t *extractCmd) extractSubtitle(videoDir string, inputFileName string, outDir string, streamIndex int, language string) error {
	outSubtitle := path.Join(outDir, fmt.Sprintf("%d.%s", streamIndex, language)+".srt")
	if _, err := os.Stat(outSubtitle); err == nil {
		log.Warn("ignore: ", outSubtitle)
		return nil
	}
	in := path.Join(videoDir, inputFileName)
	return ffmpeg.ExtractSubtitle(in, outSubtitle, streamIndex)
}

func (t *extractCmd) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use: "extract",
		Run: func(cmd *cobra.Command, args []string) {
			if t.Flag.InDir == "" {
				log.Fatal("missing argument: --in")
			}
			if err := filepath.WalkDir(t.Flag.InDir, func(fileFullPath string, d fs.DirEntry, err error) error {
				baseDir := path.Dir(fileFullPath)
				if d.IsDir() {
					return nil
				}
				ext := path.Ext(d.Name())
				if ext == ".mkv" || ext == ".mp4" {
					log.Info("extract: ", path.Join(baseDir, d.Name()))
					videoInfo, err := ffmpeg.GetVideoInfo(fileFullPath)
					if err != nil {
						log.Error(err)
						return errors.Wrap(err, "invalid video")
					}
					outDir := path.Join(baseDir, fmt.Sprintf("%s.subtitle", strings.TrimRight(d.Name(), ext)))
					if err := lib.TouchDir(outDir, 0757); err != nil {
						log.Error(err)
					}
					subtitleStreams := videoInfo.GetSubtitleStream()
					if len(subtitleStreams) > 0 {
						for index, s := range subtitleStreams {
							log.Debug(s.Tags[ffmpeg.TagKey.Title], s.Tags[ffmpeg.TagKey.Language])
							t.extractSubtitle(baseDir, d.Name(), outDir, index, s.Tags[ffmpeg.TagKey.Language])
						}
					}
				}
				return nil
			}); err != nil {
				log.Fatal(err)
			}
		},
	}
	t.initFlag(cmd)
	return cmd
}

func init() {
	rootCmd.AddCommand((&extractCmd{}).Build())
}
