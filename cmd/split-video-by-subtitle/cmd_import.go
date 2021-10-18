package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/juxuny/data-utils/lib"
	"github.com/juxuny/data-utils/log"
	"github.com/juxuny/data-utils/model"
	"github.com/juxuny/data-utils/srt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os"
	"path"
	"strings"
)

type importCmd struct {
	Flag struct {
		globalFlag
		In string
	}

	db *model.DB
}

func (t *importCmd) initFlag(cmd *cobra.Command) {
	initGlobalFlag(cmd, &t.Flag.globalFlag)
	cmd.PersistentFlags().StringVar(&t.Flag.In, "in", ".", "data directory")
}

func (t *importCmd) saveBlockList(fileName string, blocks []srt.Block) error {
	baseName := path.Base(fileName)
	ext := path.Ext(baseName)
	name := strings.TrimRight(baseName, ext)
	var count int
	if err := t.db.Model(&model.EngMovie{}).Where("name = ?", name).Count(&count).Error; err != nil {
		return errors.Wrap(err, "query eng_movie failed")
	}
	if count > 0 {
		return nil // the same srt data exists, ignore it
	}
	if err := t.db.Model(&model.EngSubtitle{}).Where("file_name = ?", baseName).Count(&count).Error; err != nil {
		return errors.Wrap(err, "query eng_subtitle failed")
	}
	if count > 0 {
		return nil // ignore
	}
	return t.db.Begin(func(db *gorm.DB) error {
		var movie = model.EngMovie{
			Id:         0,
			Name:       name,
			ParentId:   0,
			CreateTime: lib.Time.NowPointer(),
		}
		if err := db.Create(&movie).Error; err != nil {
			log.Error(err)
			return errors.Wrap(err, "save movie info failed")
		}
		var subtitle = model.EngSubtitle{
			Id:         0,
			MovieId:    movie.Id,
			Ext:        ext,
			FileName:   baseName,
			CreateTime: lib.Time.NowPointer(),
		}
		if err := db.Create(&subtitle).Error; err != nil {
			log.Error(err)
			return errors.Wrap(err, "save subtitle info failed")
		}
		for _, b := range blocks {
			subtitleBlock := model.EngSubtitleBlock{
				Id:             0,
				SubtitleId:     subtitle.Id,
				BlockId:        b.Id,
				StartTime:      b.StartTime.Format(srt.IntervalFormat),
				EndTime:        b.EndTime.Format(srt.IntervalFormat),
				DurationExtend: "",
				Content:        b.Content(),
				CreateTime:     lib.Time.NowPointer(),
			}
			if err := db.Create(&subtitleBlock).Error; err != nil {
				log.Error(err)
				return errors.Wrap(err, "save subtitle block failed")
			}
		}

		return nil
	})
}

func (t *importCmd) convertToSrt(dir string) error {
	log.Info("convert ass to srt...")
	assFileList, err := loadFileList(dir, "ass")
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range assFileList {
		ext := path.Ext(f)
		out := strings.TrimRight(f, ext) + ".srt"
		if _, err := os.Stat(out); err == nil {
			log.Info("ignore: ", out)
			continue
		}
		if err := runCommand("ffmpeg", "-i", f, out); err != nil {
			return errors.Wrap(err, "convert ass to srt failed, "+f)
		}
	}
	return nil
}

func (t *importCmd) importTvSeries(name string) error {
	ext := path.Ext(name)
	tvName := strings.TrimRight(name, ext)
	series, err := os.ReadDir(path.Join(t.Flag.In, name))
	if err != nil {
		return errors.Wrap(err, "load TV series failed")
	}
	for _, item := range series {
		if !item.IsDir() {
			continue
		}
		resourcePath := path.Join(t.Flag.In, name, item.Name())
		if err := t.importSubtitleResource(tvName, resourcePath, model.ResTypeTv); err != nil {
			log.Error(err)
			return errors.Wrap(err, "import subtitle resource failed")
		}
	}
	return nil
}

