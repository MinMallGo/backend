package global

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path"
	"sync"
)

var Suger *zap.SugaredLogger

var logOnce sync.Once

func LoggerInstance() {
	logOnce.Do(func() {
		sugerInit()
	})
}

func sugerInit() {
	core := zapcore.NewTee(cores()...)
	logger := zap.New(core, zap.AddCallerSkip(1))
	logger.Sync()
	Suger = logger.Sugar()
}

func cores() []zapcore.Core {
	cores := make([]zapcore.Core, 0, 8)
	// file_name,zapcoreLevel
	m := []map[string]zapcore.Level{
		{"mall.log": zapcore.InfoLevel},
		{"err.log": zapcore.ErrorLevel},
		{"panic.log": zapcore.PanicLevel},
	}
	encoder := getEncoder()

	for _, val := range m {
		for filename, level := range val {
			core := zapcore.NewCore(encoder, zapcore.AddSync(createLog(filename)), level)
			cores = append(cores, core)
		}
	}
	return cores
}

func getEncoder() zapcore.Encoder {
	return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
}

func createLog(filename string) *os.File {
	logPath, _ := os.Getwd()
	logDir := path.Join(logPath, "log")
	_, err := os.Stat(logDir)
	if errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir(logDir, os.ModePerm)
		if err != nil {
			panic("mkdir log with error:" + err.Error())
		}
	}

	fmt.Println(path.Join(logDir, filename))
	file, err := os.OpenFile(path.Join(logDir, filename), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic("init log with error: " + err.Error())
	}
	return file
}
