package model

import (
	"fmt"
	"papa/lib"
	"testing"
)

func TestNewInsertExecutor(t *testing.T) {
	executor := NewInsertExecutor(true)
	executor.Model(AdminUser{})
	for i := 0; i < 10; i++ {
		executor.Add(AdminUser{
			Id:       lib.ID(135 + i),
			Username: fmt.Sprintf("user%d", i),
			Status:   true,
		})
	}
	for i := 0; i < 10; i++ {
		executor.Add(map[string]interface{}{
			"id":       145 + i,
			"username": fmt.Sprintf("user%d", i),
		})
	}
	executor.OnDuplicateUpdate(map[string]struct{}{
		"createTime": {},
		"lastLogin":  {},
	})
	//statement, values, err := executor.GetSql()
	//if err != nil {
	//	t.Fatal(err)
	//}
	//t.Log(statement)
	//t.Log(values)
	_, err := executor.Exec(DB)
	if err != nil {
		t.Fatal(err)
	}
}
