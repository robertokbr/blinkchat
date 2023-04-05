package logger

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/robertokbr/blinkchat/src/utils"
)

var (
	info    *log.Logger
	warning *log.Logger
	err     *log.Logger
	debug   *DebugLogger
)

type DebugLogger struct {
	Output func(c int, s string) error
}

func NewDebugLogger(outputFunction func(c int, s string) error) *DebugLogger {
	return &DebugLogger{
		Output: outputFunction,
	}
}

func fakeOutput(c int, s string) error {
	return nil
}

func init() {
	godotenv.Load()
	debug = utils.If(
		os.Getenv("LOG_LEVEL") == "debug",
		NewDebugLogger(log.New(os.Stdout, "[DEBUG]: ", log.Ldate|log.Ltime|log.Lshortfile).Output),
		NewDebugLogger(fakeOutput),
	)
	info = log.New(os.Stdout, "[INFO]: ", log.Ldate|log.Ltime|log.Lshortfile)
	warning = log.New(os.Stdout, "[WARNING]: ", log.Ldate|log.Ltime|log.Lshortfile)
	err = log.New(os.Stderr, "[ERROR]: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func Debug(message string) {
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
