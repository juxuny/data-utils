package data_utils

import (
	"github.com/juxuny/data-utils/log"
	"github.com/juxuny/data-utils/model"
	"github.com/pkg/errors"
	"runtime/debug"
	"time"
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

type HandlerBuilder interface {
	New() (JobHandler, error)
	JobType() model.JobType
}

type QueueConfig struct {
	DriverType QueueDriverType `json:"driver_type"`
	DbConfig   model.Config
	BatchSize  int
}

type queue struct {
	config           QueueConfig
	builderSet       map[HandlerBuilder]struct{}
	queueDriver      queueDriver
	logger           log.ILogger
	allowJobTypeList []model.JobType
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
	if t.config.BatchSize < 1 {
		return errors.Errorf("batch size cannot be 0")
	}
	var jobList JobList
	var err error
	for {
		for builder := range t.builderSet {
			jobType := builder.JobType()
			jobList, err = t.Dequeue(t.config.BatchSize, jobType)
			if err != nil {
				log.Error(err)
				continue
			}
			var newJobList JobList
			for _, job := range jobList {
				err = func() error {
					defer func() {
						if err := recover(); err != nil {
							log.Error(err)
							debug.PrintStack()
							return
						}
					}()
					handler, err := builder.New()
					if err != nil {
						log.Error(err)
						return errors.Wrap(err, "create handler failed")
					}
					if err := handler.ParseMetaData(job.MetaData); err != nil {
						log.Error(err)
						return errors.Errorf("parse meta data failed")
					}
					err = t.queueDriver.UpdateState(model.JobStateRunning, "", job.Id)
					if err != nil {
						log.Error(err)
						return errors.Wrap(err, "update state failed")
					}
					newJobList, err = handler.Run()
					if err != nil {
						log.Error(err)
						return errors.Wrap(err, "run failed")
					}
					err = t.Enqueue(newJobList...)
					if err != nil {
						log.Error(err)
						return errors.Wrap(err, "enqueue job failed")
					}
					return nil
				}()
				var newState model.JobState
				var result string
				if err != nil {
					newState = model.JobStateFailed
					result = err.Error()
				} else {
					newState = model.JobStateSucceed
				}
				err = t.queueDriver.UpdateState(newState, result, job.Id)
				if err != nil {
					log.Error(err)
				}
			}
		}
		time.Sleep(time.Second)
	}
}

func (t *queue) Dequeue(num int, jobType ...model.JobType) (jobList []Job, err error) {
	if t.queueDriver == nil {
		if err := t.initDriver(); err != nil {
			log.Error(err)
			return nil, errors.Wrap(err, "init queue driver failed")
		}
	}
	return t.queueDriver.Dequeue(num, jobType...)
}

func (t *queue) Enqueue(jobList ...Job) (err error) {
	if err := t.initDriver(); err != nil {
		t.logger.Error(err)
		return errors.Wrap(err, "init queue driver failed")
	}
	if err := t.queueDriver.Enqueue(jobList...); err != nil {
		log.Error(err)
		return errors.Wrap(err, "enqueue failed")
	}

	return
}

func NewQueue(config QueueConfig) *queue {
	q := &queue{
		logger:     log.NewLogger("queue"),
		config:     config,
		builderSet: make(map[HandlerBuilder]struct{}),
	}
	return q
}

func (t *queue) RegisterHandler(builder HandlerBuilder) {
	for k := range t.builderSet {
		if k.JobType() == builder.JobType() {
			panic("duplicated builder: " + builder.JobType())
		}
	}
	t.builderSet[builder] = struct{}{}
}
