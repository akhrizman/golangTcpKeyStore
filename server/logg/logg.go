package logg

import (
	"fmt"
	"log"
	"os"
)

const (
	logsDir     = "C:/Users/Alex.Khrizman/go_logs/"
	logFilename = "tcpstore.log"
)

var (
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
	LogFile *os.File
)

func SetupLogging() {
	logFile, fileErr := os.OpenFile(fmt.Sprintf("%s%s", logsDir, logFilename), os.O_CREATE|os.O_WRONLY, 0666)
	if fileErr != nil {
		Error.Fatal(fileErr)
	}
	Info = log.New(logFile, "INFO", log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(logFile, "WARNING", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(logFile, "ERROR", log.Ldate|log.Ltime|log.Lshortfile)
}

func CloseLogFiles() {
	LogFile.Close()
}
