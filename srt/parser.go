package srt

import (
	"fmt"
	"github.com/jinzhu/now"
	"github.com/juxuny/data-utils/lib"
	"github.com/pkg/errors"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
	"unicode"
)

const IntervalFormat = "15:04:05.000"

var ZeroTime = now.BeginningOfDay()

type Block struct {
	Id        int64
	StartTime time.Time
	EndTime   time.Time
	Subtitle  NodeList
}

func (t Block) String() string {
	ret := make([]string, 0)
	ret = append(ret, fmt.Sprintf("%v", t.Id))
	ret = append(ret, fmt.Sprintf("%s --> %s", t.StartTime.Format(IntervalFormat), t.EndTime.Format(IntervalFormat)))
	for _, sub := range t.Subtitle {
		ret = append(ret, sub.String())
	}
	return strings.Join(ret, "\n")
}

func (t Block) Content() string {
	ret := make([]string, 0)
	for _, sub := range t.Subtitle {
		ret = append(ret, sub.String())
	}
	return strings.Join(ret, "\n")
}

func parseInterval(data string) (startTime, endTime time.Time, err error) {
	bs := strings.Split(data, "-->")
	if len(bs) != 2 {
		err = fmt.Errorf("invalid interval: %v", data)
		return
	}
	parseTime := func(timeString string) (time.Time, error) {
		var ret = ZeroTime
		s := strings.ReplaceAll(timeString, ":", ",")
		l := strings.Split(s, ",")
		var nums []int
		for _, item := range l {
			item = strings.TrimLeft(item, "0 ")
			if item == "" {
				item = "0"
			}
			item = strings.TrimSpace(item)
			if v, err := strconv.ParseInt(item, 10, 64); err != nil {
				return ret, errors.Wrapf(err, "invalid number: %s", item)
			} else {
				nums = append(nums, int(v))
			}
		}
		var millionSeconds int
		if len(nums) >= 3 {
			millionSeconds = nums[0]*60*60*1000 + nums[1]*60*1000 + nums[2]*1000
		}
		if len(nums) >= 4 {
			millionSeconds += nums[3]
		}
		ret = ret.Add(time.Millisecond * time.Duration(millionSeconds))
		return ret, nil
	}
	startTime, err = parseTime(bs[0])
	if err != nil {
		return
	}
	endTime, err = parseTime(bs[1])
	if err != nil {
		return
	}
	return
}

func parseBlock(blockData string) (Block, error) {
	var err error
	var ret Block
	splitCutSet := "\n"
	lines := strings.Split(blockData, splitCutSet)
	//lines = lib.StringSlice(lines).Filter(func(item string) bool {
	//	return item != ""
	//})
	if len(lines) > 0 {
		if v, err := strconv.ParseInt(lines[0], 10, 64); err != nil {
			return ret, errors.Wrapf(err, "invalid block id: %v", lines[0])
		} else {
			ret.Id = v
		}
	}
	if len(lines) > 1 {
		ret.StartTime, ret.EndTime, err = parseInterval(lines[1])
		if err != nil {
			return ret, errors.Wrapf(err, "parse block failed: id=%v", ret.Id)
		}
	}
	if len(lines) > 2 {
		subtitleData := strings.Join(lines[2:], "\n")
		ret.Subtitle, err = ParseSubtitle([]byte(subtitleData))
		if err != nil {
			return ret, errors.Wrapf(err, "parse xml failed, id=%v", ret.Id)
		}
	}
	return ret, nil
}

func preHandle(data []byte) []byte {
	data = lib.Byte.Drop(data, func(r rune) bool {
		return r != 0
	})
	i := 0
	for i < len(data) {
		if unicode.IsDigit(rune(data[i])) {
			break
		}
		i += 1
	}
	return data[i:]
}

/*
修复这种情况：
1927
00:43:55,100 --> 00:43:57,420
<font face="Pingfang SC" size="20"><b><font size="18"><font color="#ffffff">阿
奇
伯
德
·
格
雷
西</font></font></b></font>

1928
00:43:55,100 --> 00:43:57,420
<font face="Pingfang SC" size="20"><b><font size="16"><font color="#ffffff">上
校

作
家</font></font></b></font>

*/
func fix(blocks []string) []string {
	ret := make([]string, 0)
	buf := blocks[0]
	i := 1
	for i < len(blocks) {
		if unicode.IsNumber(rune(blocks[i][0])) {
			ret = append(ret, buf)
			buf = blocks[i]
		} else {
			buf += "\n\n" + blocks[i]
		}
		i++
	}
	if buf != "" {
		ret = append(ret, buf)
	}
	return ret
}

func Parse(data []byte) ([]Block, error) {
	data = preHandle(data)
	data = []byte(strings.ReplaceAll(string(data), "\r\n", "\n"))
	splitCharset := "\n\n"
	//if strings.Contains(string(data), "\r\n") {
	//	splitCharset = "\r\n\r\n"
	//}
	blocks := strings.Split(string(data), splitCharset)
	blocks = lib.StringSlice(blocks).Filter(func(item string) bool {
		return item != ""
	})
	blocks = fix(blocks)
	ret := make([]Block, 0)
	for index, item := range blocks {
		if b, err := parseBlock(item); err == nil {
			ret = append(ret, b)
		} else {
			return nil, errors.Wrapf(err, "parse block failed, index=%d", index)
		}
	}
	return ret, nil
}

func ParseFile(file string) ([]Block, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, errors.Wrap(err, "load srt file failed")
	}
	return Parse(data)
}
