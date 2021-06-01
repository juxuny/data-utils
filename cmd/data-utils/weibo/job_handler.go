package weibo

import (
	"encoding/json"
	data_utils "github.com/juxuny/data-utils"
	"github.com/juxuny/data-utils/log"
	"github.com/juxuny/data-utils/model"
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
	return nil, nil
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
