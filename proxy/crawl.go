package proxy

import (
	"fmt"
	"github.com/juxuny/data-utils/log"
	"github.com/juxuny/data-utils/model"
	"github.com/juxuny/data-utils/proxy/dt"
	"github.com/juxuny/data-utils/proxy/fetcher"
	"github.com/pkg/errors"
	"strings"
)

type Fetcher interface {
	Len() int
	Page(page int, pageSize int) (ret dt.ServerList, err error)
	AllData() (ret dt.ServerList, err error)
	Init() (err error)
	Reset()
}

var fetcherMap = map[model.Provider]Fetcher{}

func InitFetcher(handler dt.SaveHandler) error {
	fetcherMap[model.Provider66Ip] = fetcher.New66Ip(handler)
	return nil
}

func Fetch(blockOnError bool, provider ...model.Provider) (ret dt.ServerList, err error) {
	if len(fetcherMap) == 0 {
		panic("call InitFetcher first")
	}
	var loadDataFunc = func(p model.Provider) error {
		if err := fetcherMap[p].Init(); err != nil {
			return errors.Wrap(err, "init fetcher failed: "+fmt.Sprintf("%v", p))
		}
		if list, err := fetcherMap[p].AllData(); err != nil {
			return errors.Wrap(err, "fetch all data failed: "+fmt.Sprintf("%v", p))
		} else if len(list) > 0 {
			ret = append(ret, list...)
		}
		return nil
	}
	if len(provider) == 0 {
		for p := range fetcherMap {
			if err := loadDataFunc(p); err != nil {
				log.Error(err)
				if blockOnError {
					return ret, errors.Wrap(err, "load data func error, provider: "+fmt.Sprintf("%v", p))
				}
			}
		}
	} else {
		for _, p := range provider {
			if err := loadDataFunc(p); err != nil {
				log.Error(err)
				if blockOnError {
					return ret, errors.Wrap(err, "load data func error, provider: "+fmt.Sprintf("%v", p))
				}
			}
		}
	}

	return ret, nil
}

func SaveProxyData(db *model.DB, data dt.ServerList) error {
	var saveBatch = func(batch dt.ServerList) error {
		insertStatement := fmt.Sprintf("INSERT IGNORE INTO %s (ip, port, provider_name, socks5, http, https)", model.Proxy{}.TableName())
		holders := strings.Trim(strings.Repeat("(?, ?, ?, ?, ?, ?), ", len(batch)), ", ")
		var values []interface{}
		for _, item := range batch {
			socks5 := 0
			http := 0
			https := 0
			switch item.Schema {
			case dt.SchemaTypeSocks5:
				socks5, http, https = 1, 0, 0
			case dt.SchemaTypeHttp:
				socks5, http, https = 0, 1, 1
			case dt.SchemaTypeHttps:
				socks5, http, https = 0, 0, 1
			}
			values = append(
				values,
				item.Ip, item.Port, item.Provider,
				socks5, http, https,
			)
		}
		return db.Exec(insertStatement+" VALUES "+holders, values...).Error
	}
	batchSize := 5000
	for i := 0; i < len(data); i += 5000 {
		end := i + batchSize
		if end > len(data) {
			end = len(data)
		}
		if err := saveBatch(data[i:end]); err != nil {
			log.Error(err)
			return errors.Wrap(err, "save batch failed")
		}
	}
	return nil
}
