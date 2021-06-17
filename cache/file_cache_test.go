package cache

import "testing"

func TestNewFileCache(t *testing.T) {
	c := NewFileCache("tmp")
	if err := c.Set("cookie", "you"); err != nil {
		t.Fatal(err)
	}
	v, err := c.Get("cookie")
	if err != nil {
		t.Fatal(err)
	}
	if v != "you" {
		t.Fatalf("the real value is 'you', but got '%s'", v)
	}
}
