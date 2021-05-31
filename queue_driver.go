package data_utils

import "github.com/juxuny/data-utils/model"

type QueueDriverType string

const (
	QueueDriverTypeMysql = QueueDriverType("mysql")
	QueueDriverTypeRedis = QueueDriverType("redis")
)

func (t QueueDriverType) ToString() string {
	return string(t)
}

type queueDriver interface {
	GetById(ids ...int64) (jobList []model.Job, err error)
	UpdateState(state model.JobState, ids ...int64) (err error)
	Dequeue(num int, jobType ...model.JobType) (list JobList, err error)
	Enqueue(jobList ...Job) (err error)
}

func NewQueueDriver(config QueueConfig) (ret queueDriver, err error) {
	switch config.DriverType {
	case QueueDriverTypeMysql:
		ret, err = NewMysqlDriver(config.DbConfig)
	default:
		ret, err = NewMysqlDriver()
	}
	return ret, err
}
