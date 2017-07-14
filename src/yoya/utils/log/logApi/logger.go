package logApi

import (
	"fmt"
	"time"

	fluent "yoya/utils/log/hooks/fluent"
	"yoya/utils/log/logrus"
)

type LogConfig struct {
	IP      string // fluentip
	Port    int    // fluentport
	Tag     string // App name
	Level   string // log level
	APPIP   string
	APPPort int
}

const (
	Panic = "Panic"
	Fatal = "Fatal"
	Error = "Error"
	Warn  = "Warn"
	Info  = "Info"
	Debug = "Debug"
)

var logConfig LogConfig

var logger *logrus.Logger

func InitLog(config *LogConfig) error {
	var err error
	if config == nil {
		return fmt.Errorf("LogConfig is nil")
	} else {
		logConfig.IP = config.IP
		logConfig.Port = config.Port
		logConfig.Tag = config.Tag
		logConfig.Level = config.Level
		logConfig.APPIP = config.APPIP
		logConfig.APPPort = config.APPPort
	}
	logger, err = NewLog()
	if err != nil {
		return fmt.Errorf("InitLog err!")
	}

	return nil
}

func NewLog() (*logrus.Logger, error) {
	var level logrus.Level
	hook := fluent.NewHook(logConfig.IP, logConfig.Port)
	if hook == nil {
		return nil, fmt.Errorf("Newlog err!")
	}
	switch logConfig.Level {
	case Panic:
		{
			hook.SetLevels([]logrus.Level{
				logrus.PanicLevel,
			})
			level = logrus.PanicLevel
		}
	case Fatal:
		{
			hook.SetLevels([]logrus.Level{
				logrus.PanicLevel,
				logrus.FatalLevel,
			})
			level = logrus.FatalLevel
		}
	case Error:
		{
			hook.SetLevels([]logrus.Level{
				logrus.PanicLevel,
				logrus.FatalLevel,
				logrus.ErrorLevel,
			})
			level = logrus.ErrorLevel
		}
	case Warn:
		{
			hook.SetLevels([]logrus.Level{
				logrus.PanicLevel,
				logrus.FatalLevel,
				logrus.ErrorLevel,
				logrus.WarnLevel,
			})
			level = logrus.WarnLevel
		}
	case Info:
		{
			hook.SetLevels([]logrus.Level{
				logrus.PanicLevel,
				logrus.FatalLevel,
				logrus.ErrorLevel,
				logrus.WarnLevel,
				logrus.InfoLevel,
			})
			level = logrus.InfoLevel
		}
	case Debug:
		{
			hook.SetLevels([]logrus.Level{
				logrus.PanicLevel,
				logrus.FatalLevel,
				logrus.ErrorLevel,
				logrus.WarnLevel,
				logrus.InfoLevel,
				logrus.DebugLevel,
			})
			level = logrus.DebugLevel
		}
	}
	logger := logrus.New(level)
	hook.SetTag(logConfig.Tag)
	logger.Hooks.Add(hook)

	return logger, nil
}

func Debugs(place string, message string) {
	if len(message) > 0 {
		logger.Debug(fmt.Sprintf(`{"ip":"%s","port":%d,"time":%d,"place":"%s","value":"%s"}`,
			logConfig.APPIP,
			logConfig.APPPort,
			time.Now().UnixNano()/1000000,
			place,
			message))
	} else {
		logger.Debug(fmt.Sprintf(`{"ip":"%s","port":%d,"time":%d,"place":"%s"}`,
			logConfig.APPIP,
			logConfig.APPPort,
			time.Now().UnixNano()/1000000,
			place))
	}
}

func Infos(place string, message string) {
	if len(message) > 0 {
		logger.Info(fmt.Sprintf(`{"ip":"%s","port":%d,"time":%d,"place":"%s","value":%s}`,
			logConfig.APPIP,
			logConfig.APPPort,
			time.Now().UnixNano()/1000000,
			place,
			message))
	} else {
		logger.Info(fmt.Sprintf(`{"ip":"%s","port":%d,"time":%d,"place":"%s"}`,
			logConfig.APPIP,
			logConfig.APPPort,
			time.Now().UnixNano()/1000000,
			place))
	}
}

func Warns(place string, message string) {
	if len(message) > 0 {
		logger.Warn(fmt.Sprintf(`{"ip":"%s","port":%d,"time":%d,"place":"%s","value":"%s"}`,
			logConfig.APPIP,
			logConfig.APPPort,
			time.Now().UnixNano()/1000000,
			place,
			message))
	} else {
		logger.Warn(fmt.Sprintf(`{"ip":"%s","port":%d,"time":%d,"place":"%s"}`,
			logConfig.APPIP,
			logConfig.APPPort,
			time.Now().UnixNano()/1000000,
			place))
	}
}

func Errors(place string, message string) {
	if len(message) > 0 {
		logger.Error(fmt.Sprintf(`{"ip":"%s","port":%d,"time":%d,"place":"%s","value":"%s"}`,
			logConfig.APPIP,
			logConfig.APPPort,
			time.Now().UnixNano()/1000000,
			place,
			message))
	} else {
		logger.Error(fmt.Sprintf(`{"ip":"%s","port":%d,"time":%d,"place":"%s"}`,
			logConfig.APPIP,
			logConfig.APPPort,
			time.Now().UnixNano()/1000000,
			place))
	}
}

func Fatals(place string, message string) {
	if len(message) > 0 {
		logger.Fatal(fmt.Sprintf(`{"ip":"%s","port":%d,"time":%d,"place":"%s","value":"%s"}`,
			logConfig.APPIP,
			logConfig.APPPort,
			time.Now().UnixNano()/1000000,
			place,
			message))
	} else {
		logger.Fatal(fmt.Sprintf(`{"ip":"%s","port":%d,"time":%d,"place":"%s"}`,
			logConfig.APPIP,
			logConfig.APPPort,
			time.Now().UnixNano()/1000000,
			place))
	}
}

func Panics(place string, message string) {
	if len(message) > 0 {
		logger.Panic(fmt.Sprintf(`{"ip":"%s","port":%d,"time":%d,"place":"%s","value":"%s"}`,
			logConfig.APPIP,
			logConfig.APPPort,
			time.Now().UnixNano()/1000000,
			place,
			message))
	} else {
		logger.Panic(fmt.Sprintf(`{"ip":"%s","port":%d,"time":%d,"place":"%s"}`,
			logConfig.APPIP,
			logConfig.APPPort,
			time.Now().UnixNano()/1000000,
			place))
	}
}
