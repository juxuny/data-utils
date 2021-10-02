package main

import "github.com/juxuny/data-utils/model"

type SearchResult struct {
	List []ResultItem `json:"list"`
}

type ResultItem struct {
	Movie     model.EngMovie         `json:"movie"`
	Subtitle  model.EngSubtitle      `json:"subtitle"`
	Block     model.EngSubtitleBlock `json:"block"`
	VideoFile string                 `json:"video_file"`
}
