package logger

import (
	"DockerProtector/lib/color"
	"DockerProtector/lib/path"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

// Settings stores config for logger
type Settings struct {
	Path       string `yaml:"path"`
	Name       string `yaml:"name"`
	Ext        string `yaml:"ext"`
	TimeFormat string `yaml:"time-format"`
	Debug      bool   `yame:"debug"`
}

var (
	logFile            *os.File
	defaultPrefix      = ""
	defaultCallerDepth = 2
	logger             *log.Logger
	mu                 sync.Mutex
	logPrefix          = ""
	levelFlags         = []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
	debug              bool
)

type logLevel int

// log levels
const (
	DEBUG logLevel = iota
	INFO
	WARNING
	ERROR
	FATAL
)

const flags = log.LstdFlags

// Setup initializes logger
func Setup(settings *Settings) {
	var (
		err error
		mw io.Writer
	)
	dir := path.GetCurrentAbPath() + settings.Path
	fileName := fmt.Sprintf("%s-%s.%s",
		settings.Name,
		time.Now().Format(settings.TimeFormat),
		settings.Ext)

	logFile, err := mustOpen(fileName, dir)
	if err != nil {
		log.Fatalf("logging.Setup err: %s", err)
	}

	// 只输出到日志文件
	mw = io.MultiWriter(logFile)

	logger = log.New(mw, defaultPrefix, flags)
	debug = settings.Debug
}

func getPrefix(level logLevel) string {
	_, file, line, ok := runtime.Caller(defaultCallerDepth)
	if ok {
		if debug { // 开启调试模式时，输出文件名和行号
			logPrefix = fmt.Sprintf("%s %s:%d", levelFlags[level], filepath.Base(file), line)
		} else { // 当不开启调试模式时，过滤 DEBUG WARN 类型日志
			logPrefix = fmt.Sprintf("%s", levelFlags[level])
		}
	} else {
		logPrefix = fmt.Sprintf("%s", levelFlags[level])
	}
	return logPrefix
}

// Debug prints debug log
func Debug(v ...interface{}) {
	if !debug { return } // 当不开启调试模式时，过滤 DEBUG WARN 类型日志
	mu.Lock()
	defer mu.Unlock()
	v = append([]interface{}{ getPrefix(DEBUG) }, v...)
	logger.Println(v...)
}

func Debugf(format string, v ...interface{}) {
	if !debug { return } // 当不开启调试模式时，过滤 DEBUG WARN 类型日志
	mu.Lock()
	defer mu.Unlock()
	v = append([]interface{}{ getPrefix(DEBUG) }, v...)
	logger.Printf(format, v...)
}

// Info prints normal log
func Info(v ...interface{}) {
	mu.Lock()
	defer mu.Unlock()
	v = append([]interface{}{ getPrefix(INFO) }, v...)
	logger.Println(v...)
}

func Infof(format string, v ...interface{}) {
	mu.Lock()
	defer mu.Unlock()
	v = append([]interface{}{ getPrefix(INFO) }, v...)
	logger.Printf(format, v...)
}

// Warn prints warning log
func Warn(v ...interface{}) {
	// 当不开启调试模式时，过滤 DEBUG WARN 类型日志
	if !debug { return }
	mu.Lock()
	defer mu.Unlock()
	v = append([]interface{}{ getPrefix(WARNING) }, v...)
	logger.Println(v...)
}

func Warnf(format string, v ...interface{}) {
	// 当不开启调试模式时，过滤 DEBUG WARN 类型日志
	if !debug { return }
	mu.Lock()
	defer mu.Unlock()
	v = append([]interface{}{ getPrefix(WARNING) }, v...)
	logger.Printf(format, v...)
}

// Error prints error log
func Error(v ...interface{}) {
	mu.Lock()
	defer mu.Unlock()
	v = append([]interface{}{ getPrefix(ERROR) }, v...)
	logger.Println(v...)
}

func Errorf(format string, v ...interface{}) {
	mu.Lock()
	defer mu.Unlock()
	v = append([]interface{}{ getPrefix(ERROR) }, v...)
	logger.Printf(format, v...)
}

// Fatal prints error log then stop the program
func Fatal(v ...interface{}) {
	mu.Lock()
	defer mu.Unlock()
	color.Red( fmt.Sprint(v...) )
	v = append([]interface{}{ getPrefix(FATAL) }, v...)
	logger.Fatalln(v...)
}

func Fatalf(format string, v ...interface{}) {
	mu.Lock()
	defer mu.Unlock()
	color.Red( fmt.Sprintf(format, v...) )
	v = append([]interface{}{ getPrefix(FATAL) }, v...)
	logger.Fatalln(fmt.Sprintf(format, v...))
}