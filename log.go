package main

import (
	"io"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
)

func (l *Log) Setup() {
	logFile, err := os.OpenFile(l.LogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	l.Multi = io.MultiWriter(logFile, os.Stdout)
	log.SetOutput(l.Multi)
}

func (l *Log) Log(msg *Event) {

	var logLevelMap = map[int]string{
		0: "🐞 DEBUG > ",
		1: "ℹ️ INFO  > ",
		2: "⚠️ WARN  > ",
		4: "🚨 ERROR > ",
	}

	switch msg.Level {
	case 4: // Handle Fatal Errors

		_, file, line, _ := runtime.Caller(2)
		file = strings.Split(file, "/")[len(strings.Split(file, "/"))-1]
		caller := file + ":" + strconv.Itoa(line)

		for i, arg := range msg.Message {
			if i == len(msg.Message)-1 {
				log.Fatalf("%s %s - ( %s )\n", logLevelMap[msg.Level], arg, caller)
			} else {
				log.Printf("%s %s - ( %s )\n", logLevelMap[msg.Level], arg, caller)
			}
		}
	default:
		for _, arg := range msg.Message {
			log.Printf("%s %s\n", logLevelMap[msg.Level], arg)
		}
	}

}

func (l *Log) Debug(args ...interface{}) {
	l.Log(&Event{0, args})
}

func (l *Log) Info(args ...interface{}) {
	l.Log(&Event{1, args})
}

func (l *Log) Warn(args ...interface{}) {
	l.Log(&Event{2, args})
}

func (l *Log) Error(args ...interface{}) {
	l.Log(&Event{4, args})
}
