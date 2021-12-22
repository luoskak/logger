package logger

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strings"
	"time"
)

var sourceDir string

func init() {
	_, file, _, _ := runtime.Caller(0)
	sourceDir = regexp.MustCompile(`logger.logger\.go`).ReplaceAllString(file, "")
}

type LogMode uint

const (
	Silent LogMode = iota + 1
	Normal
	Warning
	Debug
)

type logLevel int

const (
	_ logLevel = iota
	lErr
	lWarn
	lDebug
	lInfo
)

type ColorType int

const (
	_ ColorType = iota + 29
	ColorBlack
	ColorRed
	ColorGreen
	ColorYellow
	ColorBlue
	ColorPurple
	ColorCyan
	ColorWhite
)

var defaultLogMode = Normal

func SetMode(mode LogMode) {
	defaultLogMode = mode
}

var defaultLogger = NewLogger()

type Logger struct {
	name      string
	sourceDir string
}

func newLogger() *Logger {
	return &Logger{sourceDir: sourceDir}
}

func NewLogger(args ...interface{}) *Logger {
	length := len(args)
	if length == 1 {
		if name, ok := args[0].(string); ok {
			l := newLogger()
			l.name = name
			return l
		}
	} else if length > 1 {
		if format, ok := args[0].(string); ok {
			l := newLogger()
			l.name = fmt.Sprintf(format, args[1:]...)
			return l
		}
	}
	l := newLogger()
	return l
}

func NewLoggerWithSourceDir(sourceDir string, args ...interface{}) *Logger {
	if sourceDir == "" {
		panic("")
	}
	l := NewLogger(args...)
	l.sourceDir = sourceDir
	return l
}

func Format(f interface{}, v ...interface{}) string {
	var msg string
	switch f := f.(type) {
	case string:
		msg = f
		if len(v) == 0 {
			return msg
		}
		if !strings.Contains(msg, "%") {
			msg += strings.Repeat(" %v", len(v))
		}
	default:
		msg = fmt.Sprint(f)
		if len(v) == 0 {
			return msg
		}
		msg += strings.Repeat(" %v", len(v))
	}
	if len(v) > 0 {
		return fmt.Sprintf(msg, v...)
	}
	return fmt.Sprint(msg)
}

func (l Logger) writeRuntimeMsg() string {
	for i := 2; i < 15; i++ {
		_, file, line, ok := runtime.Caller(i)
		if ok && (!strings.HasPrefix(file, l.sourceDir) || strings.HasSuffix(file, "_test.go")) {
			return fmt.Sprintf("\033[%d;1m%s:%d\033[0m", ColorBlack, file, line)
		}
	}
	return ""
}

func writeMsg(level logLevel, msg string) string {
	ts := time.Now().Format("2006-01-02 15:04:05")
	switch level {
	case lInfo:
		msg = fmt.Sprintf("%s \033[%d;1m[日志]\033[0m %s", ts, ColorBlue, msg)
	case lDebug:
		msg = fmt.Sprintf("%s \033[%d;1m[调试]\033[0m %s", ts, ColorGreen, msg)
	case lWarn:
		msg = fmt.Sprintf("%s \033[%d;1m[警告]\033[0m %s", ts, ColorYellow, msg)
	case lErr:
		msg = fmt.Sprintf("%s \033[%d;1m[错误]\033[0m %s", ts, ColorRed, msg)
	}
	return msg
}

func Info(f interface{}, v ...interface{}) {
	defaultLogger.Info(f, v...)
}

func Warn(f interface{}, v ...interface{}) {
	defaultLogger.Warn(f, v...)
}

func Error(f interface{}, v ...interface{}) {
	defaultLogger.Error(f, v...)
}

// Info 普通日志打印
func (l Logger) Info(f interface{}, v ...interface{}) {
	if defaultLogMode > Silent {
		bys := bytes.NewBufferString("")
		if defaultLogMode > Normal {
			bys.WriteString(l.writeRuntimeMsg())
			bys.WriteByte('\n')
		}
		if l.name != "" {
			bys.WriteString(fmt.Sprintf("Logger: %s >>>>>\n", l.name))
		}
		bys.Write([]byte(writeMsg(lInfo, Format(f, v...))))
		bys.Write([]byte("\n<<<<<\n"))
		os.Stdout.Write(bys.Bytes())
	}
}

func (l Logger) Warn(f interface{}, v ...interface{}) {
	if defaultLogMode > Normal {
		bys := bytes.NewBufferString("")
		if defaultLogMode > Warning {
			bys.WriteString(l.writeRuntimeMsg())
			bys.WriteByte('\n')
		}
		if l.name != "" {
			bys.WriteString(fmt.Sprintf("Logger: %s >>>>>\n", l.name))
		}
		bys.Write([]byte(writeMsg(lWarn, Format(f, v...))))
		bys.Write([]byte("\n<<<<<\n"))
		os.Stdout.Write(bys.Bytes())
	}

}

func (l Logger) Error(f interface{}, v ...interface{}) {
	bys := bytes.NewBuffer([]byte(l.writeRuntimeMsg()))
	bys.WriteByte('\n')
	if l.name != "" {
		bys.WriteString(fmt.Sprintf("Logger: %s >>>>>\n", l.name))
	}
	bys.Write([]byte(writeMsg(lErr, Format(f, v...))))
	bys.Write([]byte("\n<<<<<\n"))
	os.Stdout.Write(bys.Bytes())
}

func (l Logger) Level() uint {
	return uint(defaultLogMode)
}
