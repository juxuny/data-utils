package data_utils

import (
	"github.com/juxuny/env/ks"
)

var EnvKey = struct {
}{}

func init() {
	ks.InitKeyName(&EnvKey, true)
}