func (t *importCmd) importSubtitleResource(name string, dir string, resType model.ResType) error {
	log.Debug("importing: ", dir)
	subtitleDirList, err := os.ReadDir(dir)
	if err != nil {
		log.Error(err)
		return errors.Wrap(err, "load subtitle dir list failed")
	}
	if err := t.db.Begin(func(db *gorm.DB) error {
		var movie model.EngMovie
		// check if movie exists
		if err := db.Where("name = ?", name).First(&movie).Error; err != nil {
			if !model.IsErrNoDataInDb(err) {
				log.Error(err)
				return errors.Wrap(err, "query database failed")
			} else {
				// not found movie, create a new record
				movie = model.EngMovie{
					Name:       name,
					ParentId:   0,
					ResType:    resType,
					CreateTime: lib.Time.NowPointer(),
					Status:     model.StatusEnable,
				}
				if err := db.Create(&movie).Error; err != nil {
					log.Error(err)
					return errors.Wrap(err, "init movie failed")
				}
			}
		}
		for _, subtitleDir := range subtitleDirList {
			dirExt := path.Ext(subtitleDir.Name())
			if dirExt != ".subtitle" {
				continue
			}
			srtList, err := os.ReadDir(path.Join(dir, subtitleDir.Name()))
			if err != nil {
				log.Error(err)
				return errors.Wrap(err, "load subtitle dir failed")
			}

			for _, item := range srtList {
				ext := path.Ext(item.Name())
				if ext != ".srt" {
					continue
				}
				subtitleItem := model.EngSubtitle{
					MovieId:    movie.Id,
					SubName:    strings.TrimRight(subtitleDir.Name(), dirExt),
					Ext:        ext,
					FileName:   item.Name(),
					CreateTime: lib.Time.NowPointer(),
				}
				if err := db.Where("sub_name = ? AND movie_id = ?", subtitleDir.Name(), movie.Id).FirstOrCreate(&subtitleItem).Error; err != nil {
					log.Error(err)
					return errors.Wrap(err, "init subtitle")
				}
				file := path.Join(dir, subtitleDir.Name(), item.Name())
				log.Debug(file)
				blocks, err := srt.ParseFile(file)
				if err != nil {
					log.Error(err)
					continue
				}
				insertExecutor := model.NewInsertExecutor(true).Model(model.EngSubtitleBlock{})
				for _, b := range blocks {
					if insertExecutor.DataLen() > 100 {
						if _, err := insertExecutor.Exec(db); err != nil {
							log.Error(err)
							return errors.Wrap(err, "save blocks failed")
						}
						model.NewInsertExecutor(true).Model(model.EngSubtitleBlock{})
					}
					insertExecutor.Add(model.EngSubtitleBlock{
						SubtitleId:     subtitleItem.Id,
						BlockId:        b.Id,
						StartTime:      b.StartTime.Format(lib.TimeInMillionLayout),
						EndTime:        b.EndTime.Format(lib.TimeInMillionLayout),
						DurationExtend: "",
						Content:        b.Content(),
						CreateTime:     lib.Time.NowPointer(),
					})
				}
				if insertExecutor.DataLen() > 0 {
					if _, err := insertExecutor.Exec(db); err != nil {
						log.Error(err)
						return errors.Wrap(err, "save blocks failed")
					}
				}
			}
		}
		return nil
	}); err != nil {
		return errors.Wrap(err, "save subtitle failed")
	}

	return nil
}

func (t *importCmd) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use: "import",
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			t.db, err = model.Open()
			if err != nil {
				log.Fatal(err)
			}

			fileList, err := os.ReadDir(t.Flag.In)
			if err != nil {
				log.Fatal(err)
			}
			for i, moviePath := range fileList {
				log.Info(fmt.Sprintf("%.02f%%", float64(i+1)/float64(len(fileList))*100))
				stat, err := os.Stat(path.Join(t.Flag.In, moviePath.Name()))
				if err != nil {
					log.Error(err)
					continue
				}
				if !stat.IsDir() {
					log.Warn("ignore: ", moviePath.Name(), " is not a directory")
					continue
				}
				ext := path.Ext(moviePath.Name())
				if ext == ".d" {
					// It's a TV series
					if err := t.importTvSeries(moviePath.Name()); err != nil {
						log.Error(err)
					}
				} else {
					if err := t.importSubtitleResource(moviePath.Name(), path.Join(t.Flag.In, moviePath.Name()), model.ResTypeNormal); err != nil {
						log.Error(err)
					}
				}
				log.Debug("finished: ", moviePath.Name())
				//if err := t.convertToSrt(path.Join(moviePath.Name(), "subtitle")); err != nil {
				//	log.Fatal(err)
				//}
			}
		},
	}
	t.initFlag(cmd)
	return cmd
}

func init() {
	rootCmd.AddCommand((&importCmd{}).Build())
}
