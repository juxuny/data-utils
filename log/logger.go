package log

import (
	"fmt"
	"log"
	"os"
)

type Level int

const (
	LevelDebug = 1
	LevelInfo  = 2
	LevelWarn  = 3
	LevelError = 4
	LevelNone  = 10000 // 不输出
)

var (
	currentLevel = Level(LevelDebug)
)

func SetLevel(lv Level) {
	currentLevel = lv
}

var (
	enableSql = true
)

func DisableSQL() {
	enableSql = false
}

type ILogger interface {
	SetPrefix(string)
	Println(...interface{})
	Error(...interface{})
	Errorf(format string, v ...interface{})
	Warn(...interface{})
	Warnf(format string, v ...interface{})
	Info(...interface{})
	Infof(format string, v ...interface{})
	Debug(...interface{})
	Debugf(format string, v ...interface{})
	Printf(string, ...interface{})
	Print(...interface{})
	SQL(...interface{})
	Consuming(...interface{})
	Output(...interface{})
	Flush()
}

var defaultLogger = NewLogger("DEFAULT", 4)

type Logger struct {
	l         *log.Logger
	CallDepth int
}

func (t *Logger) Flush() {
}

func (t *Logger) SetPrefix(s string) {
	t.l.SetPrefix(s)
}

func (t *Logger) Println(v ...interface{}) {
	t.Output("[INFO] " + fmt.Sprint(v...))
}

func (t *Logger) Error(v ...interface{}) {
	if currentLevel <= LevelError {
		t.Output("[ERROR] \033[0;31m" + fmt.Sprint(v...) + "\033[0m")
	}
}

func (t *Logger) Errorf(format string, v ...interface{}) {
	if currentLevel <= LevelError {
		t.Output("[ERROR] \033[0;31m" + fmt.Sprintf(format, v...) + "\033[0m")
	}
}

func (t *Logger) Warn(v ...interface{}) {
	if currentLevel <= LevelWarn {
		t.Output("[WARN] \033[0;33m" + fmt.Sprint(v...) + "\033[0m")
	}
}

func (t *Logger) Warnf(format string, v ...interface{}) {
	if currentLevel <= LevelWarn {
		t.Output("[WARN] \033[0;33m" + fmt.Sprintf(format, v...) + "\033[0m")
	}
}

func (t *Logger) Info(v ...interface{}) {
	if currentLevel <= LevelInfo {
		t.Output("[INFO] \033[0;32m" + fmt.Sprint(v...) + "\033[0m")
	}
}

func (t *Logger) Infof(format string, v ...interface{}) {
	if currentLevel <= LevelInfo {
		t.Output("[INFO] \033[0;32m" + fmt.Sprintf(format, v...) + "\033[0m")
	}
}

func (t *Logger) Debug(v ...interface{}) {
	if currentLevel <= LevelDebug {
		t.Output("[DEBUG] " + fmt.Sprint(v...) + "")
	}
}

func (t *Logger) Debugf(format string, v ...interface{}) {
	if currentLevel <= LevelDebug {
		t.Output("[DEBUG] " + fmt.Sprint(v...) + "")
	}
}

func (t *Logger) Printf(format string, v ...interface{}) {
	if currentLevel <= LevelInfo {
		t.Output("[INFO] " + fmt.Sprintf(format, v...))
	}
}

func (t *Logger) SQL(v ...interface{}) {
	if enableSql {
		t.Output("[SQL] \033[0;35m" + fmt.Sprint(v...) + "\033[0m")
	}
}

func (t *Logger) Consuming(v ...interface{}) {
	if enableSql {
		t.Output("[CONSUMING] \033[0;36m" + fmt.Sprint(v...) + "\033[0m")
	}
}

func NewLogger(prefix string, callDepth ...int) ILogger {
	cd := 3
	if len(callDepth) > 0 {
		cd = callDepth[0]
	}
	ret := &Logger{
		l:         log.New(os.Stdout, prefix+" ", log.Ltime|log.Llongfile|log.Ldate|log.LstdFlags),
		CallDepth: cd,
	}
	return ret
}

func (t *Logger) Print(v ...interface{}) {
	t.Output("[INFO] " + fmt.Sprint(v...))
}

func (t *Logger) Output(v ...interface{}) {
	_ = t.l.Output(t.CallDepth, fmt.Sprint(v...))
}

func Debug(v ...interface{}) {
	defaultLogger.Debug(v...)
}

func Debugf(format string, v ...interface{}) {
	defaultLogger.Debugf(format, v...)
}

func Info(v ...interface{}) {
	defaultLogger.Info(v...)
}

func Infof(format string, v ...interface{}) {
	defaultLogger.Infof(format, v...)
}

func Warn(v ...interface{}) {
	defaultLogger.Warn(v...)
}

func Warnf(format string, v ...interface{}) {
	defaultLogger.Warnf(format, v...)
}

func Error(v ...interface{}) {
	defaultLogger.Error(v...)
}

func Errorf(format string, v ...interface{}) {
	defaultLogger.Errorf(format, v...)
}

func Fatal(v ...interface{}) {
	defaultLogger.Error(v...)
	os.Exit(-1)
}

func Fatalf(format string, v ...interface{}) {
	defaultLogger.Errorf(format, v...)
	os.Exit(-1)
}
