package model

import (
	"github.com/jinzhu/now"
	"github.com/juxuny/data-utils/lib"
	"time"
)

func GetZero() time.Time {
	return now.BeginningOfDay()
}

type EngSubtitleBlock struct {
	Id             int64      `json:"id" gorm:"int(11);primary_key;auto_increment"`
	SubtitleId     int64      `json:"subtitleId" gorm:"type:int(11)"`
	BlockId        int64      `json:"blockId" gorm:"type:int(11)"`
	StartTime      string     `json:"startTime" gorm:"type:varchar(20)"`
	EndTime        string     `json:"endTime" gorm:"type:varchar(20)"`
	DurationExtend string     `json:"durationExtend" gorm:"type:varchar(200)"`
	Content        string     `json:"content" gorm:"type:text"`
	CreateTime     *time.Time `json:"createTime" gorm:"type:timestamp;default"`
}

func (t EngSubtitleBlock) MoveToBeginning() (ret EngSubtitleBlock, err error) {
	ret = t
	var zero = GetZero()
	layout := "2006-01-02 15:04:05.000"
	startTime, err := lib.Time.Parse(layout, zero.Format(lib.DayLayout)+" "+t.StartTime)
	if err != nil {
		return ret, err
	}
	endTime, err := lib.Time.Parse(layout, zero.Format(lib.DayLayout)+" "+t.EndTime)
	if err != nil {
		return ret, err
	}
	detail := startTime.Sub(zero)
	detail -= detail % time.Second
	ret.StartTime = startTime.Add(-detail).Format(lib.TimeInMillionLayout)
	ret.EndTime = endTime.Add(-detail).Format(lib.TimeInMillionLayout)
	//log.Debug(detail, startTime, endTime, " ", t.StartTime, " ", t.EndTime, " ", ret.StartTime, " ", ret.EndTime)
	return
}

func (t EngSubtitleBlock) ExpandSubtitleDuration(seconds int64) (ret EngSubtitleBlock, err error) {
	ret = t
	var zero = GetZero()
	layout := "2006-01-02 15:04:05.000"
	startTime, err := lib.Time.Parse(layout, zero.Format(lib.DayLayout)+" "+t.StartTime)
	if err != nil {
		return ret, err
	}
	startTime = startTime.Add(-time.Duration(seconds) * time.Second)
	if startTime.Before(zero) {
		startTime = zero
	}
	endTime, err := lib.Time.Parse(layout, zero.Format(lib.DayLayout)+" "+t.EndTime)
	if err != nil {
		return ret, err
	}
	endTime = endTime.Add(time.Duration(seconds) * time.Second)
	ret.StartTime = startTime.Format(lib.TimeInMillionLayout)
	ret.EndTime = endTime.Format(lib.TimeInMillionLayout)
	return
}

type EngSubtitleBlockList []EngSubtitleBlock
