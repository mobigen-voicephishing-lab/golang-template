package common

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mobigen/golang-web-template/common/appdata"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger struct have embedding logrus.Logger
type Logger struct {
	*logrus.Logger
}

// Logger log variable
var l *Logger

func init() {
	l = &Logger{logrus.New()}
	l.SetOutput(os.Stdout)
	f := &Formatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
		ShowFields:      true,
	}
	l.SetFormatter(f)
	l.SetReportCaller(true)
}

// Setting log setting
func (l *Logger) Setting(conf *appdata.LogConfiguration, appHome string) error {
	var writers []io.Writer
	switch conf.Output {
	case appdata.LogOutStdout:
		writers = append(writers, os.Stdout)
	case appdata.LogOutFile:
		writer, err := l.getFileOutput(conf, appHome)
		if err != nil {
			return fmt.Errorf("ERROR. Can't Initialize Log Output Setup")
		}
		writers = append(writers, writer)
	case appdata.LogOutBoth:
		writers = append(writers, os.Stdout)
		writer, err := l.getFileOutput(conf, appHome)
		if err != nil {
			return fmt.Errorf("ERROR. Can't Initialize Log Output Setup")
		}
		writers = append(writers, writer)
	default:
		return fmt.Errorf("ERROR. Not Supported Log Output[ %s ]", conf.Output)
	}
	lv, err := appdata.CheckLogLevel(conf.Level)
	if err != nil {
		return err
	}
	if len(writers) == 0 {
		l.SetOutput(os.Stdout)
	} else if len(writers) == 1 {
		l.SetOutput(writers[0])
	} else {
		l.SetOutput(io.MultiWriter(writers...))
	}
	l.SetLogLevel(logrus.Level(lv))
	return nil
}

// SetFileOutput
func (l *Logger) getFileOutput(conf *appdata.LogConfiguration, appHome string) (io.Writer, error) {
	savePath := conf.SavePath
	if !filepath.IsAbs(savePath) {
		savePath = filepath.Join(appHome, savePath)
	}
	fileName := savePath + "/" + conf.FileName

	lj := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    int(conf.SizePerFileMb),
		MaxBackups: int(conf.MaxOfDay),
		MaxAge:     int(conf.MaxAge),
		Compress:   conf.Compress,
	}
	return lj, nil
}

// GetInstance return logger instance
func (Logger) GetInstance() *Logger {
	return l
}

// SetLogLevel set log level
func (l *Logger) SetLogLevel(lv logrus.Level) {
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
func (l *Logger) GetLogLevel() string {
	text, _ := l.GetLevel().MarshalText()
	return string(text)
}

// Start Print Start Banner
func (l *Logger) Start() {
	l.Errorf("%s", LINE90)
	l.Errorf(" ")
	l.Errorf("                         START. %s:%s-%s",
		strings.ToUpper(appdata.Name), appdata.Version, appdata.BuildHash)
	l.Errorf(" ")
	l.Errorf("%90s", "Copyright(C) 2026 Mobigen Corporation.  ")
	l.Errorf(" ")
	l.Errorf("%s", LINE90)
}

// Shutdown Print Shutdown
func (l *Logger) Shutdown() {
	l.Errorf("%s", LINE90)
	l.Errorf(" ")
	l.Errorf("                        %s Bye Bye.", strings.ToUpper(appdata.Name))
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
func MakeTestLogger(tb testing.TB) *Logger {
	l = &Logger{logrus.New()}
	l.SetOutput(os.Stdout)
	f := &Formatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
		ShowFields:      true,
	}
	l.SetFormatter(f)
	l.SetReportCaller(false)
	return l
}
