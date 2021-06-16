package weibo

import (
	"encoding/json"
	data_utils "github.com/juxuny/data-utils"
	"github.com/juxuny/data-utils/log"
	"github.com/juxuny/data-utils/model"
	"github.com/pkg/errors"
)

type jobHandler struct {
	metaData MetaData
	config   data_utils.QueueConfig
}

func (j *jobHandler) JobType() model.JobType {
	return model.JobTypeWeibo
}

func (j *jobHandler) ParseMetaData(data string) error {
	return json.Unmarshal([]byte(data), &j.metaData)
}

func (j *jobHandler) Run() (data_utils.JobList, error) {
	log.Info(j.metaData.Url)
	var ret = make(data_utils.JobList, 0)
	db, err := model.Open(j.config.DbConfig)
	if err != nil {
		log.Error(err)
		return nil, errors.Wrap(err, "connect database failed")
	}
	defer func() {
		_ = db.Close()
	}()
	for parser := range parserMap {
		if parser.CheckValid(j.metaData) {
			if err := parser.Prepare(db); err != nil {
				log.Errorf("prepare parser failed: %v", err)
				continue
			}
			list, err := parser.Parse(j.metaData)
			if err != nil {
				log.Error(err)
				continue
			}
			if len(list) > 0 {
				ret = append(ret, list...)
			}
		}
	}
	return ret, nil
}

func NewJobHandler(config data_utils.QueueConfig) *jobHandler {
	h := &jobHandler{
		config: config,
	}
	return h
}

type handlerBuilder struct {
	config data_utils.QueueConfig
}

func (t *handlerBuilder) JobType() model.JobType {
	return model.JobTypeWeibo
}

func NewHandlerBuilder(config data_utils.QueueConfig) data_utils.HandlerBuilder {
	return &handlerBuilder{
		config: config,
	}
}

func (t *handlerBuilder) New() (data_utils.JobHandler, error) {
	return NewJobHandler(t.config), nil
}
