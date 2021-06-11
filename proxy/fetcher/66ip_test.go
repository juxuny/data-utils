package fetcher

import "testing"

func TestNew66Ip(t *testing.T) {
	instance := New66Ip()
	if err := instance.Init(); err != nil {
		t.Fatal(err)
	}
}
