package data_utils

import (
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
