package logger

import "fmt"

var current_log_level = 0

func SetLogLevel(lv int) {
	current_log_level = lv
}

const (
	LOG_LEVEL_NONE  = 0
	LOG_LEVEL_INFO  = 1
	LOG_LEVEL_DEBUG = 2
)

func LOG_INFO(format string, a ...interface{}) {
	if current_log_level >= LOG_LEVEL_INFO {
		fmt.Printf(format, a...)
	}
}

func LOG_DEBUG(format string, a ...interface{}) {
	if current_log_level >= LOG_LEVEL_DEBUG {
		fmt.Printf(format, a...)
	}
}
