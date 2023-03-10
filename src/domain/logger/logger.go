package logger

import (
	"fmt"
	"log"
	"os"
)

var (
	info    *log.Logger
	warning *log.Logger
	err     *log.Logger
	debug   *log.Logger
)

func init() {
	debug = log.New(os.Stdout, "[DEBUG]: ", log.Ldate|log.Ltime|log.Lshortfile)
	info = log.New(os.Stdout, "[INFO]: ", log.Ldate|log.Ltime|log.Lshortfile)
	warning = log.New(os.Stdout, "[WARNING]: ", log.Ldate|log.Ltime|log.Lshortfile)
	err = log.New(os.Stderr, "[ERROR]: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func Debug(message string) {
	if os.Getenv("LOG_LEVEL") != "debug" {
		return
	}

	debug.Output(2, message)
}

func Info(message string) {
	info.Output(2, message)
}

func Warning(message string) {
	warning.Output(2, message)
}

func Error(message string) {
	err.Output(2, message)
}

func Debugf(format string, v ...interface{}) {
	debug.Output(2, fmt.Sprintf(format, v...))
}

func Infof(format string, v ...interface{}) {
	info.Output(2, fmt.Sprintf(format, v...))
}

func Warningf(format string, v ...interface{}) {
	warning.Output(2, fmt.Sprintf(format, v...))
}

func Errorf(format string, v ...interface{}) {
	err.Output(2, fmt.Sprintf(format, v...))
}
