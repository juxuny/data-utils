package main

import (
	"encoding/json"
	"fmt"
	"github.com/juxuny/data-utils/ffmpeg"
	"github.com/juxuny/data-utils/lib"
	"github.com/juxuny/data-utils/log"
	"github.com/juxuny/data-utils/model"
	"github.com/juxuny/data-utils/srt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"
	"unicode"
)

type searchCmd struct {
	Flag struct {
		globalFlag
		SearchDir          string
		Key                string
		ResultDir          string // save the search result temporary
		Mode               string // search mode
		Limit              int    // limit the number of result
		Ext                string
		OutExt             string
		NeedContainChinese bool // check if contain Chinese words
		FontSize           int
	}

	outDir string
	db     *model.DB
}

type SearchMode string

const (
	SearchModeBool   = SearchMode("bool")
	SearchModeNormal = SearchMode("normal")
	SearchModeLike   = SearchMode("like")
)

func (t SearchMode) IsValid() bool {
	switch t {
	case SearchModeNormal, SearchModeBool, SearchModeLike:
		return true
	}
	return false
}

func (t *searchCmd) initFlag(cmd *cobra.Command) {
	initGlobalFlag(cmd, &t.Flag.globalFlag)
	cmd.PersistentFlags().StringVar(&t.Flag.SearchDir, "search-dir", ".", "search directory, will search video in this directory")
	cmd.PersistentFlags().StringVar(&t.Flag.Key, "key", "", "search key")
	cmd.PersistentFlags().StringVar(&t.Flag.ResultDir, "out", "tmp", "temporary directory, save the last search result here")
	cmd.PersistentFlags().StringVar(&t.Flag.Mode, "mode", "normal", "search mode, bool or normal")
	cmd.PersistentFlags().IntVar(&t.Flag.Limit, "limit", 10, "limit number of result")
	cmd.PersistentFlags().StringVar(&t.Flag.Ext, "ext", "mkv", "video type")
	cmd.PersistentFlags().StringVar(&t.Flag.OutExt, "out-ext", "mp4", "output video extension")
	cmd.PersistentFlags().BoolVar(&t.Flag.NeedContainChinese, "contain-chinese", false, "whether need contain Chinese words")
	cmd.PersistentFlags().IntVar(&t.Flag.FontSize, "font-size", 18, "font size of subtitle")
}

func (t *searchCmd) loadFileList(dir string) (list []string, err error) {
	var fileList []fs.FileInfo
	fileList, err = ioutil.ReadDir(dir)
	if err != nil {
		log.Error(err)
		return nil, errors.Wrap(err, "load dir failed")
	}
	for _, f := range fileList {
		if !f.IsDir() {
			continue
		}
		list = append(list, path.Join(dir, f.Name()))
	}
	return
}

func (t *searchCmd) getVideoFileList(dir, ext string) (list []string, err error) {
	l, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Error(err)
		return nil, errors.Wrap(err, "get video file list")
	}
	for _, f := range l {
		if !f.IsDir() {
			continue
		}
		fileExt := path.Ext(f.Name())
		if strings.Trim(fileExt, ".") != ext {
			continue
		}
		list = append(list, path.Join(dir, f.Name()))
	}
	return
}

func (t *searchCmd) findVideoFile(fileList []string, videoName string, srtName string) (file string, err error) {
	for _, item := range fileList {
		tmpVideoName := path.Base(item)
		tmpVideoName = lib.String.TrimSubStringRight(tmpVideoName, t.Flag.Ext)
		tmpVideoName = strings.Trim(tmpVideoName, ".")
		if tmpVideoName == videoName {
			videoFileList, err := t.getVideoFileList(item, t.Flag.Ext)
			if err != nil {
				log.Error(err)
				continue
			}
			for _, vf := range videoFileList {
				log.Debug(vf)
				videoDir := path.Dir(vf)
				name := path.Base(vf)
				ext := path.Ext(vf)
				name = lib.String.TrimSubStringRight(name, ext)
				srtDir := path.Join(videoDir, name+SubtitleDirSuffix)
				log.Debug("===>", srtDir)
				if stat, err := os.Stat(path.Join(srtDir, srtName)); err == nil && !stat.IsDir() {
					return vf, nil
				}
			}
		}
	}
	return "", lib.ErrNotFound
}

func (t *searchCmd) getPadConfig() string {
	return fmt.Sprintf("iw:(iw*(16/9)):0:((iw*(16/9))-ih)/2:black")
}

func (t *searchCmd) hasChinese(in string) bool {
	for _, x := range in {
		if unicode.Is(unicode.Han, x) {
			return true
		}
	}
	return false
}

func (t *searchCmd) getMarginV(rawVideoFile string) int {
	videoInfo, err := ffmpeg.GetVideoInfo(rawVideoFile)
	if err != nil {
		log.Error(err)
		return 0
	}
	width, height, err := videoInfo.GetVideoSize()
	if err != nil {
		log.Error(err)
		return 0
	}
	outputHeight := int64(float64(width) * (float64(16) / float64(9)))

	delta := (outputHeight - height) >> 1
	return int(delta) >> 2
}

