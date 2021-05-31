package data_utils

import (
	"github.com/jinzhu/gorm"
	"github.com/juxuny/data-utils/log"
	"github.com/juxuny/data-utils/model"
	"github.com/pkg/errors"
)

type queueMysqlDriver struct {
	db     *model.DB
	logger log.ILogger
}

func (t *queueMysqlDriver) GetById(ids ...int64) (jobList []model.Job, err error) {
	if len(ids) == 0 {
		return []model.Job{}, nil
	}
	if err := t.db.Where("id IN (?)", ids).Find(&jobList).Error; err != nil {
		if !model.IsErrNoDataInDb(err) {
			t.logger.Error(err)
			return nil, errors.Wrap(err, "read table failed: job")
		}
	}
	return
}

func (t *queueMysqlDriver) UpdateState(state model.JobState, ids ...int64) (err error) {
	if len(ids) == 0 {
		return nil
	}
	if err := t.db.Model(&model.Job{}).Where("in IN (?)", ids).Updates(map[string]interface{}{
		"state": state,
	}).Error; err != nil {
		t.logger.Error(err)
		return errors.Wrap(err, "update table failed: job")
	}
	return nil
}

func (t *queueMysqlDriver) Dequeue(num int, jobType ...model.JobType) (list JobList, err error) {
	var records []model.Job
	db := t.db.Limit(num)
	if len(jobType) > 0 {
		db = db.Where("job_type IN (?) AND state = ?", jobType, model.JobStateWaiting)
	}
	if err := db.Find(&records).Error; err != nil {
		if !model.IsErrNoDataInDb(err) {
			t.logger.Error(err)
			return nil, errors.Wrap(err, "read table failed: job")
		}
	}
	for _, item := range records {
		list = append(list, Job{
			Id:       item.Id,
			MetaData: item.MetaData,
		})
	}
	return
}

func (t *queueMysqlDriver) Enqueue(jobList ...Job) (err error) {
	return t.db.Begin(func(db *gorm.DB) error {
		for _, item := range jobList {
			if err := db.Create(&model.Job{
				JobType:  item.JobType,
				State:    model.JobStateWaiting,
				MetaData: item.MetaData,
			}).Error; err != nil {
				t.logger.Error(err)
				return errors.Wrap(err, "write table failed: job")
			}
		}
		return nil
	})
}

func NewMysqlDriver(config ...model.Config) (driver queueDriver, err error) {
	logger := log.NewLogger("mysql-driver")
	var finalConfig model.Config
	if len(config) > 0 {
		finalConfig = config[0]
	} else {
		finalConfig, err = model.GetEnvConfig()
		if err != nil {
			logger.Error(err)
			return nil, errors.Wrap(err, "get database config from environment failed")
		}
	}
	db, err := model.Open(finalConfig)
	if err != nil {
		logger.Error(err)
		return nil, errors.Wrap(err, "create database connect failed")
	}
	ret := &queueMysqlDriver{
		db: db,
	}
	return ret, nil
}
