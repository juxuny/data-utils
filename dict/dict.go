// 把四六级词库加载到内存 https://github.com/mahavivo/english-wordlists
package dict

import (
	"fmt"
	dataUtils "github.com/juxuny/data-utils"
	"github.com/pkg/errors"
	"io/ioutil"
	"strings"
)

type Word struct {
	Name                  string `json:"name"`                  // 单词
	Description           string `json:"description"`           // 解析
	PhoneticTranscription string `json:"phoneticTranscription"` // 音标
}

func (t *Word) String() string {
	if t.PhoneticTranscription != "" {
		return fmt.Sprintf("%s %s %s", t.Name, t.PhoneticTranscription, t.Description)
	}
	return fmt.Sprintf("%s %s", t.Name, t.Description)
}

type Dict struct {
	Data map[string]Word
}

func NewDict() *Dict {
	return &Dict{
		Data: make(map[string]Word),
	}
}

func (t *Dict) String() string {
	ret := make([]string, 0)
	for _, w := range t.Data {
		ret = append(ret, w.String())
	}
	return strings.Join(ret, "\n")
}

func LoadDict(dataFile string) (*Dict, error) {
	d := NewDict()
	data, err := ioutil.ReadFile(dataFile)
	if err != nil {
		return nil, errors.Wrap(err, "load dict data failed")
	}
	lines := strings.Split(string(data), "\n")
	lines = dataUtils.StringFilter(lines, func(l string) bool {
		return len(l) > 0 && l[0] >= 'a' && l[0] <= 'z'
	})
	for _, line := range lines {
		//log.Debug(line)
		w := parseWord(line)
		//log.Debug(dataUtils.ToJson(w))
		d.Data[w.Name] = w
	}
	return d, nil
}

func parseWord(line string) Word {
	w := Word{}
	if strings.Contains(line, "[") && strings.Contains(line, "]") {
		l := strings.Split(line, "[")
		w.Name = strings.TrimSpace(l[0])
		l = strings.Split(l[1], "]")
		w.Description = strings.Trim(l[1], " \n")
		w.PhoneticTranscription = fmt.Sprintf("[%s]", l[0])
	} else {
		index := strings.Index(line, " ")
		w.Name = strings.TrimSpace(line[:index])
		w.Description = strings.TrimSpace(line[index+1:])
	}
	return w
}
