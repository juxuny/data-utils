package model

import (
	"fmt"
	"github.com/juxuny/data-utils/global_key"
	"github.com/juxuny/data-utils/log"
	"github.com/juxuny/env"
	"github.com/pkg/errors"
	"reflect"
	"runtime/debug"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var logger = log.NewLogger("model")

func init() {
	defaultConfig, _ = GetEnvConfig()
}

type Config struct {
	DbHost     string `json:"db_host"`
	DbPort     int    `json:"db_port"`
	DbUser     string `json:"db_user"`
	DbPassword string `json:"db_password"`
	DbName     string `json:"db_name"`
	DbDebug    bool   `json:"db_debug"`
}

func GetEnvConfig() (config Config, err error) {
	config = Config{
		DbHost:     env.GetString(global_key.EnvKey.DbHost, "127.0.0.1"),
		DbPort:     env.GetInt(global_key.EnvKey.DbPort, 3306),
		DbUser:     env.GetString(global_key.EnvKey.DbUser, "root"),
		DbPassword: env.GetString(global_key.EnvKey.DbPwd, ""),
		DbName:     env.GetString(global_key.EnvKey.DbName),
		DbDebug:    env.GetBool(global_key.EnvKey.DbDebug, false),
	}
	return config, nil
}

var defaultConfig Config

func Open(config ...Config) (db *DB, err error) {
	finalConfig := defaultConfig
	log.Debug(finalConfig)
	if len(config) > 0 {
		finalConfig = config[0]
	}
	d, err := gorm.Open("mysql", fmt.Sprintf(
		"%s:%s@(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		finalConfig.DbUser,
		finalConfig.DbPassword,
		finalConfig.DbHost,
		finalConfig.DbPort,
		finalConfig.DbName,
	))
	if err != nil {
		logger.Error(err)
		return nil, errors.Wrap(err, "failed to connect to database")
	}
	d.SingularTable(true) // 表名都用奇数
	d.LogMode(finalConfig.DbDebug)
	return &DB{
		DB: d,
	}, nil
}

// 创建数据库事务，代码崩溃自动回滚
func (t *DB) Begin(cb func(*gorm.DB) error) error {
	tx := t.DB.Begin()
	if err := tx.Error; err != nil {
		return err
	}

	// 下面的代码发生错误时，回滚事务
	defer func() {
		if err := recover(); err != nil {
			logger.Error("Transaction_Error: %v", err)
			debug.PrintStack()
			if err := tx.Rollback().Error; err != nil {
				logger.Errorf("Transaction_Rollback_Error: %v", err)
			}
			panic(err)
		}
	}()

	if err := cb(tx); err != nil {
		logger.Print(err)
		tx.Rollback()
		return err
	} else if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

// 执行原生SQL并存储到[]map中
func (t *DB) RawToMap(sql string, data *[]map[string]interface{}, args ...interface{}) (columns []string, err error) {

	rows, err := t.DB.Raw(sql, args...).Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	makeResultReceiver := func(length int) []interface{} {
		result := make([]interface{}, 0, length)
		for i := 0; i < length; i++ {
			var current interface{}
			current = struct{}{}
			result = append(result, &current)
		}
		return result
	}

	columns, _ = rows.Columns()
	length := len(columns)
	for rows.Next() {
		current := makeResultReceiver(length)
		if err := rows.Scan(current...); err != nil {
			panic(err)
		}
		value := make(map[string]interface{})
		for i := 0; i < length; i++ {
			key := columns[i]
			val := *(current[i]).(*interface{})
			if val == nil {
				value[key] = nil
				continue
			}
			vType := reflect.TypeOf(val)
			switch vType.String() {
			case "int64":
				value[key] = val.(int64)
			case "string":
				value[key] = val.(string)
			case "time.Time":
				value[key] = val.(time.Time)
			case "[]uint8":
				value[key] = string(val.([]uint8))
			default:
				// TODO remember add other data type
				fmt.Printf("unsupport data type '%s' now\n", vType)
				value[key] = val
			}
		}
		*data = append(*data, value)
	}

	return columns, nil
}

var IsErrNoDataInDb = gorm.IsRecordNotFoundError

type DB struct {
	*gorm.DB
}
