package model

import (
	"time"
)

type ResType int

const (
	ResTypeNormal = ResType(1) // normal movie resource
	ResTypeTv     = ResType(2) // TV series
)

type Status int

const (
	StatusEnable  = Status(1)
	StatusDisable = Status(0)
)

type EngMovie struct {
	Id         int64      `json:"id" gorm:"type:int(11);primary_key;auto_increment"`
	Name       string     `json:"name" gorm:"type:varchar(100)"`
	ParentId   int64      `json:"parentId" gorm:"type:int(11)"`
	ResType    ResType    `json:"resType" gorm:"type:int(11)"`
	CreateTime *time.Time `json:"createTime" gorm:"createTime;default"`
	Status     Status     `json:"status" gorm:"type:tinyint(2)"`
	DeletedAt  *time.Time `json:"deletedAt" gorm:"type:timestamp;default"`
}

type EngMovieList []EngMovie
