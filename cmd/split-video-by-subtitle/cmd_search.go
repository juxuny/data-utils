package main

import (
	"encoding/json"
	"fmt"
	"github.com/juxuny/data-utils/ffmpeg"
	"github.com/juxuny/data-utils/lib"
	"github.com/juxuny/data-utils/log"
	"github.com/juxuny/data-utils/model"
	"github.com/juxuny/data-utils/srt"
	"github.com/spf13/cobra"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
	"unicode"
)

type searchCmd struct {
	Flag struct {
		globalFlag
		SearchDir string
		Key       string
		ResultDir string // save the search result temporary
		Mode      string // search mode
		Limit     int    // limit the number of result
		Ext       string
		OutExt    string
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
}

func (t *searchCmd) loadFileList(dir string) (list []string, err error) {
	if err := filepath.WalkDir(dir, func(file string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		ext := strings.Trim(path.Ext(d.Name()), ".")
		if ext == t.Flag.Ext {
			log.Debug("detect: ", file)
			list = append(list, file)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return
}

func (t *searchCmd) findVideoFile(fileList []string, videoName string) (file string, err error) {
	for _, item := range fileList {
		baseName := path.Base(item)
		ext := path.Ext(baseName)
		baseName = strings.TrimRight(baseName, ext)
		if strings.Contains(videoName, baseName) {
			return item, nil
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
		//base := path.Base(item.VideoFile)
		//ext := path.Ext(base)
		//name := strings.TrimRight(base, ext)
		name := getFileNameWithoutExt(item.Subtitle.FileName)
		start, err := lib.Time.Parse(lib.TimeInMillionLayout, item.Block.StartTime)
		if err != nil {
			log.Error(err)
			continue
		}
		end, err := lib.Time.Parse(lib.TimeInMillionLayout, item.Block.EndTime)
		if err != nil {
			log.Error(err)
			continue
		}
		duration := srt.ZeroTime.Add(end.Sub(start)).Format(timeLayout)
		out := path.Join(t.outDir, fmt.Sprintf("%s.split.%d.%s", name, item.Block.BlockId, t.Flag.OutExt))
		script += fmt.Sprintf(
			"if [ ! -f %s ]; then ffmpeg -y -i '%s' -ss %s -t %s '%s'; fi\n",
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
			"ffmpeg -y -i '%s' -vf subtitles=%s:force_style=\\'FontSize=16\\' '%s'\n",
			splitFile,
			srtFile,
			out,
		)
	}
	// padding
	outList := make([]string, 0)
	for _, item := range result.List {
		name := getFileNameWithoutExt(item.Subtitle.FileName)
		inFile := path.Join(t.outDir, fmt.Sprintf("%s.subtitle.%d.%s", name, item.Block.BlockId, t.Flag.OutExt))
		out := path.Join(t.outDir, fmt.Sprintf("%s.pad.%d.%s", name, item.Block.BlockId, t.Flag.OutExt))
		if t.hasChinese(out) {
			if t.hasChinese(out) {
				outList = append(outList, out)
			}
			script += fmt.Sprintf(
				"ffmpeg -y -i '%s' -vf 'pad=%s' '%s'\n",
				inFile,
				t.getPadConfig(),
				out,
			)
		}
	}
	// concat
	mergedVideo := path.Join(t.outDir, "merged.mp4")
	//script += fmt.Sprintf(
	//	"ffmpeg -y -i 'concat:%s' -vf select=concatdec_select -af aselect=concatdec_select,aresample=async=1 '%s'\n",
	//	strings.Join(outList, "|"),
	//	mergedVideo,
	//)
	script += fmt.Sprintf(
		"ffmpeg -y -i 'concat:%s' '%s'\n",
		strings.Join(outList, "|"),
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
		videoFile, err := t.findVideoFile(fileList, item.Movie.Name)
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
		srtFormatData := fmt.Sprintf("%d\n%s --> %s\n%s", beginningBlock.BlockId, beginningBlock.StartTime, beginningBlock.EndTime, beginningBlock.Content)
		if err := ioutil.WriteFile(outSrt, []byte(srtFormatData), 0644); err != nil {
			log.Fatal(err)
		}
	}
	// generate script
	t.generateSplitScript(result)
	log.Info("save search result in: " + out)
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
