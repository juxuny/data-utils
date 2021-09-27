package model

import (
	"fmt"
	"github.com/fatih/camelcase"
	"github.com/jinzhu/gorm"
	"github.com/juxuny/data-utils/log"
	"github.com/pkg/errors"
	"reflect"
	"strings"
)

type InsertExecutor interface {
	Model(model interface{}) InsertExecutor
	Add(data ...interface{}) InsertExecutor
	GetSql() (statement string, values []interface{}, err error)
	Exec(tx *gorm.DB) (affected int64, err error)
	// 唯一键值冲突时更新
	OnDuplicateUpdate(map[string]struct{}) InsertExecutor
	// 唯一键值冲突时累加
	OnDuplicateIncrease(map[string]struct{}) InsertExecutor
	DataLen() int
}

type InsertExecutorImpl struct {
	insertIgnore        bool
	model               reflect.Value
	data                []map[string]interface{}
	fieldMap            map[string]struct{}
	onDuplicateUpdate   map[string]struct{} // 唯一索引冲突直接更新的字段
	onDuplicateIncrease map[string]struct{} // 唯一索引冲突时累加的字段
	tableName           string
}

func NewInsertExecutor(insertIgnore bool) InsertExecutor {
	executor := &InsertExecutorImpl{
		fieldMap:     make(map[string]struct{}),
		insertIgnore: insertIgnore,
	}
	return executor
}

func (t *InsertExecutorImpl) DataLen() int {
	return len(t.data)
}

func (t *InsertExecutorImpl) OnDuplicateUpdate(update map[string]struct{}) InsertExecutor {
	t.onDuplicateUpdate = update
	return t
}

func (t *InsertExecutorImpl) OnDuplicateIncrease(update map[string]struct{}) InsertExecutor {
	t.onDuplicateIncrease = update
	return t
}

func (t *InsertExecutorImpl) Model(model interface{}) InsertExecutor {
	t.model = reflect.ValueOf(model)
	tt := t.model.Type()
	for i := 0; i < tt.NumField(); i++ {
		f := tt.Field(i)
		column, ignore := getFieldNameFromTag(f)
		if !ignore {
			t.fieldMap[column] = struct{}{}
		}
	}
	//get table name
	l := camelcase.Split(tt.Name())
	for i := 0; i < len(l); i++ {
		l[i] = strings.ToLower(l[i])
	}
	t.tableName = strings.Join(l, "_")
	return t
}

func (t *InsertExecutorImpl) Add(data ...interface{}) InsertExecutor {
	for _, item := range data {
		m := convertToMap(item)
		t.data = append(t.data, m)
	}
	return t
}

func (t *InsertExecutorImpl) GetSql() (statement string, values []interface{}, err error) {
	if len(t.data) == 0 {
		return "", nil, errors.Errorf("data is empty")
	}
	if t.tableName == "" {
		return "", nil, errors.Errorf("tableName is not set, call Model() first")
	}
	usedField := make(map[string]bool)
	for _, item := range t.data {
		for k := range item {
			usedField[k] = true
		}
	}
	fieldList := make([]string, 0)
	for f := range usedField {
		fieldList = append(fieldList, f)
	}
	if t.insertIgnore {
		statement = "INSERT IGNORE INTO "
	} else {
		statement = "INSERT INTO "
	}
	statement += fmt.Sprintf("`%s` (%s) VALUES ", t.tableName, wrapAndJoin(fieldList, "`", ", "))
	var holders []string
	for _, m := range t.data {
		holder := strings.Trim(strings.Repeat("?, ", len(fieldList)), ", ")
		holders = append(holders, "("+holder+")")
		for _, f := range fieldList {
			if v, found := m[f]; found {
				values = append(values, v)
			} else {
				values = append(values, nil)
			}
		}
	}
	statement += strings.Join(holders, ", ")
	if !t.insertIgnore && len(t.onDuplicateUpdate) > 0 {
		update := ""
		for k := range t.onDuplicateUpdate {
			if update != "" {
				update += ", "
			}
			update += fmt.Sprintf("%s=VALUES(%s)", k, k)
		}
		statement += " ON DUPLICATE KEY UPDATE " + update
	} else if !t.insertIgnore && len(t.onDuplicateIncrease) > 0 {
		update := ""
		for k := range t.onDuplicateIncrease {
			if update != "" {
				update += ", "
			}
			update += fmt.Sprintf("%s = %s + VALUES(%s)", k, k, k)
		}
		statement += " ON DUPLICATE KEY UPDATE " + update
	}
	return
}

func (t *InsertExecutorImpl) Exec(tx *gorm.DB) (affected int64, err error) {
	statement, values, err := t.GetSql()
	if err != nil {
		log.Error(err)
		return 0, errors.Wrap(err, "get sql statement failed")
	}
	result := tx.Exec(statement, values...)
	if err := result.Error; err != nil {
		log.Error(err)
		return 0, errors.Wrap(err, "insert batch data failed")
	}
	return result.RowsAffected, nil
}

func getFieldNameFromTag(t reflect.StructField) (field string, ignore bool) {
	gormTag := t.Tag.Get("gorm")
	if gormTag == "-" {
		return "", true
	}
	if t.Type.Kind() == reflect.Struct || t.Type.Kind() == reflect.Slice {
		return "", true
	}
	l := strings.Split(gormTag, ";")
	column := ""
	for _, item := range l {
		kv := strings.Split(item, ":")
		if len(kv) == 2 && strings.ToLower(strings.Trim(kv[0], " ")) == "column" {
			column = strings.Trim(kv[1], " ")
			break
		}
	}
	if column != "" {
		return column, false
	}

	jsonTag := t.Tag.Get("json")
	if jsonTag == "" {
		jsonTag = t.Name
	}
	l = camelcase.Split(jsonTag)
	for i := 0; i < len(l); i++ {
		l[i] = strings.ToLower(l[i])
	}
	return strings.Join(l, "_"), false
}

// gorm tag里是否有default
func hasDefault(t reflect.StructField) bool {
	gormTag := t.Tag.Get("gorm")
	kvs := strings.Split(gormTag, ";")
	for _, item := range kvs {
		kv := strings.Split(item, ":")
		if len(kv) > 0 && strings.ToLower(strings.Trim(kv[0], " ")) == "default" {
			return true
		}
	}
	return false
}

func convertToMap(data interface{}) (ret map[string]interface{}) {
	vv := reflect.ValueOf(data)
	tt := reflect.TypeOf(data)
	ret = make(map[string]interface{})
	if vv.Kind() == reflect.Struct {
		for i := 0; i < tt.NumField(); i++ {
			f := tt.Field(i)
			column, ignore := getFieldNameFromTag(f)
			if ignore {
				continue
			}
			v := vv.Field(i)
			if v.IsZero() && hasDefault(f) {
				continue
			}
			ret[column] = v.Interface()
		}
	} else if vv.Kind() == reflect.Map {
		for _, k := range vv.MapKeys() {
			v := vv.MapIndex(k)
			ret[fmt.Sprintf("%v", k.Interface())] = v.Interface()
		}
	} else if vv.Kind() == reflect.Ptr {
		return convertToMap(vv.Elem().Interface())
	}
	return ret
}

func wrapAndJoin(list []string, bracket string, sep string) string {
	ret := ""
	for _, item := range list {
		if ret != "" {
			ret += sep
		}
		ret += fmt.Sprintf("%s%s%s", bracket, item, bracket)
	}
	return ret
}
