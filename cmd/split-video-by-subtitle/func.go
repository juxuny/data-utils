package main

import (
	"encoding/json"
	"fmt"
	"github.com/juxuny/data-utils/canvas"
	"github.com/juxuny/data-utils/dict"
	"github.com/juxuny/data-utils/log"
	"github.com/juxuny/data-utils/srt"
	"github.com/pkg/errors"
	"image"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

const intervalLayout = "15:04:05.000"
const timeLayout = "15:04:05"

var ZeroTime, _ = time.Parse(intervalLayout, "00:00:00.000")

type CetFilter func(content string) (words []dict.Word)

type SplitConfigData struct {
	Id        int64               `json:"id"`
	StartTime string              `json:"startTime"`
	EndTime   string              `json:"endTime"`
	Words     map[string]struct{} `json:"words"`
}

// 把时间太接近的片段合并
func mergeSplitConfigData(configMap map[int64]*SplitConfigData) (ret []*SplitConfigData) {
	input := make([]*SplitConfigData, 0)
	maxId := int64(0)
	minId := int64(-1)
	for id := range configMap {
		if minId == -1 {
			minId = id
		}
		if id < minId {
			minId = id
		}
		if id > maxId {
			maxId = id
		}
	}
	for i := minId; i <= maxId; i++ {
		if v, found := configMap[i]; found {
			input = append(input, v)
		}
	}
	//ret = make([]*SplitConfigData, 0)
	//i := 0
	//for i < len(input) {
	//	if i < len(input)-1 {
	//		endTime, err := time.Parse(intervalLayout, input[i].EndTime)
	//		if err != nil {
	//			log.Warn(err)
	//			continue
	//		}
	//		nextStart, err := time.Parse(intervalLayout, input[i+1].StartTime)
	//		if err != nil {
	//			log.Warn(err)
	//			continue
	//		}
	//		if endTime.Add(time.Second * time.Duration(splitFlag.ExpandSeconds)).After(nextStart.Add(-time.Second * time.Duration(splitFlag.ExpandSeconds))) {
	//			log.Info("merged: ", input[i+1].Id, " >> ", input[i].Id)
	//			merged := &SplitConfigData{
	//				Id:        input[i].Id,
	//				StartTime: input[i].StartTime,
	//				EndTime:   input[i+1].EndTime,
	//				Words:     input[i].Words,
	//			}
	//			for w := range input[i+1].Words {
	//				merged.Words[w] = struct{}{}
	//			}
	//			ret = append(ret, input[i])
	//			i += 1
	//		}
	//	} else {
	//		ret = append(ret, input[i])
	//	}
	//	i++
	//}
	return input
}

// 扩展时间范围，防止视频过短
func expandDuration(data []*SplitConfigData) {
	for _, item := range data {
		startTime, err := time.Parse(intervalLayout, item.StartTime)
		if err != nil {
			log.Warn(err)
			continue
		}
		endTime, err := time.Parse(intervalLayout, item.EndTime)
		if err != nil {
			log.Warn(err)
			continue
		}
		if int64(startTime.Sub(ZeroTime).Seconds()) < int64(splitFlag.ExpandSeconds) {
			continue
		}
		duration := int(endTime.Sub(startTime).Seconds())
		if duration > splitFlag.ExpandSeconds*2 {
			continue
		}
		startTime = startTime.Add(-time.Second * time.Duration(splitFlag.ExpandSeconds))
		endTime = endTime.Add(time.Second * time.Duration(splitFlag.ExpandSeconds))
		item.StartTime = startTime.Format(intervalLayout)
		item.EndTime = endTime.Format(intervalLayout)
	}
}

func runCommand(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	log.Debug("run command: ", cmd.String())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return err
	}
	return cmd.Wait()
}

func getSplitScriptFileName() string {
	ext := path.Ext(splitFlag.InputFile)
	return strings.TrimRight(splitFlag.InputFile, ext) + ".sh"
}

func getPadConfig() string {
	return fmt.Sprintf("%d:%d:0:(%d-ih)/2:black", splitFlag.Width, splitFlag.Height, splitFlag.Height)
}

