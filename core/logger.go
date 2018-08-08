package core

import (
	"io"
	"log"
	"os"
)

var (
	// Logger Instance
	Logger log.Logger

	// LogFile Pointer
	LogFile *os.File
)

const (
	prefix     = "CHAT_LACKS: "
	loggerFlag = log.Ldate | log.Ltime | log.Lshortfile
	fileFlag   = os.O_RDWR | os.O_CREATE | os.O_APPEND
)

func defaultLogger() *log.Logger {
	log.Println("falling back to default stdout")
	return log.New(os.Stdout, prefix, loggerFlag)
}

// CreateLogger creates New Logger
// logOutput is the name of the logfile to log to
// default logOutput is stdout
func CreateLogger(logOutput string) (*log.Logger, *os.File) {
	if logOutput == "" {
		log.Println("log output not set")
		return defaultLogger(), nil
	}
	if logOutput == "stdout" {
		return defaultLogger(), nil
	}
	LogFile, err := os.OpenFile(logOutput, fileFlag, 0666)
	if err != nil {
		log.Printf("failed to open log file: %s \n", err.Error())
		return defaultLogger(), nil
	}
	w := io.MultiWriter(LogFile, os.Stdout)
	return log.New(w, prefix, loggerFlag), LogFile
}
