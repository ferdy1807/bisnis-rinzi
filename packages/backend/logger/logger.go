package logger

import (
	"fmt"
	"log"
	"os"
)

type Logger struct {
	infoLog  *log.Logger
	errorLog *log.Logger
}

var Log *Logger

func InitLogger() {
	Log = &Logger{
		infoLog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.Lshortfile),
		errorLog: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func Info(format string, v ...interface{}) {
	if Log == nil {
		InitLogger()
	}
	Log.infoLog.Output(2, fmt.Sprintf(format, v...))
}

func Error(format string, v ...interface{}) {
	if Log == nil {
		InitLogger()
	}
	Log.errorLog.Output(2, fmt.Sprintf(format, v...))
}
