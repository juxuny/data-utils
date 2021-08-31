package data_utils

import (
	"encoding/json"
	"github.com/juxuny/data-utils/log"
	"io/ioutil"
	"runtime/debug"
	"strings"
)

func GetListFromFile(fileName string) (lines []string, err error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	l := strings.Split(string(data), "\n")
	lines = make([]string, 0)
	for _, line := range l {
		tmp := strings.Trim(line, "\r\t ")
		if tmp != "" {
			lines = append(lines, tmp)
		}
	}
	return lines, nil
}

func RecoverRun(f func()) {
	if err := recover(); err != nil {
		log.Error(err)
		debug.PrintStack()
		return
	}
	f()
}

func ToJson(v interface{}) string {
	data, _ := json.Marshal(v)
	return string(data)
}

func StringFilter(l []string, filter ...func(l string) bool) []string {
	var isOk = func(data string, filters []func(l string) bool) bool {
		if len(filters) > 0 {
			for _, f := range filters {
				if !f(data) {
					return false
				}
			}
		}
		return true
	}
	var ret = make([]string, 0)
	for _, item := range l {
		if isOk(item, filter) {
			ret = append(ret, item)
		}
	}
	return ret
}
