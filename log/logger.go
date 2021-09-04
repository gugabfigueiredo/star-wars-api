package log

import (
	"fmt"
	"github.com/rs/zerolog"
)

type Logger struct {
	*zerolog.Logger
	context map[string]interface{}
	
}

func (l *Logger) C(tags ...interface{}) *Logger {

	if len(tags) == 0 {
		return l
	}

	if len(tags)%2 == 1 {
		panic("logger.C: odd argument count")
	}

	logger := &Logger{
		Logger: l.Logger,
		context: l.context,
	}
	for i := 0; i < len(tags); i +=2 {
		tag, ok := tags[i].(string)
		if !ok {
			panic(fmt.Sprintf("logging tag is not a string. tag: %v", tag))
		}
		logger.context[tag] = tags[i+1]
	}

	return logger
}

func (l *Logger) D(message string, tags ...interface{}) {
	l.chainLog(l.Debug(), message, tags...)
}

func (l *Logger) F(message string, tags ...interface{}) {
	l.chainLog(l.Fatal(), message, tags...)
}

func (l *Logger) E(message string, tags ...interface{}) {
	l.chainLog(l.Error(), message, tags...)
}

func (l *Logger) W(message string, tags ...interface{}) {
	l.chainLog(l.Warn(), message, tags...)
}

func (l *Logger) L(message string, tags ...interface{}) {
	l.chainLog(l.Log(), message, tags...)
}

func (l *Logger) I(message string, tags ...interface{}) {
	l.chainLog(l.Info(), message, tags...)
}

func (l *Logger) chainLog(e *zerolog.Event, message string, tags ...interface{}) {

	if tags != nil {
		logger := l.C(tags...)

		for key, value := range logger.context {
			e = e.Str(key, fmt.Sprintf("%v", value))
		}
	}

	e.Msg(message)
}