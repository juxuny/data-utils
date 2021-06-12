package fetcher

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/djimenez/iconv-go"
	"github.com/juxuny/data-utils/curl"
	"github.com/juxuny/data-utils/log"
	"github.com/juxuny/data-utils/model"
	"github.com/juxuny/data-utils/proxy/dt"
	"github.com/pkg/errors"
	"net"
	"strconv"
)

type _66ip struct {
	TotalPage   int
	data        dt.ServerList
	saveHandler dt.SaveHandler
}

func New66Ip(saveHandler dt.SaveHandler) (ret *_66ip) {
	return &_66ip{
		saveHandler: saveHandler,
	}
}

func (t *_66ip) Len() int {
	return len(t.data)
}

func (t *_66ip) Page(page, pageSize int) (ret dt.ServerList, err error) {
	offset := (page - 1) * pageSize
	end := offset + pageSize
	if len(t.data) < end {
		end = len(t.data)
	}
	return t.data[offset:end], nil
}

func (t *_66ip) AllData() (ret dt.ServerList, err error) {
	return t.data, nil
}

func (t *_66ip) getUrl(page int) string {
	return fmt.Sprintf("http://www.66ip.cn/%d.html", page)
}

func (t *_66ip) parsePage(page int, cb ...func(selection *goquery.Document) error) (ret dt.ServerList, err error) {
	content, err := curl.Get(t.getUrl(page))
	if err != nil {
		return nil, err
	}
	converter, err := iconv.NewConverter("gb2312", "utf-8")
	if err != nil {
		return nil, errors.Wrap(err, "create converter failed")
	}
	var out string
	out, err = converter.ConvertString(string(content))
	doc, err := goquery.NewDocumentFromReader(bytes.NewBufferString(out))
	if err != nil {
		return nil, err
	}
	doc.Find(
		"div.layui-row.layui-col-space15 table tbody tr",
	).Each(func(i int, selection *goquery.Selection) {
		ipValue := selection.Find("td:nth-child(1)").Text()
		portValue := selection.Find("td:nth-child(2)").Text()
		ip := net.ParseIP(ipValue)
		port, errPort := strconv.ParseInt(portValue, 10, 64)
		if errPort == nil && ip != nil {
			log.Debug("ip: ", ipValue, " port: ", portValue)
			ret = append(ret, dt.ServerItem{
				Schema:   dt.SchemaTypeHttp,
				Ip:       ipValue,
				Port:     int(port),
				Provider: model.Provider66Ip,
			})
		}
	})
	if len(cb) > 0 {
		for _, f := range cb {
			if err := f(doc); err != nil {
				return nil, errors.Wrap(err, "callback error")
			}
		}
	}
	return ret, nil
}

func (t *_66ip) Init() error {
	var getTotalPageNum = func(doc *goquery.Document) error {
		doc.Find("#PageList a").Each(func(i int, selection *goquery.Selection) {
			pageStr := selection.Text()
			if p, err := strconv.ParseInt(pageStr, 10, 64); err == nil {
				if t.TotalPage < int(p) {
					t.TotalPage = int(p)
				}
			}
		})
		return nil
	}
	data, err := t.parsePage(1, getTotalPageNum)
	if err != nil {
		return errors.Wrap(err, "init data failed")
	}
	for _, item := range data {
		t.data = append(t.data, item)
	}
	if t.saveHandler != nil {
		if err := t.saveHandler.SaveServerList(data); err != nil {
			return errors.Wrapf(err, "save server list failed, page=%v", 1)
		}
	}
	for i := 2; i <= t.TotalPage; i++ {
		log.Info("load page: ", i)
		data, err := t.parsePage(i, getTotalPageNum)
		if err != nil {
			return errors.Wrap(err, "get page data failed")
		}
		for _, item := range data {
			t.data = append(t.data, item)
		}
		if t.saveHandler != nil {
			if err := t.saveHandler.SaveServerList(data); err != nil {
				return errors.Wrapf(err, "save server list failed, page=%v", i)
			}
		}
	}
	return nil
}

func (t *_66ip) Reset() {
	t.TotalPage = 0
	t.data = make(dt.ServerList, 0)
}
