package model

import "time"

type EngMovie struct {
	Id         int64      `json:"id" gorm:"type:int(11);primary_key;auto_increment"`
	Name       string     `json:"name" gorm:"type:varchar(100)"`
	ParentId   int64      `json:"parentId" gorm:"type:int(11)"`
	CreateTime *time.Time `json:"createTime" gorm:"createTime;default"`
}

type EngMovieList []EngMovie
