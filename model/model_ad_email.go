package model

import "time"

type AdEmail struct {
	Id              int64      `json:"id" gorm:"type:bigint(20);PRIMARY_KEY;AUTO_INCREMENT"`
	Email           string     `json:"email"`
	Count           int64      `json:"count"`
	CreatedAt       *time.Time `json:"created_at" gorm:"TYPE:TIMESTAMP;DEFAULT"`
	UpdatedAt       *time.Time `json:"updated_at" gorm:"TYPE:TIMESTAMP;DEFAULT"`
	LastError       string     `json:"last_error" gorm:"TYPE:TEXT"`
	LastErrorTime   *time.Time `json:"last_error_time" gorm:"TYPE:TIMESTAMP;DEFAULT"`
	LastSuccessTime *time.Time `json:"last_success_time" gorm:"TYPE:TIMESTAMP;DEFAULT"`
}

func (AdEmail) TableName() string {
	return "ad_email"
}

type AdEmailList []AdEmail

func (t AdEmailList) GetIdList() []int64 {
	ret := make([]int64, len(t))
	for i, item := range t {
		ret[i] = item.Id
	}
	return ret
}

func (t AdEmailList) GetEmailList() []string {
	ret := make([]string, len(t))
	for i, item := range t {
		ret[i] = item.Email
	}
	return ret
}
