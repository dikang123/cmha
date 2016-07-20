package log

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/upmio/horus-server/structs"
)

type Level int

const (
	PanicLevel Level = iota
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
)

var (
	LogLever Level = InfoLevel
	logger         = log.New(os.Stderr, "", log.Ltime|log.Ldate)
)

func ConfigLevel(lvl string) error {
	level, err := ParseLevel(lvl)
	if err != nil {
		return err
	}

	LogLever = level
	return nil
}

func (level Level) String() string {
	switch level {
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warning"
	case ErrorLevel:
		return "error"
	case FatalLevel:
		return "fatal"
	case PanicLevel:
		return "panic"
	}

	return "unknown"
}

func ParseLevel(lvl string) (Level, error) {
	switch lvl {
	case "panic":
		return PanicLevel, nil
	case "fatal":
		return FatalLevel, nil
	case "error":
		return ErrorLevel, nil
	case "warn", "warning":
		return WarnLevel, nil
	case "info":
		return InfoLevel, nil
	case "debug":
		return DebugLevel, nil
	}

	var l Level
	return l, fmt.Errorf("not a valid logrus Level: %q", lvl)
}

func LogInit(level string, out io.Writer) error {
	if err := ConfigLevel(level); err != nil {
		return err
	}
	logger = log.New(out, "", log.Lmicroseconds|log.Ldate)
	return nil
}

func Pannic(v ...interface{}) {
	if LogLever >= PanicLevel {
		logger.SetPrefix("[panic] ")
		logger.Println(v...)
	}
}
func Pannicf(format string, v ...interface{}) {
	if LogLever >= PanicLevel {
		logger.SetPrefix("[panic] ")
		logger.Printf(format, v...)
	}
}
func Error(v ...interface{}) {
	if LogLever >= ErrorLevel {
		logger.SetPrefix("[error] ")
		logger.Println(v...)
	}
}

func Errorf(format string, v ...interface{}) {
	if LogLever >= ErrorLevel {
		logger.SetPrefix("[error] ")
		logger.Printf(format, v...)
	}
}

func Info(v ...interface{}) {
	if LogLever >= InfoLevel {
		logger.SetPrefix("[info] ")
		logger.Println(v...)
	}
}

func Infof(format string, v ...interface{}) {
	if LogLever >= InfoLevel {
		logger.SetPrefix("[info] ")
		logger.Printf(format, v...)
	}
}

func Warn(v ...interface{}) {
	if LogLever >= WarnLevel {
		logger.SetPrefix("[warn] ")
		logger.Println(v...)
	}
}

func Warnf(format string, v ...interface{}) {
	if LogLever >= WarnLevel {
		logger.SetPrefix("[warn] ")
		logger.Printf(format, v...)
	}
}

func Debug(v ...interface{}) {
	if LogLever >= DebugLevel {
		logger.SetPrefix("[debug] ")
		logger.Println(v...)
	}
}

func Debugf(format string, v ...interface{}) {
	if LogLever >= DebugLevel {
		logger.SetPrefix("[debug] ")
		logger.Printf(format, v...)
	}
}

func LogInfo(req *http.Request) {
	Info("the req :", req.Method, req.URL.Path)
}

func LogError(err error, req *http.Request) {
	Error("the error :", err.Error(), req.Method, req.URL.Path)
}

func ErrHanlde(w http.ResponseWriter, req *http.Request, err error) {
	bts, _ := json.Marshal(structs.CommonRes{Err: err.Error()})
	http.Error(w, string(bts), http.StatusInternalServerError)
	Errorf("the req %s exec error:%s\n", req.Method+":"+req.URL.Path, err.Error())
}
