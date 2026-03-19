package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mobigen/golang-web-template/internal/infrastructure/config"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger 로거 인터페이스. 어댑터 레이어에서 구체 구현체 대신 이 타입에 의존한다.
type Logger interface {
	Debugf(format string, args ...any)
	Infof(format string, args ...any)
	Warnf(format string, args ...any)
	Errorf(format string, args ...any)
}

// LogrusLogger logrus 기반 Logger 구현체
type LogrusLogger struct {
	*logrus.Logger
}

// defaultLogger 패키지 수준 기본 로거 인스턴스
var defaultLogger *LogrusLogger

func init() {
	defaultLogger = &LogrusLogger{logrus.New()}
	defaultLogger.SetOutput(os.Stdout)
	f := &Formatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
		ShowFields:      true,
	}
	defaultLogger.SetFormatter(f)
	defaultLogger.SetReportCaller(true)
}

// CheckLogLevel check loglevel and return logrus log level
func CheckLogLevel(lv string) (int, error) {
	switch lv {
	case config.LogLevelDebug:
		return int(logrus.DebugLevel), nil
	case config.LogLevelInfo:
		return int(logrus.InfoLevel), nil
	case config.LogLevelWarn:
		return int(logrus.WarnLevel), nil
	case config.LogLevelError:
		return int(logrus.ErrorLevel), nil
	case config.LogLevelSilent:
		return int(logrus.FatalLevel), nil
	default:
		return -1, fmt.Errorf("ERROR. Not Supported Log Level")
	}
}

// Setting :
func (l *LogrusLogger) Setting(conf *config.LogConfiguration, appHome string) error {
	var writers []io.Writer
	for _, out := range conf.Output {
		switch out {
		case config.LogOutStdout:
			writers = append(writers, os.Stdout)
		case config.LogOutFile:
			savePath := conf.LogRotate.SavePath
			if strings.TrimSpace(savePath) != "" {
				if !filepath.IsAbs(savePath) {
					savePath = filepath.Join(appHome, savePath)
				}
			}
			fileName := filepath.Join(savePath, conf.LogRotate.FileName)
			lj := &lumberjack.Logger{
				Filename:   fileName,
				MaxSize:    conf.LogRotate.SizePerFileMb,
				MaxBackups: conf.LogRotate.MaxOfDay,
				MaxAge:     conf.LogRotate.MaxAge,
				Compress:   conf.LogRotate.Compress,
			}
			writers = append(writers, lj)
		default:
			return fmt.Errorf("ERROR. Not Supported Log Output[ %s ]", out)
		}
	}
	if len(writers) == 0 {
		l.SetOutput(os.Stdout)
	} else if len(writers) == 1 {
		l.SetOutput(writers[0])
	} else {
		l.SetOutput(io.MultiWriter(writers...))
	}
	lv, err := CheckLogLevel(conf.Level)
	if err != nil {
		return err
	}
	l.SetLogLevel(logrus.Level(lv))
	return nil
}

// GetInstance return logger instance
func (LogrusLogger) GetInstance() *LogrusLogger {
	return defaultLogger
}

// SetLogLevel set log level
func (l *LogrusLogger) SetLogLevel(lv logrus.Level) {
	switch lv {
	case logrus.ErrorLevel:
		l.SetLevel(lv)
	case logrus.WarnLevel:
		l.SetLevel(lv)
	case logrus.InfoLevel:
		l.SetLevel(lv)
	case logrus.DebugLevel:
		l.SetLevel(lv)
	default:
		l.Errorf("ERROR. Not Supported Log Level[ %d ]", lv)
	}
}

// GetLogLevel get log level
func (l *LogrusLogger) GetLogLevel() string {
	text, _ := l.GetLevel().MarshalText()
	return string(text)
}

// Start Print Start Banner
func (l *LogrusLogger) Start() {
	l.Errorf("%s", LINE90)
	l.Errorf(" ")
	l.Errorf("                         START. %s:%s-%s",
		strings.ToUpper(config.Name), config.Version, config.BuildHash)
	l.Errorf(" ")
	l.Errorf("%90s", "Copyright(C) 2026 Mobigen Corporation.  ")
	l.Errorf(" ")
	l.Errorf("%s", LINE90)
}

// Shutdown Print Shutdown
func (l *LogrusLogger) Shutdown() {
	l.Errorf("%s", LINE90)
	l.Errorf(" ")
	l.Errorf("                        %s Bye Bye.", strings.ToUpper(config.Name))
	l.Errorf(" ")
	l.Errorf("%90s", "Copyright(C) 2026 Mobigen Corporation.  ")
	l.Errorf(" ")
	l.Errorf("%s", LINE90)
}

// For test
// testingWriter is an io.Writer that writes through t.Log.
type testingWriter struct {
	tb testing.TB
}

func (tw *testingWriter) Write(b []byte) (int, error) {
	tw.tb.Log(strings.TrimSpace(string(b)))
	return len(b), nil
}

// MakeTestLogger creates a custom format logrus.Logger
func MakeTestLogger(tb testing.TB) *LogrusLogger {
	defaultLogger = &LogrusLogger{logrus.New()}
	defaultLogger.SetOutput(os.Stdout)
	f := &Formatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
		ShowFields:      true,
	}
	defaultLogger.SetFormatter(f)
	defaultLogger.SetReportCaller(false)
	return defaultLogger
}
