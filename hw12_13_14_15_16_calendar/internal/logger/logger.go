package logger

import (
	"fmt"
	"strings"
	"time"
)

type Logger struct {
	level string
}

func New(level string) *Logger {
	return &Logger{level: strings.ToUpper(level)}
}

func (l *Logger) levelAllowed(level string) bool {
	allowed := map[string]int{
		"DEBUG": 1,
		"INFO":  2,
		"ERROR": 3,
	}

	current := allowed[l.level]
	want := allowed[strings.ToUpper(level)]
	return want >= current
}

func (l *Logger) Debug(msg string) {
	if l.levelAllowed("DEBUG") {
		fmt.Println(time.Now().Format("15:04:05"), "[DEBUG]", msg)
	}
}

func (l *Logger) Info(msg string) {
	if l.levelAllowed("INFO") {
		fmt.Println(time.Now().Format("15:04:05"), "[INFO]", msg)
	}
}

func (l *Logger) Error(msg string) {
	if l.levelAllowed("ERROR") {
		fmt.Println(time.Now().Format("15:04:05"), "[ERROR]", msg)
	}
}
