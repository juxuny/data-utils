package model

import "time"

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

type EngSubtitleBlockList []EngSubtitleBlock