func getOutPutDir() string {
	ext := path.Ext(splitFlag.InputFile)
	return strings.TrimRight(splitFlag.InputFile, ext) + ".d"
}

func generateSplitScript(srtFile string, data []*SplitConfigData) error {
	script := ""
	ext := path.Ext(splitFlag.InputFile)
	outScript := getSplitScriptFileName()
	outDir := strings.TrimRight(splitFlag.InputFile, ext) + ".d"
	cropVideo := path.Join(outDir, strings.TrimRight(path.Base(splitFlag.InputFile), ext)+".crop."+splitFlag.OutExt)
	//outVideo := path.Join(outDir, strings.TrimRight(path.Base(splitFlag.InputFile), ext)+"."+splitFlag.OutExt)
	if stat, err := os.Stat(outDir); os.IsNotExist(err) {
		if err := os.MkdirAll(outDir, 0755); err != nil {
			return errors.Wrap(err, "create output directory failed")
		}
	} else if !stat.IsDir() {
		log.Fatal(outDir + "is not a director")
	}
	// 添加字幕
	//script += fmt.Sprintf("if [  ! -f %s ]; then ", outVideo)
	//script += fmt.Sprintf("ffmpeg -y -i '%s' -vf subtitles='%s' '%s';", splitFlag.InputFile, srtFile, outVideo)
	//script += fmt.Sprintf("fi\n")

	script += fmt.Sprintf("set -e\n")

	// 增加黑边
	script += fmt.Sprintf("if [  ! -f %s ]; then ", cropVideo)
	script += fmt.Sprintf("ffmpeg -y -i '%s' -vf 'pad=%s' '%s';", splitFlag.InputFile, getPadConfig(), cropVideo)
	script += fmt.Sprintf("fi\n")

	// 增加字幕
	//script += fmt.Sprintf("if [  ! -f %s ]; then ", outVideo)
	//script += fmt.Sprintf("ffmpeg -y -i '%s' -vf subtitles=%s:force_style=\\'MarginV=80\\' '%s';", cropVideo, srtFile, outVideo)
	//script += fmt.Sprintf("fi\n")

	//generate split command
	for _, d := range data {
		outSplit := path.Join(outDir, strings.TrimRight(path.Base(splitFlag.InputFile), ext)+fmt.Sprintf(".%d", d.Id))
		start, err := time.Parse(intervalLayout, d.StartTime)
		if err != nil {
			log.Warn(err, "ignore block: ", d.Id)
			continue
		}
		end, err := time.Parse(intervalLayout, d.EndTime)
		if err != nil {
			log.Warn(err, "ignore block: ", d.Id)
			continue
		}
		duration := ZeroTime.Add(end.Sub(start)).Format(timeLayout)
		script += fmt.Sprintf(
			"ffmpeg -y -i '%s' -ss %s -t %s '%s'\n",
			cropVideo,
			start.Format(timeLayout),
			duration,
			outSplit+".crop.mp4",
		)
		script += fmt.Sprintf(
			"ffmpeg -y -i '%s' -vf subtitles=%s:force_style=\\'MarginV=80\\' '%s'\n", outSplit+".crop.mp4", outSplit+".srt", outSplit+".subtitle.mp4",
		)
		script += fmt.Sprintf(
			"ffmpeg -y -i '%s' -f mpegts '%s'\n",
			outSplit+".subtitle.mp4",
			outSplit+".crop.ts",
		)
		script += fmt.Sprintf(
			"ffmpeg -y -i 'concat:%s|%s' -vf select=concatdec_select -af aselect=concatdec_select,aresample=async=1 '%s'\n",
			outSplit+".begin.ts",
			strings.Trim(strings.Repeat(outSplit+".crop.ts|", splitFlag.Repeat), "|"),
			//outSplit+".end.ts",
			outSplit+"."+splitFlag.OutExt,
		)
	}

	if err := ioutil.WriteFile(outScript, []byte(script), 0755); err != nil {
		return errors.Wrap(err, "create script file failed")
	}
	return nil
}

