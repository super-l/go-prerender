package internal

import (
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
	"time"
)

var stdoutLogger = logrus.New()
var fileLogger = logrus.New()

type sLogger struct{}

var SLogger = sLogger{}

func init() {
	initStdoutLogger()
	initFileLogger()
}

// Custom log format definition
type MyFormatter struct{}

func (s *MyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := time.Now().Local().Format("01/02 15:04:05")
	msg := fmt.Sprintf("%s [%s] %s\n", timestamp, strings.ToUpper(entry.Level.String()), entry.Message)
	return []byte(msg), nil
}

func initStdoutLogger() {
	stdoutLogger.Out = os.Stdout
	stdoutLogger.SetLevel(logrus.DebugLevel)

	//自定writer就行， hook 交给 lfshook
	stdoutLogger.SetFormatter(new(MyFormatter))
}

func initFileLogger() {
	fileLogger.SetLevel(logrus.InfoLevel)
	fileLogger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "01-02 15:04:05",
	})

	writer, _ := rotatelogs.New(
		"logs/"+"%Y%m%d%H"+"-"+"cookie.log",
		rotatelogs.WithLinkName("logs/"+"cookie.log"),
		rotatelogs.WithMaxAge(30*24*time.Hour),                   // 一个月
		rotatelogs.WithRotationTime(time.Duration(12*time.Hour)), // 按分钟
	)
	fileLogger.SetOutput(writer)
}

func (s sLogger) WriteFile(filename string, content string) {
	_, err := os.Stat("logs")
	if err != nil && os.IsNotExist(err) {
		_ = os.MkdirAll("logs", 0755)
	}

	filepath := fmt.Sprintf("logs/%s", filename)
	f, err := os.OpenFile(filepath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0660)
	if err != nil {
		return
	}
	_, _ = f.WriteString(content)
	_ = f.Close()
}

func (s sLogger) WriteFile2(filename string, content string) {
	_, err := os.Stat("logs")
	if err != nil && os.IsNotExist(err) {
		_ = os.MkdirAll("logs", 0755)
	}

	filepath := fmt.Sprintf("logs/%s", filename)
	f, err := os.OpenFile(filepath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0660)
	if err != nil {
		return
	}
	defer f.Close()
	_, _ = f.WriteString(content)
}

func (s sLogger) GetLogger() *logrus.Logger {
	return s.GetFileLogger()
}

func (s sLogger) GetStdoutLogger() *logrus.Logger {
	return stdoutLogger
}

func (s sLogger) GetFileLogger() *logrus.Logger {
	return fileLogger
}
