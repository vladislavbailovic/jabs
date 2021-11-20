package types

type LogLevel int

const (
	LOG_DEBUG LogLevel = iota
	LOG_INFO
	LOG_NOTICE
	LOG_WARNING
	LOG_ERROR
)