// suffix 'begin' or 'end'.
func getCoverImageFile(blockId int64, suffix string) string {
	ext := path.Ext(splitFlag.InputFile)
	outDir := strings.TrimRight(splitFlag.InputFile, ext) + ".d"
	outImg := path.Join(outDir, strings.TrimRight(path.Base(splitFlag.InputFile), ext)+fmt.Sprintf(".%d.%s", blockId, suffix)+".jpg")
	return outImg
}

func generateCoverImageOne(outImg string, data SplitConfigData, dictData *dict.Dict) error {
	words := make([]dict.Word, 0)
	for w := range data.Words {
		if v, b := dictData.Data[w]; b {
			words = append(words, v)
		}
	}
	if _, err := os.Stat(outImg); os.IsNotExist(err) {
		log.Info("generate cover:", outImg)
		c := canvas.NewCanvas(splitFlag.CoverWidth, splitFlag.CoverHeight)
		//if strings.Index(splitFlag.CoverBg, "#") == 0 {
		//	if err := c.DrawColor(splitFlag.CoverBg); err != nil {
		//		log.Warn(err)
		//		return err
		//	}
		//} else {
		//	coverBgExt := path.Ext(splitFlag.CoverBg)
		//	imageType := canvas.ImageTypeJpeg
		//	if strings.ToLower(coverBgExt) == ".png" {
		//		imageType = canvas.ImageTypePng
		//	}
		//	bg := canvas.CreateImageView(splitFlag.CoverBg, c.Width, c.Height, imageType)
		//	if err := c.Draw(bg); err != nil {
		//		return errors.Wrapf(err, "generate cover failed, block id=%d", data.Id)
		//	}
		//}
		//_ = c.DrawColor("#000000")
		_ = c.DrawImageFromFile(image.Pt(0, 0), "tmp/image/sun.png")
		textPadding := 100
		lv := canvas.CreateListView(textPadding, textPadding, []canvas.View{})
		for _, w := range words {
			if err := lv.AppendChild(canvas.CreateTextView(w.Name, splitFlag.FontFile, splitFlag.CoverFontSize, splitFlag.CoverFontColor)); err != nil {
				log.Warn(err)
				continue
			}
			//if err := lv.AppendChild(canvas.CreateWrapTextView(w.Description, splitFlag.FontFile, splitFlag.DescriptionFontSize, splitFlag.DescriptionFontColor, splitFlag.Width-textPadding*2, nil)); err != nil {
			//	log.Warn(err)
			//	continue
			//}
		}
		var h = 810 // 电影视频源高度
		rect := lv.Measure()
		if rect.Dy() > h {
			h = rect.Dy()
		}
		log.Debug(h)
		h += textPadding * 2
		paddingHeight := (splitFlag.Height - h) / 2

		box := canvas.CreateCenterLayout(canvas.CenterTypeHorizontal, image.Pt(0, paddingHeight), splitFlag.CoverWidth, h, lv)
		center := canvas.CreateCenterLayout(canvas.CenterTypeVertical, image.Pt(0, 0), splitFlag.CoverWidth, splitFlag.CoverHeight, box)
		if err := c.Draw(center); err != nil {
			log.Error(err)
			return errors.Wrapf(err, "render word failed, block id: %d", data.Id)
		}
		if err := c.Save(outImg, canvas.ImageTypeJpeg); err != nil {
			log.Error(err)
			return errors.Wrap(err, "save cover img failed: "+outImg)
		}
	} else {
		log.Info("ignore begin cover: ", outImg)
	}
	return nil
}

