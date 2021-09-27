package model

import "testing"

func TestNewUpdateExecutor(t *testing.T) {
	executor := NewUpdateExecutor("uid", "user_purse")
	executor.Where("gold_num = ?", 0)
	executor.Add(908102, "gold_num", 100).
		Add(908103, "gold_num", 200)
	_, err := executor.Exec(DB)
	if err != nil {
		t.Fatal(err)
	}
}
