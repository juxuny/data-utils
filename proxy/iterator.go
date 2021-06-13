package proxy

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/juxuny/data-utils/log"
	"github.com/juxuny/data-utils/model"
	"github.com/juxuny/data-utils/proxy/dt"
	"github.com/pkg/errors"
)

type Iterator interface {
	Len() int
	Next() (address string, err error)
	Reset()
	Init(address ...string)
}

type memoryIterator struct {
	data   []string
	start  int
	length int
	index  int
}

func NewMemoryProxyAddressIterator() *memoryIterator {
	ret := &memoryIterator{}
	return ret
}

func (t *memoryIterator) Init(address ...string) {
	if len(t.data) < len(address) {
		t.data = make([]string, len(address))
	}
	t.start = 0
	t.length = len(address)
	t.index = t.start
	copy(t.data, address)
}

func (t *memoryIterator) Next() (address string, err error) {
	if t.index >= len(t.data) {
		return "", dt.ErrEOF
	}
	address = t.data[t.index]
	t.index += 1
	return
}

func (t *memoryIterator) Reset() {
	t.start = 0
}

func (t *memoryIterator) Len() int {
	return t.length
}

type IteratorOption struct {
	IgnoreConnectFailed bool
}

type databaseProxyAddressIterator struct {
	*memoryIterator
	db        *model.DB
	lastId    int64
	batchSize int64
	opt       IteratorOption
}

func NewDatabaseProxyAddressIterator(db *model.DB, opt ...IteratorOption) (ret *databaseProxyAddressIterator) {
	ret = &databaseProxyAddressIterator{
		batchSize:      1000, // 默认读多少个数据
		db:             db,
		memoryIterator: NewMemoryProxyAddressIterator(),
	}
	if len(opt) > 0 {
		ret.opt = opt[0]
	}
	return
}

func (t *databaseProxyAddressIterator) getDb() *gorm.DB {
	db := t.db.Where("id > ?", t.index).Order("id ASC")
	if t.opt.IgnoreConnectFailed {
		db = db.Where("latency > 0 || latency IS NULL")
	}
	return db
}

func (t *databaseProxyAddressIterator) Len() int {
	var total int
	if err := t.getDb().Model(&model.Proxy{}).Count(&total).Error; err != nil {
		log.Error(err)
		return 0
	}
	return total
}

func (t *databaseProxyAddressIterator) Reset() {
	t.memoryIterator.Reset()
	t.index = 0
}

func (t *databaseProxyAddressIterator) Init(address ...string) {
	panic("no implement Init")
}

func (t *databaseProxyAddressIterator) loadFromDatabase() (err error) {
	var data []model.Proxy
	db := t.getDb()
	if err := db.Find(&data).Error; err != nil {
		if !model.IsErrNoDataInDb(err) {
			log.Error(err)
			return errors.Wrap(err, "load data from database failed")
		}
	}
	if len(data) == 0 {
		return dt.ErrEOF
	}
	var addressList []string
	for _, item := range data {
		if item.Id > t.lastId {
			t.lastId = item.Id
		}
		proxyItem := dt.ServerItem{
			Schema:   "",
			Ip:       item.Ip,
			Port:     item.Port,
			Provider: item.ProviderName,
		}
		if item.Socks5 {
			proxyItem.Schema = dt.SchemaTypeSocks5
		} else if item.Http || item.Https {
			proxyItem.Schema = dt.SchemaTypeHttp
		}
		addressList = append(addressList, fmt.Sprintf("%s://%s:%d", proxyItem.Schema, proxyItem.Ip, proxyItem.Port))
	}
	t.memoryIterator.Init(addressList...)
	return nil
}

func (t *databaseProxyAddressIterator) Next() (address string, err error) {
	address, err = t.memoryIterator.Next()
	if err == nil {
		return address, nil
	} else if err == dt.ErrEOF {
		if err := t.loadFromDatabase(); err != nil {
			if err == dt.ErrEOF {
				return "", dt.ErrEOF
			}
		}
		return t.memoryIterator.Next()
	} else {
		log.Error(err)
		return "", err
	}
}
