package data_utils

import (
	"github.com/juxuny/data-utils/log"
	"github.com/juxuny/data-utils/model"
	"github.com/pkg/errors"
)

type Job struct {
	Id       int64         `json:"id"`
	JobType  model.JobType `json:"job_type"`
	MetaData string        `json:"meta_data"`
}

type JobList []Job

type JobHandler interface {
	JobType() model.JobType
	ParseMetaData(data string) error
	Run() (JobList, error)
}

type QueueConfig struct {
	DriverType QueueDriverType `json:"driver_type"`
	DbConfig   model.Config
}

type queue struct {
	config      QueueConfig
	handlerSet  map[JobHandler]struct{}
	queueDriver queueDriver
	logger      log.ILogger
}

func (t *queue) initDriver() error {
	var err error
	if t.queueDriver == nil {
		t.queueDriver, err = NewQueueDriver(t.config)
	}
	return err
}

func (t *queue) StartDaemon() error {
	if err := t.initDriver(); err != nil {
		t.logger.Error(err)
		return errors.Wrap(err, "init queue driver failed")
	}
	return nil
}

func (t *queue) Enqueue(jobList ...Job) (err error) {
	if err := t.initDriver(); err != nil {
		t.logger.Error(err)
		return errors.Wrap(err, "init queue driver failed")
	}
	return t.queueDriver.Enqueue(jobList...)
}

func NewQueue(config QueueConfig) *queue {
	q := &queue{
		logger:     log.NewLogger("queue"),
		config:     config,
		handlerSet: make(map[JobHandler]struct{}),
	}
	return q
}

func (t *queue) RegisterHandler(j JobHandler) {
	t.handlerSet[j] = struct{}{}
}
