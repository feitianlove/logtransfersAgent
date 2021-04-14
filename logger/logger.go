package logger

import (
	goliblogger "github.com/feitianlove/golib/common/logger"
	"github.com/feitianlove/logtransfersAgent/config"
	"github.com/sirupsen/logrus"
)

var (
	Ctrl    *logrus.Logger
	Console *logrus.Logger
)

func init() {
	Ctrl = goliblogger.NewLoggerInstance()
	Console = goliblogger.NewLoggerInstance()
}
func InitCtrlLog(conf *goliblogger.LogConf) error {
	logger, err := goliblogger.InitLogger(conf)
	if err != nil {
		return err
	}
	Ctrl = logger
	return nil
}
func InitConsoleLog(conf *goliblogger.LogConf) error {
	logger, err := goliblogger.InitLogger(conf)
	if err != nil {
		return err
	}
	Console = logger
	return nil
}

func InitLog(cfg *config.Config) error {
	if err := InitConsoleLog(cfg.ConsoleLog); err != nil {
		return err
	}
	if err := InitCtrlLog(cfg.CtrlLog); err != nil {
		return err
	}
	return nil
}