func (t *searchCmd) generateSplitScript(result SearchResult) {
	script := "#!/bin/bash\nset -e\n"
	// generate split script
	for _, item := range result.List {

		name := getFileNameWithoutExt(item.Subtitle.FileName)
		blockExpand, err := item.Block.ExpandSubtitleDuration(1)
		if err != nil {
			log.Error(err)
			return
		}
		start, err := lib.Time.Parse(lib.TimeInMillionLayout, blockExpand.StartTime)
		if err != nil {
			log.Error(err)
			continue
		}
		end, err := lib.Time.Parse(lib.TimeInMillionLayout, blockExpand.EndTime)
		if err != nil {
			log.Error(err)
			continue
		}
		end = end.Add(500 * time.Millisecond)
		duration := srt.ZeroTime.Add(end.Sub(start)).Format(timeLayout)
		out := path.Join(t.outDir, fmt.Sprintf("%s.split.%d.%s", name, item.Block.BlockId, t.Flag.OutExt))
		script += fmt.Sprintf(
			"if [ ! -f '%s' ]; then ffmpeg -y -i '%s' -ss %s -t %s '%s'; fi\n",
			out,
			item.VideoFile,
			start.Format(lib.TimeLayout),
			duration,
			//t.getPadConfig(),
			out,
		)
	}
	// add subtitle
	for _, item := range result.List {
		name := getFileNameWithoutExt(item.Subtitle.FileName)
		splitFile := path.Join(t.outDir, fmt.Sprintf("%s.split.%d.%s", name, item.Block.BlockId, t.Flag.OutExt))
		out := path.Join(t.outDir, fmt.Sprintf("%s.subtitle.%d.%s", name, item.Block.BlockId, t.Flag.OutExt))
		srtFile := path.Join(t.outDir, fmt.Sprintf("%s.%d.srt", name, item.Block.BlockId))
		script += fmt.Sprintf(
			"ffmpeg -y -i '%s' -vf subtitles='%s':force_style=\\'FontSize=%d\\' '%s'\n",
			splitFile,
			srtFile,
			t.Flag.FontSize,
			out,
		)
	}
	// padding
	outList := make([]string, 0)
	for _, item := range result.List {
		name := getFileNameWithoutExt(item.Subtitle.FileName)
		inFile := path.Join(t.outDir, fmt.Sprintf("%s.subtitle.%d.%s", name, item.Block.BlockId, t.Flag.OutExt))
		out := path.Join(t.outDir, fmt.Sprintf("%s.pad.%d.%s", name, item.Block.BlockId, t.Flag.OutExt))
		if !t.Flag.NeedContainChinese || t.hasChinese(out) {
			outList = append(outList, out)
		}
		script += fmt.Sprintf(
			"ffmpeg -y -i '%s' -vf 'pad=%s' '%s'\n",
			inFile,
			t.getPadConfig(),
			out,
		)
	}
	// generate concat.list
	concatFile := path.Join(t.outDir, "concat.list")
	concatData := ""
	for _, item := range outList {
		concatData += fmt.Sprintf("file '%s'\n", item)
	}
	if err := ioutil.WriteFile(concatFile, []byte(concatData), 0755); err != nil {
		log.Fatal(err)
	}
	// concat
	mergedVideo := path.Join(t.outDir, "merged.mp4")
	//script += fmt.Sprintf(
	//	"ffmpeg -y -i 'concat:%s' -vf select=concatdec_select -af aselect=concatdec_select,aresample=async=1 '%s'\n",
	//	strings.Join(outList, "|"),
	//	mergedVideo,
	//)
	script += fmt.Sprintf(
		"# ffmpeg -safe 0 -y -f concat -i '%s' '%s'\n",
		concatFile,
		mergedVideo,
	)
	scriptFile := path.Join(t.outDir, "split.sh")
	if err := ioutil.WriteFile(scriptFile, []byte(script), 0755); err != nil {
		log.Fatal(err)
	}
	log.Info("generate split script: ", scriptFile)
}

