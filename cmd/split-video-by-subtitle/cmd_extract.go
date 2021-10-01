package main

import (
	"github.com/juxuny/data-utils/ffmpeg"
	"github.com/juxuny/data-utils/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"io/fs"
	"path"
	"path/filepath"
	"strings"
)

type extractCmd struct {
	Flag struct {
		globalFlag
		InDir  string // working dir
		OutDir string // extract subtitle and save here
	}
}

func (t *extractCmd) initFlag(cmd *cobra.Command) {
	initGlobalFlag(cmd, &t.Flag.globalFlag)
	cmd.PersistentFlags().StringVar(&t.Flag.InDir, "in", "", "working dir")
	cmd.PersistentFlags().StringVar(&t.Flag.OutDir, "out", "", "extract subtitle and save as .srt")
}

func (t *extractCmd) extractSubtitle(videoDir string, inputFileName string, outDir string, streamIndex int) error {
	ext := path.Ext(inputFileName)
	baseName := strings.TrimRight(inputFileName, ext)
	outSubtitle := path.Join(outDir, baseName+".srt")
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
			if t.Flag.OutDir == "" {
				log.Fatal("missing argument: --out")
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
					for _, s := range videoInfo.Streams {
						if s.CodecType == ffmpeg.CodecTypeSubtitle && s.Tags[ffmpeg.TagKey.Language] == "eng" {
							log.Debug(s.Tags[ffmpeg.TagKey.Title], s.Tags[ffmpeg.TagKey.Language])
							t.extractSubtitle(baseDir, d.Name(), t.Flag.OutDir, s.Index)
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
