package ffmpeg

import "github.com/juxuny/env/ks"

var TagKey = struct {
	Language string
	Title    string
}{}

func init() {
	ks.InitKeyName(&TagKey, false)
}