func (t *searchCmd) searchAndSaveResult() {
	var blocks model.EngSubtitleBlockList
	mode := "BOOLEAN"
	var query string
	var values []interface{}
	searchMode := SearchMode(t.Flag.Mode)
	if searchMode == SearchModeNormal {
		mode = "NATURAL LANGUAGE"
		query = fmt.Sprintf("SELECT * FROM eng_subtitle_block WHERE MATCH(content) AGAINST(? IN %s MODE) LIMIT ?", mode)
		values = append(values, t.Flag.Key, t.Flag.Limit)
	} else if searchMode == SearchModeBool {
		query = fmt.Sprintf("SELECT * FROM eng_subtitle_block WHERE MATCH(content) AGAINST(? IN %s MODE) LIMIT ?", mode)
		values = append(values, t.Flag.Key, t.Flag.Limit)
	} else if searchMode == SearchModeLike {
		query = fmt.Sprintf("SELECT * FROM eng_subtitle_block WHERE content LIKE ? LIMIT ?")
		values = append(values, fmt.Sprintf("%%%s%%", t.Flag.Key), t.Flag.Limit)
	} else {
		log.Fatal("invalid --mode=" + t.Flag.Mode)
	}
	rows, err := t.db.Raw(query, values...).Rows()
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		item := model.EngSubtitleBlock{}
		if err := t.db.ScanRows(rows, &item); err != nil {
			log.Error(err)
			break
		}
		blocks = append(blocks, item)
	}
	_ = rows.Close()
	log.Debug(len(blocks))

	movieMap := make(map[int64]model.EngMovie)
	subtitleMap := make(map[int64]model.EngSubtitle)
	subtitleIds := make([]int64, 0)
	for _, item := range blocks {
		subtitleIds = append(subtitleIds, item.SubtitleId)
	}
	var subtitles model.EngSubtitleList
	var movies model.EngMovieList
	var movieIds = make([]int64, 0)
	if len(subtitleIds) > 0 {
		if err := t.db.Where("id IN (?)", subtitleIds).Find(&subtitles).Error; err != nil {
			log.Fatal(err)
		}
		for _, item := range subtitles {
			subtitleMap[item.Id] = item
			movieIds = append(movieIds, item.MovieId)
		}
	}

	if len(movieIds) > 0 {
		if err := t.db.Where("id IN (?)", movieIds).Find(&movies).Error; err != nil {
			log.Fatal(err)
		}
		for _, item := range movies {
			movieMap[item.Id] = item
		}
	}

	fileList, err := t.loadFileList(t.Flag.SearchDir)
	if err != nil {
		log.Fatal(err)
	}
	var result SearchResult
	for _, block := range blocks {
		item := ResultItem{
			Block: block,
		}
		if subtitle, b := subtitleMap[block.SubtitleId]; b {
			item.Subtitle = subtitle
		}
		if movie, b := movieMap[item.Subtitle.MovieId]; b {
			item.Movie = movie
		}
		//log.Debug(path.Join(item.Movie.Name, item.Subtitle.FileName))
		videoFile, err := t.findVideoFile(fileList, item.Movie.Name, item.Subtitle.FileName)
		if err != nil {
			log.Info(err, " ", item.Movie.Name)
			continue
		}
		item.VideoFile = videoFile
		result.List = append(result.List, item)
	}

	// save result
	t.outDir = path.Join(t.Flag.ResultDir, time.Now().Format("20060102_150405"))
	if err := os.MkdirAll(t.outDir, 0755); err != nil {
		log.Fatal(err)
	}
	out := path.Join(t.outDir, "result.json")
	data, _ := json.Marshal(result)
	if err := ioutil.WriteFile(out, data, 0644); err != nil {
		log.Fatal("save result failed: " + err.Error())
	}
	// save split subtitle
	for _, item := range result.List {
		ext := path.Ext(item.Subtitle.FileName)
		baseName := strings.TrimRight(item.Subtitle.FileName, ext)
		outSrt := path.Join(t.outDir, baseName+fmt.Sprintf(".%d.srt", item.Block.BlockId))

		beginningBlock, err := item.Block.MoveToBeginning()
		if err != nil {
			log.Fatal(err)
		}
		//beginningBlockExpend, err := beginningBlock.ExpandSubtitleDuration(1)
		//if err != nil {
		//	log.Fatal(err)
		//}
		srtFormatData := fmt.Sprintf(
			"%d\n%s --> %s\n%s",
			beginningBlock.BlockId,
			beginningBlock.StartTime,
			beginningBlock.EndTime,
			beginningBlock.Content,
		)
		if err := ioutil.WriteFile(outSrt, []byte(srtFormatData), 0644); err != nil {
			log.Fatal(err)
		}
	}
	// generate script
	t.generateSplitScript(result)
	log.Info("save search result in: " + out)
	log.Info("result: ", len(blocks))
	for _, b := range blocks {
		log.Info(b.Content)
	}
}

func (t *searchCmd) Build() *cobra.Command {
	cmd := &cobra.Command{
		Use: "search",
		Run: func(cmd *cobra.Command, args []string) {
			// check arguments
			if t.Flag.Key == "" {
				log.Fatal("missing argument: --key")
			}
			searchMode := SearchMode(t.Flag.Mode)
			if !searchMode.IsValid() {
				log.Fatal("invalid --mode, only allow: bool or normal")
			}

			// init database connection
			var err error
			t.db, err = model.Open()
			if err != nil {
				log.Fatal(err)
			}
			defer func(db *model.DB) {
				err := db.Close()
				if err != nil {
					log.Fatal(err)
				}
			}(t.db)

			//search
			t.searchAndSaveResult()
		},
	}
	t.initFlag(cmd)
	return cmd
}

func init() {
	rootCmd.AddCommand((&searchCmd{}).Build())
}
