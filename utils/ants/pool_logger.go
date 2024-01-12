package ants

import (
	"fmt"
	
	"github.com/ZYallers/golib/utils/logger"
	"github.com/panjf2000/ants/v2"
)

type PoolLogger interface {
	ants.Logger
	SendMessage(msg string)
	LogName() string
}

type poolLogger struct{}

func (l *poolLogger) LogName() string {
	return "ants-pool"
}

func (l *poolLogger) Printf(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	logger.Use(l.LogName()).Info(s)
	l.SendMessage(s)
}

func (l *poolLogger) SendMessage(msg string) {}
