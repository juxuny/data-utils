package proxy

import (
	"github.com/juxuny/data-utils/model"
	"github.com/juxuny/data-utils/proxy/dt"
)

type saveServerListHandler struct {
	db *model.DB
}

func NewSaveServerListHandler(db *model.DB) *saveServerListHandler {
	ret := &saveServerListHandler{db: db}
	return ret
}

func (t *saveServerListHandler) SaveServerList(list dt.ServerList) error {
	return SaveProxyData(t.db, list)
}
