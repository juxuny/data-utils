package cache

import (
	"encoding/base64"
	"fmt"
	"github.com/juxuny/data-utils/log"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path"
	"sync"
)

type fileCache struct {
	dir  string
	data sync.Map
	*sync.Mutex
}

func NewFileCache(dir string) Cache {
	c := &fileCache{
		dir:   dir,
		Mutex: &sync.Mutex{},
	}
	return c
}

func (t *fileCache) touchDir() error {
	stat, err := os.Stat(t.dir)
	if os.IsNotExist(err) {
		_ = os.MkdirAll(t.dir, 0775)
	} else {
		if !stat.IsDir() {
			return fmt.Errorf("%s is a directory", t.dir)
		}
	}
	return nil
}

func (t *fileCache) genKey(key string) string {
	return base64.StdEncoding.EncodeToString([]byte(key))
}

func (t *fileCache) genFileName(key string) string {
	fileName := t.genKey(key) + ".dat"
	return path.Join(t.dir, fileName)
}

func (t *fileCache) Get(key string) (data string, err error) {
	t.Lock()
	defer t.Unlock()
	fileName := t.genFileName(key)
	_, err = os.Stat(fileName)
	if os.IsNotExist(err) {
		return "", ErrNotFound
	}
	byteData, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Error(err)
		return "", errors.Wrap(err, "read data file failed")
	}
	return string(byteData), nil
}

func (t *fileCache) Set(key, value string) (err error) {
	t.Lock()
	defer t.Unlock()
	if err := t.touchDir(); err != nil {
		log.Error(err)
		return errors.Wrap(err, "create storage directory failed")
	}
	fileName := t.genFileName(key)
	err = ioutil.WriteFile(fileName, []byte(value), 0664)
	return err
}
