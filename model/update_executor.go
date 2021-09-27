package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/juxuny/data-utils/log"
	"github.com/pkg/errors"
	"strings"
)

type UpdateExecutor interface {
	Where(where string, args ...interface{}) UpdateExecutor
	Add(key interface{}, fieldName string, value interface{}) UpdateExecutor
	Exec(db *gorm.DB) (affected int64, err error)
	GetSql() (string, []interface{}, error)
	DataLen() int
}

type updateExecutorImpl struct {
	actionSet   map[string]map[interface{}]interface{}
	where       string
	whereValues []interface{}
	keyColumn   string
	tableName   string
}

func NewUpdateExecutor(keyColumn string, tableName string) UpdateExecutor {
	executor := &updateExecutorImpl{
		actionSet: make(map[string]map[interface{}]interface{}),
		keyColumn: keyColumn,
		tableName: tableName,
	}
	return executor
}

func (t *updateExecutorImpl) DataLen() int {
	return len(t.actionSet)
}

func (t *updateExecutorImpl) Where(where string, args ...interface{}) UpdateExecutor {
	t.where = where
	t.whereValues = args
	return t
}

func (t *updateExecutorImpl) GetSql() (string, []interface{}, error) {
	if len(t.actionSet) == 0 {
		return "", nil, errors.Errorf("data set is empty")
	}
	keyMap := make(map[interface{}]bool)
	update := make([]string, 0)
	for field, action := range t.actionSet {
		if len(action) == 0 {
			return "", nil, errors.Errorf("field '%s' has no update action", field)
		}
		when := make([]string, 0)
		for key, value := range action {
			keyMap[key] = true
			when = append(when, fmt.Sprintf("WHEN %d THEN %v", key, value))
		}
		statement := fmt.Sprintf("`%s` = `%s` + (CASE %s %s)", field, field, t.keyColumn, strings.Join(when, " ")+" ELSE 0 END")
		update = append(update, statement)
	}
	if len(keyMap) == 0 {
		return "", nil, errors.Errorf("key list is empty")
	}
	where := ""
	for key := range keyMap {
		if where != "" {
			where += ", "
		}
		where += fmt.Sprintf("%v", key)
	}
	where = fmt.Sprintf("WHERE %s IN (%s)", t.keyColumn, where)
	if t.where != "" {
		where += " AND " + fmt.Sprintf("(%s)", t.where)
	}
	ret := fmt.Sprintf("UPDATE `%s` SET %s %s", t.tableName, strings.Join(update, ", "), where)
	return ret, t.whereValues, nil
}

func (t *updateExecutorImpl) addInterfaceValue(raw interface{}, add interface{}) (ret interface{}) {
	if raw == nil {
		return add
	}
	switch v := raw.(type) {
	case int:
		addValue, _ := add.(int)
		return v + addValue
	case int8:
		addValue, _ := add.(int8)
		return v + addValue
	case int16:
		addValue, _ := add.(int16)
		return v + addValue
	case int32:
		addValue, _ := add.(int32)
		return v + addValue
	case int64:
		addValue, _ := add.(int64)
		return v + addValue
	case uint8:
		addValue, _ := add.(uint8)
		return v + addValue
	case uint16:
		addValue, _ := add.(uint16)
		return v + addValue
	case uint32:
		addValue, _ := add.(uint32)
		return v + addValue
	case uint64:
		addValue, _ := add.(uint64)
		return v + addValue
	case float64:
		addValue, _ := add.(float64)
		return v + addValue
	case float32:
		addValue, _ := add.(float32)
		return v + addValue
	}
	return raw
}

func (t *updateExecutorImpl) Add(key interface{}, fieldName string, value interface{}) UpdateExecutor {
	if set, b := t.actionSet[fieldName]; b {
		set[key] = t.addInterfaceValue(set[key], value)
		t.actionSet[fieldName] = set
	} else {
		set = map[interface{}]interface{}{
			key: value,
		}
		t.actionSet[fieldName] = set
	}
	return t
}

func (t *updateExecutorImpl) Exec(db *gorm.DB) (affected int64, err error) {
	statement, values, err := t.GetSql()
	if err != nil {
		log.Error(err)
		return 0, errors.Wrap(err, "generate sql statement failed")
	}
	scope := db.Exec(statement, values...)
	affected = scope.RowsAffected
	err = scope.Error
	return
}
