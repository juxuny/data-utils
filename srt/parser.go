package srt

import (
	"fmt"
	"github.com/jinzhu/now"
	"github.com/juxuny/data-utils/lib"
	"github.com/pkg/errors"
	"strconv"
	"strings"
	"time"
)

type Block struct {
	Id        int64
	StartTime time.Time
	EndTime   time.Time
	Subtitle  []*Node
}

func (t Block) String() string {
	ret := make([]string, 0)
	ret = append(ret, fmt.Sprintf("%v", t.Id))
	ret = append(ret, fmt.Sprintf("%s --> %s", t.StartTime.Format("15:04:05.000"), t.EndTime.Format("15:04:05.000")))
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
		var ret = now.BeginningOfDay()
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
	lines := strings.Split(blockData, "\n")
	lines = lib.StringSlice(lines).Filter(func(item string) bool {
		return item != ""
	})
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
		ret.Subtitle, err = ParseSubtitle([]byte(lines[2]))
		if err != nil {
			return ret, errors.Wrapf(err, "parse xml failed, id=%v", ret.Id)
		}
	}
	return ret, nil
}

func Parse(data []byte) ([]Block, error) {
	blocks := strings.Split(string(data), "\n\n")
	blocks = lib.StringSlice(blocks).Filter(func(item string) bool {
		return item != ""
	})
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
