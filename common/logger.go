package common

import (
	"fmt"
	"log"
	"os"
)

var infoLogger = createLogger("info.log", "[info]")
var errorLogger = createLogger("error.log", "[error]")

func createLogger(file string, prefix string) *log.Logger {
	var f *os.File
	var err error
	if f, err = openFile(file); err != nil {
		if f, err = os.Create(file); err != nil {
			defer f.Close()
			panic(err)
		}
	}
	return log.New(f, prefix, log.LstdFlags|log.Lmicroseconds)
}

func openFile(file string) (*os.File, error) {
	if f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666); err != nil && os.IsNotExist(err) {
		return nil, err
	} else {
		return f, nil
	}
}

func Info(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	log.Println(message)
	infoLogger.Println(message)
}

func Error(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	log.Println(message)
	errorLogger.Println(message)
}

func Fatal(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	errorLogger.Println(message)
	log.Fatalln(message)
}
