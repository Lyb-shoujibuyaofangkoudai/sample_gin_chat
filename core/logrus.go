package core

import (
	"bytes"
	"fmt"
	"gin_chat/global"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"path"
)

const (
	red    = 31
	yellow = 33
	blue   = 36
	gray   = 37
)

type LogFormatter struct{}

func (*LogFormatter) printLog(entry *logrus.Entry, b *bytes.Buffer, timestamp, logPrefix string, levelColor int) {
	env := viper.GetString("system.env")
	if env == "debug" {
		if entry.HasCaller() {
			//	 自定义文件路径
			funcVal := entry.Caller.Function
			fileVal := fmt.Sprintf("%s:%d", path.Base(entry.Caller.File), entry.Caller.Line)
			//	 自定义输出格式
			fmt.Fprintf(b, "%s[%s] \x1b[%dm[%s]\x1b[0m %s %s %s \n", logPrefix, timestamp, levelColor, entry.Level, fileVal, funcVal, entry.Message)
		} else {
			fmt.Fprintf(b, "%s[%s] \x1b[%dm[%s]\x1b[0m %s \n", logPrefix, timestamp, levelColor, entry.Level, entry.Message)
		}
	} else if env == "release" {
		if entry.HasCaller() {
			//	 自定义文件路径
			funcVal := entry.Caller.Function
			fileVal := fmt.Sprintf("%s:%d", path.Base(entry.Caller.File), entry.Caller.Line)
			//	 自定义输出格式
			fmt.Fprintf(b, "%s[%s] [%s] [%s] [%s] %s \n", logPrefix, timestamp, entry.Level, fileVal, funcVal, entry.Message)
		} else {
			fmt.Fprintf(b, "%s[%s] [%s]  %s \n", logPrefix, timestamp, entry.Level, entry.Message)
		}
	}
}

func (t *LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var levelColor int
	switch entry.Level {
	case logrus.DebugLevel, logrus.TraceLevel:
		levelColor = gray
	case logrus.WarnLevel:
		levelColor = yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		levelColor = red
	default:
		levelColor = blue
	}

	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	logPrefix := viper.GetString("logger.prefix")
	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	t.printLog(entry, b, timestamp, logPrefix, levelColor)
	//if entry.HasCaller() {
	//	//	 自定义文件路径
	//	funcVal := entry.Caller.Function
	//	fileVal := fmt.Sprintf("%s:%d", path.Base(entry.Caller.File), entry.Caller.Line)
	//
	//	//	 自定义输出格式
	//	fmt.Fprintf(b, "%s[%s] \x1b[%dm[%s]\x1b[0m %s %s %s \n", logPrefix, timestamp, levelColor, entry.Level, fileVal, funcVal, entry.Message)
	//} else {
	//	fmt.Fprintf(b, "%s[%s] \x1b[%dm[%s]\x1b[0m %s \n", logPrefix, timestamp, levelColor, entry.Level, entry.Message)
	//}
	return b.Bytes(), nil
}

func InitLogger() *logrus.Logger {
	mLog := logrus.New()                                    // 实例化
	mLog.SetOutput(os.Stdout)                               // 设置输出类型
	mLog.SetReportCaller(viper.GetBool("logger.show_line")) // 开启返回函数名和行号
	mLog.SetFormatter(&LogFormatter{})                      // 设置自定义的formatter
	level, err := logrus.ParseLevel(viper.GetString("logger.level"))
	if err != nil {
		level = logrus.InfoLevel
	}
	mLog.SetLevel(level) // 设置最低日志级别
	InitDefaultLogger()
	if viper.GetString("system.env") == "release" {
		// 生产环境就写入日志到文件
		LogWriteToFile(mLog)
	}
	return mLog
}

func InitDefaultLogger() {
	// 全局log
	logrus.SetOutput(os.Stdout)
	logrus.SetReportCaller(viper.GetBool("logger.show_line"))
	logrus.SetFormatter(&LogFormatter{})
	level, err := logrus.ParseLevel(viper.GetString("logger.level"))
	if err != nil {
		level = logrus.InfoLevel
	}
	logrus.SetLevel(level) // 设置最低日志级别
}

func LogWriteToFile(l *logrus.Logger) {
	// todo 写入日志待完善，没有校验文件是否存在，日志分割等操作
	file, err := os.OpenFile("log/logrus.log", os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		l.Out = file
	} else {
		l.Info("Failed to log to file, using default stderr")
	}
}

// LogrusGormLogger gorm日志记录器
type LogrusGormLogger struct {
	*logrus.Logger
}

// Printf 实现gorm/logger.Writer接口
func (m *LogrusGormLogger) Printf(format string, v ...interface{}) {
	logStr := fmt.Sprintf(format, v...)
	log := m.WithFields(logrus.Fields{
		"Type":    global.RouterLog,
		"Format":  format,
		"Values":  v,
		"Message": logStr,
	})
	//利用loggus记录日志
	log.Info(logStr)
}
