package logger

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

type Logger struct {
	ShowTime  bool // True: filename will be '<yyyy/mm/dd>' False: filename will be 'log'
	TypeOfLog int8 // 0: general log - 1: batch log - 2: frontend log
}

func (l *Logger) Log(msg string, errorMsg error) {
	defaultDir, err := os.Getwd()
	if err != nil {
		log.Panic("error getting Default Working Directory")
		return
	}

	logFilePath := filepath.Join(defaultDir, "logs")
	err = os.MkdirAll(logFilePath, os.ModePerm)
	if err != nil {
		panic(err)
	}

	var logFileName string
	switch l.ShowTime {
	case false:
		logFileName = ""
	case true:
		t := time.Now()
		logFileName = t.Format(time.DateOnly)
	}

	switch l.TypeOfLog {
	case 0:
		logFileName = "log"
	case 1:
		logFileName = "batch-log"
	case 2:
		logFileName = "test-log"
	}

	logFilePath = filepath.Join(logFilePath, logFileName)

	f, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Print once to console and once to the log file
	log.Printf(msg + " | ")
	log.Println(errorMsg)
	log.SetOutput(f)
	log.Print(msg)
	if errorMsg != nil {
		log.Println(errorMsg)
	}
}
