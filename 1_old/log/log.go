package log

import (
	"context"
	"log"
)

type ILog interface {
	Debug(ctx context.Context, args ...interface{})
	Info(ctx context.Context, args ...interface{})
	Warn(ctx context.Context, args ...interface{})
	Error(ctx context.Context, args ...interface{})
	Debugf(ctx context.Context, format string, args ...interface{})
	Infof(ctx context.Context, format string, args ...interface{})
	Warnf(ctx context.Context, format string, args ...interface{})
	Errorf(ctx context.Context, format string, args ...interface{})
}

type DefaultLog struct {
	level Level
}

func NewDefaultLog() *DefaultLog {
	// log.SetPrefix("datasupply: ")
	log.SetFlags(log.LstdFlags)
	return &DefaultLog{}
}

var _ ILog = new(DefaultLog)

func (logger *DefaultLog) SetLevel(level Level) {
	logger.level = level
}

func (logger *DefaultLog) Debug(ctx context.Context, args ...interface{}) {
	log.Print("debug:")
	log.Println(args...)
}

func (logger *DefaultLog) Info(ctx context.Context, args ...interface{}) {
	if logger.level > InfoLevel {
		return
	}
	log.Print("info:")
	log.Println(args...)
}

func (logger *DefaultLog) Warn(ctx context.Context, args ...interface{}) {
	if logger.level > WarnLevel {
		return
	}
	log.Print("warn:")
	log.Println(args...)
}

func (logger *DefaultLog) Error(ctx context.Context, args ...interface{}) {
	if logger.level > ErrorLevel {
		return
	}
	log.Print("error:")
	log.Println(args...)
}

func (logger *DefaultLog) Debugf(ctx context.Context, format string, args ...interface{}) {
	log.Print("debug:")
	log.Printf(format, args...)
}

func (logger *DefaultLog) Infof(ctx context.Context, format string, args ...interface{}) {
	if logger.level > InfoLevel {
		return
	}
	log.Print("info:")
	log.Printf(format, args...)
}

func (logger *DefaultLog) Warnf(ctx context.Context, format string, args ...interface{}) {
	if logger.level > WarnLevel {
		return
	}
	log.Print("warn:")
	log.Printf(format, args...)

}

func (logger *DefaultLog) Errorf(ctx context.Context, format string, args ...interface{}) {
	if logger.level > ErrorLevel {
		return
	}
	log.Print("error:")
	log.Printf(format, args...)
}
