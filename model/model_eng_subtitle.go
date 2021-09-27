package model

import "time"

type EngSubtitle struct {
	Id         int64      `json:"id" gorm:"type:int(11);primary_key;auto_increment"`
	MovieId    int64      `json:"movieId" gorm:"type:int(11)"`
	Ext        string     `json:"ext" gorm:"type:varchar(20)"` // file extension
	FileName   string     `json:"fileName" gorm:"type:varchar(200)"`
	CreateTime *time.Time `json:"createTime" gorm:"type:timestamp;default"`
}

type EngSubtitleList []EngSubtitle