// 生成重点单词汇总图片
func generateCoverImage(dictData *dict.Dict) error {
	splitDataFile := path.Dir(splitFlag.OutSrt) + string(os.PathSeparator) + path.Base(splitFlag.OutSrt) + ".split.json"
	jsonData, err := ioutil.ReadFile(splitDataFile)
	if err != nil {
		return errors.Wrap(err, "load split json failed")
	}
	var splitData []*SplitConfigData
	if err := json.Unmarshal(jsonData, &splitData); err != nil {
		return errors.Wrap(err, "parse split json failed")
	}
	//ext := path.Ext(splitFlag.InputFile)
	//outDir := strings.TrimRight(splitFlag.InputFile, ext) + ".d"
	for _, d := range splitData {
		//words := make([]dict.Word, 0)
		//for w := range d.Words {
		//	if v, b := dictData.Data[w]; b {
		//		words = append(words, v)
		//	}
		//}
		// 生成开头的封面
		outImg := getCoverImageFile(d.Id, "begin")
		if err := generateCoverImageOne(outImg, *d, dictData); err != nil {
			log.Error(err)
			continue
		}

		// 图片转ts视频
		if err := convertImageToVideo(outImg); err != nil {
			log.Error(err)
			return errors.Wrap(err, "convert image to video failed")
		}

		//// 生成结尾图片
		//outImg = getCoverImageFile(d.Id, "end")
		//if _, err := os.Stat(outImg); os.IsNotExist(err) {
		//	log.Info("generate cover:", outImg)
		//	c := canvas.NewCanvas(splitFlag.Width, splitFlag.Height)
		//	if strings.Index(splitFlag.CoverBg, "#") == 0 {
		//		if err := c.DrawColor(splitFlag.CoverBg); err != nil {
		//			log.Warn(err)
		//			continue
		//		}
		//	} else {
		//		coverBgExt := path.Ext(splitFlag.CoverBg)
		//		imageType := canvas.ImageTypeJpeg
		//		if strings.ToLower(coverBgExt) == ".png" {
		//			imageType = canvas.ImageTypePng
		//		}
		//		bg := canvas.CreateImageView(splitFlag.CoverBg, c.Width, c.Height, imageType)
		//		if err := c.Draw(bg); err != nil {
		//			return errors.Wrapf(err, "generate cover failed, block id=%d", d.Id)
		//		}
		//	}
		//	lv := canvas.CreateListView(427+50, 978+50, []canvas.View{})
		//	for _, w := range words {
		//		if err := lv.AppendChild(canvas.CreateTextView(w.Name, splitFlag.FontFile, splitFlag.CoverFontSize, splitFlag.CoverFontColor)); err != nil {
		//			log.Warn(err)
		//			continue
		//		}
		//		//log.Debug("descFontSize:", splitFlag.DescriptionFontSize, " color:", splitFlag.DescriptionFontColor)
		//		if err := lv.AppendChild(canvas.CreateWrapTextView(w.Description, splitFlag.FontFile, splitFlag.DescriptionFontSize, splitFlag.DescriptionFontColor, 1380, nil)); err != nil {
		//			log.Warn(err)
		//			continue
		//		}
		//
		//		if err := lv.AppendChild(canvas.CreateBox(image.Rect(0, 0, 50, 50), nil)); err != nil {
		//			log.Warn(err)
		//			continue
		//		}
		//	}
		//	if err := c.Draw(lv); err != nil {
		//		log.Error(err)
		//		return errors.Wrapf(err, "render word failed, block id: %d", d.Id)
		//	}
		//	if err := c.Save(outImg, canvas.ImageTypeJpeg); err != nil {
		//		log.Error(err)
		//		return errors.Wrap(err, "save cover img failed: "+outImg)
		//	}
		//} else {
		//	log.Info("ignore end cover: ", outImg)
		//}
		//
		//// 结尾图片转视频
		//if err := convertImageToVideo(outImg); err != nil {
		//	log.Error(err)
		//	return errors.Wrap(err, "convert image to video failed")
		//}
	}
	return nil
}

func convertImageToVideo(imageFile string) error {
	ext := path.Ext(imageFile)
	outVideo := strings.TrimRight(imageFile, ext) + ".ts"
	if _, err := os.Stat(outVideo); os.IsNotExist(err) {
		return runCommand("ffmpeg",
			"-y", "-r", "25", "-loop", "1", "-i", imageFile,
			"-r", "25", "-t", fmt.Sprintf("%d", splitFlag.CoverDuration), outVideo,
		)
	}
	return nil
}

