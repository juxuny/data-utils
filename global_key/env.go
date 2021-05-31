package global_key

import "github.com/juxuny/env/ks"

var EnvKey = struct {
	DbHost  string
	DbPwd   string
	DbUser  string
	DbPort  string
	DbName  string
	DbDebug string
}{}

func init() {
	ks.InitKeyName(&EnvKey, true)
}
