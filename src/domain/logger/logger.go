package logger

import (
	"fmt"
	"log"
	"os"
)

var (
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
	Debug   *log.Logger
)

func init() {
	Info = log.New(os.Stdout, "[INFO]: ", log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(os.Stdout, "[WARNING]: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(os.Stderr, "[ERROR]: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func Debugf(format string, v ...interface{}) {
	Debug.Output(2, fmt.Sprintf(format, v...))
}

func Infof(format string, v ...interface{}) {
	Info.Output(2, fmt.Sprintf(format, v...))
}

func Warningf(format string, v ...interface{}) {
	Warning.Output(2, fmt.Sprintf(format, v...))
}

func Errorf(format string, v ...interface{}) {
	Error.Output(2, fmt.Sprintf(format, v...))
}
