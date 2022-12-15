package logg

import (
	"fmt"
	"log"
	"os"
)

const (
	logsDir             = "C:/Users/Alex.Khrizman/go_logs/"
	serverLogFilename   = "tcpserver.log"
	responseLogFilename = "tcpresponses.log"
)

var (
	Info            *log.Logger
	Warning         *log.Logger
	Error           *log.Logger
	Response        *log.Logger
	ServerLogFile   *os.File
	ResponseLogFile *os.File
)

func SetupLogging() {
	SetupServerLogging()
	SetupRequestLogging()
}

func SetupServerLogging() {
	logFile, fileErr := os.OpenFile(fmt.Sprintf("%s%s", logsDir, serverLogFilename), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if fileErr != nil {
		Error.Fatal(fileErr)
	}
	Info = log.New(logFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(logFile, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(logFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func SetupRequestLogging() {
	logFile, fileErr := os.OpenFile(fmt.Sprintf("%s%s", logsDir, responseLogFilename), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if fileErr != nil {
		Error.Fatal(fileErr)
	}
	Response = log.New(logFile, "RESPONSE: ", log.Ldate|log.Ltime)
}

func CloseLogFiles() {
	ServerLogFile.Close()
	ResponseLogFile.Close()
}
