package main

import (
	"encoding/json"
	"fmt"
	"github.com/juxuny/data-utils/canvas"
	"github.com/juxuny/data-utils/dict"
	"github.com/juxuny/data-utils/log"
	"github.com/juxuny/data-utils/srt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

const intervalLayout = "15:04:05.000"
const timeLayout = "15:04:05"

var ZeroTime, _ = time.Parse(intervalLayout, "00:00:00.000")

var globalFlag = struct {
	Verbose bool
}{}

func initGlobalFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVarP(&globalFlag.Verbose, "verbose", "v", false, "display debug output")

}

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
	ret = make([]*SplitConfigData, 0)
	i := 0
	for i < len(input) {
		if i < len(input)-1 {
			endTime, err := time.Parse(intervalLayout, input[i].EndTime)
			if err != nil {
				log.Warn(err)
				continue
			}
			nextStart, err := time.Parse(intervalLayout, input[i+1].StartTime)
			if err != nil {
				log.Warn(err)
				continue
			}
			if endTime.Add(time.Second * time.Duration(splitFlag.ExpandSeconds)).After(nextStart.Add(-time.Second * time.Duration(splitFlag.ExpandSeconds))) {
				log.Info("merged: ", input[i+1].Id, " >> ", input[i].Id)
				merged := &SplitConfigData{
					Id:        input[i].Id,
					StartTime: input[i].StartTime,
					EndTime:   input[i+1].EndTime,
					Words:     input[i].Words,
				}
				for w := range input[i+1].Words {
					merged.Words[w] = struct{}{}
				}
				ret = append(ret, input[i])
				i += 1
			}
		} else {
			ret = append(ret, input[i])
		}
		i++
	}
	return
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
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func getSplitScriptFileName() string {
	ext := path.Ext(splitFlag.InputFile)
	return strings.TrimRight(splitFlag.InputFile, ext) + ".sh"
}

func generateSplitScript(srtFile string, data []*SplitConfigData) error {
	script := ""
	ext := path.Ext(splitFlag.InputFile)
	outScript := getSplitScriptFileName()
	outDir := strings.TrimRight(splitFlag.InputFile, ext) + ".d"
	outVideo := path.Join(outDir, strings.TrimRight(path.Base(splitFlag.InputFile), ext)+"."+splitFlag.OutExt)
	if stat, err := os.Stat(outDir); os.IsNotExist(err) {
		if err := os.MkdirAll(outDir, 0755); err != nil {
			return errors.Wrap(err, "create output directory failed")
		}
	} else if !stat.IsDir() {
		log.Fatal(outDir + "is not a director")
	}
	script += fmt.Sprintf("ffmpeg -y -i '%s' -vf subtitles='%s' '%s'\n", splitFlag.InputFile, srtFile, outVideo)

	//generate split command
	for _, d := range data {
		outSplit := path.Join(outDir, strings.TrimRight(path.Base(splitFlag.InputFile), ext)+fmt.Sprintf(".%d", d.Id)+"."+splitFlag.OutExt)
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
		script += fmt.Sprintf("ffmpeg -y -i '%s' -ss %s -t %s %s\n", outVideo, start.Format(timeLayout), duration, outSplit)
	}

	if err := ioutil.WriteFile(outScript, []byte(script), 0755); err != nil {
		return errors.Wrap(err, "create script file failed")
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
	ext := path.Ext(splitFlag.InputFile)
	outDir := strings.TrimRight(splitFlag.InputFile, ext) + ".d"
	for _, d := range splitData {
		words := make([]dict.Word, 0)
		for w := range d.Words {
			if v, b := dictData.Data[w]; b {
				words = append(words, v)
			}
		}
		outImg := path.Join(outDir, strings.TrimRight(path.Base(splitFlag.InputFile), ext)+fmt.Sprintf(".%d.begin", d.Id)+".jpg")
		log.Info("generate cover:", outImg)
		c := canvas.NewCanvas(750, 1206)
		if strings.Index(splitFlag.CoverBg, "#") == 0 {
			if err := c.DrawColor(splitFlag.CoverBg); err != nil {
				log.Warn(err)
				continue
			}
		} else {
			coverBgExt := path.Ext(splitFlag.CoverBg)
			imageType := canvas.ImageTypeJpeg
			if strings.ToLower(coverBgExt) == ".png" {
				imageType = canvas.ImageTypePng
			}
			bg := canvas.CreateImageView(splitFlag.CoverBg, c.Width, c.Height, imageType)
			if err := c.Draw(bg); err != nil {
				return errors.Wrapf(err, "generate cover failed, block id=%d", d.Id)
			}
		}
		lv := canvas.CreateListView(135+50, 380+50, []canvas.View{})
		for _, w := range words {
			if err := lv.AppendChild(canvas.CreateTextView(w.Name, splitFlag.FontFile, splitFlag.CoverFontSize, splitFlag.CoverFontColor)); err != nil {
				log.Warn(err)
				continue
			}
		}
		if err := c.Draw(lv); err != nil {
			log.Error(err)
			return errors.Wrapf(err, "render word failed, block id: %d", d.Id)
		}
		if err := c.Save(outImg, canvas.ImageTypeJpeg); err != nil {
			log.Error(err)
			return errors.Wrap(err, "save cover img failed: "+outImg)
		}
		if d.Id > 7 {
			break
		}
		break
	}
	return nil
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
						"<font face=\"%s\" size=\"%d\" color=\"%s\"><b>%s</b></font>",
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
