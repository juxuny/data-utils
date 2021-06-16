package weibo

import (
	"fmt"
	data_utils "github.com/juxuny/data-utils"
	"github.com/juxuny/data-utils/curl"
	"github.com/juxuny/data-utils/log"
	"github.com/juxuny/data-utils/model"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
	"strconv"
)

type friendshipsParser struct {
	*cookieManager
	db *model.DB
}

func NewFriendshipParser() *friendshipsParser {
	return &friendshipsParser{
		cookieManager: NewCookieManager(),
	}
}

func (t *friendshipsParser) CheckValid(metaData MetaData) (isOk bool) {
	parsedUrl, err := url.Parse(metaData.Url)
	if err != nil {
		log.Debug(err)
		return
	}
	if parsedUrl.Path == "/ajax/friendships/friends" {
		return true
	}
	return false
}

func (t *friendshipsParser) Prepare(db *model.DB) error {
	t.db = db
	return nil
}

func (t *friendshipsParser) saveFriendships(currentUid int64, fans FansVo) (err error) {
	fans.WeiboFans.Uid = currentUid
	fans.WeiboFans.FansId = fans.Id
	// create user
	var userInfo model.WeiboUser
	if err := t.db.First(&userInfo, "id = ?", fans.Id).Error; err != nil {
		if !model.IsErrNoDataInDb(err) {
			log.Error(err)
			return errors.Wrap(err, "check user info failed")
		} else {
			if err := t.db.Create(fans.WeiboUser).Error; err != nil {
				log.Error(err)
			}
		}
	} else {
		if err := t.db.Save(&fans.WeiboUser).Error; err != nil {
			log.Fatal(err)
		}
	}
	// create friendship
	var weiboFans model.WeiboFans
	if err := t.db.First(&weiboFans, "uid = ? AND fans_id = ?", currentUid, fans.Id).Error; err != nil {
		if !model.IsErrNoDataInDb(err) {
			log.Error(err)
			return errors.Wrap(err, "check fans relation failed")
		} else {
			if err := t.db.Create(fans.WeiboFans).Error; err != nil {
				log.Error(err)
			}
		}
	} else {
		if err := t.db.Table(model.WeiboFans{}.TableName()).Where("uid = ? AND fans_id = ?", currentUid, fans.Id).Updates(&fans.WeiboFans).Error; err != nil {
			log.Error(err)
		}
	}

	return nil
}

func (t *friendshipsParser) Parse(metaData MetaData) (jobList data_utils.JobList, err error) {
	parsedUrl, err := url.Parse(metaData.Url)
	if err != nil {
		log.Error(err)
		return
	}
	var currentUid int64
	currentUid, err = strconv.ParseInt(parsedUrl.Query().Get("uid"), 10, 64)
	if err != nil {
		log.Error(err)
		return nil, errors.Wrapf(err, "invalid uid, err:%v", err)
	}

	var header = http.Header{}
	cookie := t.cookieManager.GetCookies("")
	header.Add("Cookie", cookie)
	var resp FriendshipsResp
	if err := curl.GetJson(metaData.Url, &resp, header); err != nil {
		log.Error(err)
		return nil, errors.Wrap(err, "request friendship failed")
	}
	log.Debug(data_utils.ToJson(resp))
	if resp.Ok != 1 {
		log.Warn("resp.ok == ", resp.Ok)
		return nil, errors.Errorf("ok = %v", resp.Ok)
	}

	// save data
	for _, u := range resp.Users {
		if err := t.saveFriendships(currentUid, u); err != nil {
			log.Warn(err)
			continue
		}
	}

	values := parsedUrl.Query()
	page, err := strconv.ParseInt(values.Get("page"), 10, 64)
	if err != nil {
		log.Error(err)
		page = 1
	} else {
		page += 1
	}
	if resp.NextCursor > 0 {
		values.Set("page", fmt.Sprintf("%d", page))
		newMetaData := MetaData{
			Url: fmt.Sprintf("%s?%s", "https://weibo.com/ajax/friendships/friends", values.Encode()),
		}
		jobList = append(jobList, data_utils.Job{
			JobType:  model.JobTypeWeibo,
			MetaData: newMetaData.Encode(),
		})
	}
	for _, u := range resp.Users {
		log.Debug(u.Id)
		parsedUrl.Query().Set("uid", u.IdStr)
		parsedUrl.Query().Set("page", "1")
		newMetaData := MetaData{Url: parsedUrl.String()}
		jobList = append(jobList, data_utils.Job{
			JobType:  model.JobTypeWeibo,
			MetaData: newMetaData.Encode(),
		})
	}
	return
}
