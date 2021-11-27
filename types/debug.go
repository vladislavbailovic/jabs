package types

type LogLevel int

const (
	LOG_DEBUG LogLevel = iota
	LOG_INFO
	LOG_NOTICE
	LOG_WARNING
	LOG_ERROR
)

var LOG_LEVELS = map[LogLevel]string{
	LOG_DEBUG:   "debug",
	LOG_INFO:    "info",
	LOG_NOTICE:  "notice",
	LOG_WARNING: "warning",
	LOG_ERROR:   "error",
}