// 随机选择单词数组里的几个元素
func randomSampleWords(words []dict.Word, num int) []dict.Word {
	if len(words) == 0 {
		return []dict.Word{}
	}
	index := rand.Perm(len(words))
	i := 0
	ret := make([]dict.Word, 0)
	for i < num && i < len(words) {
		ret = append(ret, words[index[i]])
		i++
	}
	return ret
}

func generateSplitSrt(block srt.Block) {
	outDir := getOutPutDir()
	if _, err := os.Stat(outDir); os.IsNotExist(err) {
		_ = os.MkdirAll(outDir, 0766)
	}
	ext := path.Ext(splitFlag.InputFile)
	splitSrt := path.Join(outDir, strings.TrimRight(path.Base(splitFlag.InputFile), ext)+fmt.Sprintf(".%d.srt", block.Id))
	var detail = block.StartTime.Sub(srt.ZeroTime)
	detail -= detail % time.Second
	block.Id = 1
	block.StartTime = block.StartTime.Add(-detail)
	block.EndTime = block.EndTime.Add(-detail)
	if err := ioutil.WriteFile(splitSrt, []byte(block.String()), 0655); err != nil {
		log.Error(err)
	}
}

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
	splitDataMap := make(map[int64]*SplitConfigData)
	for _, b := range blocks {
		b.Subtitle.ForEach(func(n *srt.Node) bool {
			if n.Type == srt.NodeTypeTag {
				if _, b := n.Attr["size"]; b {
					n.Attr["size"] = fmt.Sprintf("%d", splitFlag.GlobalFontSize)
				}
			}
			if n.Type != srt.NodeTypeText {
				return false
			}
			words := filter(string(n.Content))
			if len(words) == 0 {
				return false
			}
			words = randomSampleWords(words, splitFlag.MaxWords)
			log.Debug("highlight words:", len(words))
			wordMap := make(map[string]struct{})
			for _, w := range words {
				wordMap[w.Name] = struct{}{}
			}
			if v, found := splitDataMap[b.Id]; !found {
				splitDataMap[b.Id] = &SplitConfigData{
					Id:        b.Id,
					StartTime: b.StartTime.Format("15:04:05.000"),
					EndTime:   b.EndTime.Format("15:04:05.000"),
					Words:     wordMap,
				}
			} else {
				for w := range wordMap {
					v.Words[w] = struct{}{}
				}
			}
			for _, w := range words {
				log.Info("found word: ", w.Name)
				//n.Content = []byte(strings.ReplaceAll(string(n.Content), w.Name, fmt.Sprintf("<font size=\"48\"><b>%s</b></font>", w.Name)))
				n.Content = dict.ReplaceWords(
					n.Content,
					[]byte(w.Name),
					[]byte(fmt.Sprintf(
						"<font size=\"%d\" color=\"%s\"><b>%s</b></font>",
						splitFlag.HighLightFontSize,
						splitFlag.HighlightColor,
						w.Name,
					)),
				)
			}
			return false
		})
		if _, found := splitDataMap[b.Id]; found {
			generateSplitSrt(b)
		}
	}

	mergedSplitData := mergeSplitConfigData(splitDataMap)
	expandDuration(mergedSplitData)
	// save block need to split
	splitDataFile := path.Dir(outFile) + string(os.PathSeparator) + path.Base(outFile) + ".split.json"
	jsonData, _ := json.Marshal(mergedSplitData)
	if err := ioutil.WriteFile(splitDataFile, jsonData, 0664); err != nil {
		return errors.Wrap(err, "convert srt file failed")
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

	// generate video split script
	if err := generateSplitScript(outFile, mergedSplitData); err != nil {
		return errors.Wrap(err, "convert srt file failed, when generate split script")
	}
	log.Info("generate script finished")
	return nil
}

func loadFileList(dir string, ext string) (ret []string, err error) {
	fileList, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, item := range fileList {
		if item.IsDir() {
			continue
		}

		if strings.Trim(path.Ext(item.Name()), ".") == strings.Trim(ext, ".") {
			ret = append(ret, path.Join(dir, item.Name()))
		}
	}
	return ret, nil
}

func getFileNameWithoutExt(file string) string {
	base := path.Base(file)
	ext := path.Ext(file)
	return strings.TrimRight(base, ext)
}
