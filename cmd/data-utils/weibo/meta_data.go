package weibo

import "encoding/json"

type MetaData struct {
	Url string `json:"url"`
}

func (t MetaData) ToJson() string {
	data, _ := json.Marshal(t)
	return string(data)
}
