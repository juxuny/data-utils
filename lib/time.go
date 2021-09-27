package lib

import "time"

type _Time struct{}

var Time = _Time{}

func (_Time) NowPointer() *time.Time {
	t := time.Now()
	return &t
}
