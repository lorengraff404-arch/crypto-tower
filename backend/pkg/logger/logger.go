package logger

import (
	"log"
	"os"
)

var (
	InfoLogger    *log.Logger
	WarningLogger *log.Logger
	ErrorLogger   *log.Logger
)

func Init() {
	InfoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func Info(message string) {
	InfoLogger.Println(message)
}

func Warning(message string) {
	WarningLogger.Println(message)
}

func Error(message string) {
	ErrorLogger.Println(message)
}
