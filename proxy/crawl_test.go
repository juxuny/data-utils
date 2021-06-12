package proxy

import (
	"github.com/juxuny/data-utils/model"
	"github.com/juxuny/data-utils/proxy/dt"
	"testing"
)

func TestSaveProxyData(t *testing.T) {
	db, err := model.Open(model.Config{
		DbHost:     "127.0.0.1",
		DbPort:     3307,
		DbUser:     "root",
		DbPassword: "123456",
		DbName:     "crawl",
		DbDebug:    true,
	})
	if err != nil {
		t.Fatal(err)
	}
	list := dt.ServerList{
		{Schema: dt.SchemaTypeHttp, Ip: "61.19.27.201", Port: 8080, Provider: model.Provider66Ip},
	}
	if err := SaveProxyData(db, list); err != nil {
		t.Fatal(err)
	}
}
