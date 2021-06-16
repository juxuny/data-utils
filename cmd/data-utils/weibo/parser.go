package weibo

import (
	data_utils "github.com/juxuny/data-utils"
	"github.com/juxuny/data-utils/model"
)

type Parser interface {
	CheckValid(metaData MetaData) (isOk bool)
	Prepare(db *model.DB) error
	Parse(metaData MetaData) (data_utils.JobList, error)
}

var parserMap = map[Parser]struct{}{
	NewFriendshipParser(): struct{}{},
}
