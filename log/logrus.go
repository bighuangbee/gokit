package log

import (
	"fmt"
	"github.com/bighuangbee/gokit/function"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var (
	logPath string
)

type Loger struct {
	Logrus *logrus.Logger
}

func New(logPath string) (*Loger, error) {
	logPath = logPath

	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		err := os.MkdirAll(logPath, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

	src, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return nil, err
	}

	loger := logrus.New()
	loger.SetReportCaller(true)

	//输出到文件
	loger.Out = src

	logerFormatter := new(LogerFormatter)
	loger.SetFormatter(logerFormatter)

	loger.AddHook(newLocalFileLogHook(logrus.ErrorLevel, logerFormatter))
	loger.AddHook(newLocalFileLogHook(logrus.InfoLevel, logerFormatter))

	loger.SetOutput(os.Stdout)

	return &Loger{loger}, nil
}

/*
*
写入本地日志文件， 按日期、日志级别分割为不同的文件
*/
func newLocalFileLogHook(level logrus.Level, formatter logrus.Formatter) logrus.Hook {

	fileName := filepath.Join(logPath + "%Y%m%d.log")

	//文件分割
	writer, err := rotatelogs.New(
		fileName,
		// 最大保存时间(30天)
		rotatelogs.WithMaxAge(30*24*time.Hour),
		// 日志分割间隔(1天)
		rotatelogs.WithRotationTime(24*time.Hour),
	)

	if err != nil {
		fmt.Errorf("config local file system for Loger error: %v", err)
	}

	return lfshook.NewHook(lfshook.WriterMap{
		level: writer,
	}, formatter)

}

func (s *Loger) Infof(format string, args ...interface{}) {
	setPrefix("Infof")
	s.Info(format, args)
}

func (s *Loger) Info(args ...interface{}) {
	s.Logrus.Info(setPrefix("Info"), args)
}

func (s *Loger) Error(args ...interface{}) {
	s.Logrus.Error(setPrefix(function.Red("Error")), args)
}

// setPrefix set the prefix of the log output
func setPrefix(level string) string {

	pc, file, line, ok := runtime.Caller(2)
	if ok {
		loc, _ := time.LoadLocation("Asia/Shanghai")
		funcName := runtime.FuncForPC(pc).Name()
		funcName = strings.TrimPrefix(filepath.Ext(funcName), ".")
		timestamp := time.Now().In(loc).Format("2006-01-02 15:04:05.000")

		return fmt.Sprintf("[%s][%s][%s:%d:%s]", timestamp, level, filepath.Base(file), line, funcName)
	}
	return ""
}

/*
日志输出格式
*/
type LogerFormatter struct{}

func (s *LogerFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	msg := fmt.Sprintf("%s \n", entry.Message)
	return []byte(msg), nil
}
