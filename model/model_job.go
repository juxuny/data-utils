package model

import "time"

type Job struct {
	Id        int64      `json:"id" gorm:"TYPE:BIGINT(21);PRIMARY_KEY;AUTO_INCREMENT"`
	JobType   JobType    `json:"job_type" gorm:"TYPE:tinyint(2)"`
	State     JobState   `json:"state" gorm:"TYPE:tinyint(2)"`
	CreatedAt *time.Time `json:"created_at" gorm:"TYPE:TIMESTAMP;DEFAULT"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"TYPE:TIMESTAMP;DEFAULT"`
	StartedAt *time.Time `json:"started_at" gorm:"TYPE:TIMESTAMP"`
	EndAt     *time.Time `json:"end_at" gorm:"TYPE:TIMESTAMP"`
	MetaData  string     `json:"meta_data" gorm:"TYPE:TEXT"`
	Result    string     `json:"result" gorm:"TYPE:TEXT"`
}
