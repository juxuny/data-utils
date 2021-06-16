package weibo

import "encoding/json"

type MetaData struct {
	Url string `json:"url"`
}

func (t MetaData) ToJson() string {
	data, _ := json.Marshal(t)
	return string(data)
}

func (t MetaData) Encode() string {
	return t.ToJson()
}

func DecodeMetaData(data string, out *MetaData) error {
	return json.Unmarshal([]byte(data), out)
}
