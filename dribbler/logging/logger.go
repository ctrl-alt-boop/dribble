package logging

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
)

const (
	LogLevelInfo = iota
	LogLevelWarn
	LogLevelError
	LogLevelPanic
)

func init() {
	newGlobalLogger("dribble", true)
}

var Log *Logger

func newGlobalLogger(appName string, removeExisting bool) *Logger {
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		os.Mkdir("logs", 0755)
	}
	filename := fmt.Sprintf("%s.log", appName)
	if removeExisting {
		os.Remove("logs/" + filename)
	}
	logfile, err := os.OpenFile("logs/"+filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	Log = &Logger{
		logger: log.New(logfile, "", log.LstdFlags), //|log.Lshortfile),
		file:   logfile,
	}
	Log.Infof("Global Logger: %s created...", appName)
	return Log
}

func GlobalLogger() *Logger {
	if Log == nil {
		_, file, line, _ := runtime.Caller(1)
		panic("Global logger not initialized: called from: " + fmt.Sprint(file, ":", line))
	}
	return Log
}

func CloseGlobalLogger() {
	if Log == nil {
		return
	}
	Log.Close()
	Log = nil
}

type Logger struct {
	logger *log.Logger
	file   *os.File
}

func (l *Logger) Write(p []byte) (n int, err error) {
	l.logger.SetPrefix("[INFO]: ")
	message := l.formatMessage(string(p))
	l.logger.Println(message)
	l.logger.SetPrefix("")
	return len(p), nil
}

func NewLogger(packageName string) *Logger {
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		os.Mkdir("logs", 0755)
	}
	logfile, err := os.Create("logs/" + packageName + ".log")
	if err != nil {
		panic(err)
	}

	return &Logger{
		logger: log.New(logfile, "", log.LstdFlags), //|log.Lshortfile),
		file:   logfile,
	}
}

func (l *Logger) formatMessage(args ...any) string {
	pc, _, _, _ := runtime.Caller(2)
	// shortFile := file
	funcName := runtime.FuncForPC(pc).Name()
	shortFuncName := strings.ReplaceAll(strings.TrimPrefix(funcName, "github.com/ctrl-alt-boop/dribble/"), "/", ".")
	// if idx := strings.LastIndex(file, "/"); idx != -1 {
	// 	shortFile = file[idx+1:]
	// }
	// if idx := strings.LastIndex(shortFile, "\\"); idx != -1 {
	// 	shortFile = shortFile[idx+1:]
	// }
	// message := fmt.Sprintf("%s:%d -> %s", shortFuncName, line, fmt.Sprint(args...))
	message := fmt.Sprintf("%s -> %s", shortFuncName, fmt.Sprint(args...))
	return message
}

func (l *Logger) formatMessageWithCallStack(skipFrames, numFrames int, args ...any) string {
	pc := make([]uintptr, numFrames)
	n := runtime.Callers(skipFrames, pc)

	var stackInfo strings.Builder
	for i := range n {
		f := runtime.FuncForPC(pc[i])
		if f == nil {
			stackInfo.WriteString("<unknown function> ")
			continue
		}

		file, line := f.FileLine(pc[i])
		shortFile := file
		if idx := strings.LastIndex(file, "/"); idx != -1 {
			shortFile = file[idx+1:]
		}
		if idx := strings.LastIndex(shortFile, "\\"); idx != -1 { // Handle Windows paths
			shortFile = shortFile[idx+1:]
		}

		stackInfo.WriteString(fmt.Sprintf("%s:%d -> ", shortFile, line))
	}

	messageContent := fmt.Sprint(args...)
	return stackInfo.String() + "\n\t" + messageContent
}

func (l *Logger) Info(args ...any) {
	l.logger.SetPrefix("[INFO]: ")
	message := l.formatMessage(args...)
	l.logger.Println(message)
	l.logger.SetPrefix("")
}

func (l *Logger) Warn(args ...any) {
	l.logger.SetPrefix("[WARN]: ")
	message := l.formatMessage(args...)
	l.logger.Println(message)
	l.logger.SetPrefix("")
}

func (l *Logger) Error(args ...any) {
	l.logger.SetPrefix("[ERROR]: ")
	message := l.formatMessage(args...)
	l.logger.Println(message)
	l.logger.SetPrefix("")
}

func (l *Logger) Fatal(args ...any) {
	l.logger.SetPrefix("[ERROR]: ")
	message := l.formatMessage(args...)
	l.logger.Fatal(message)
	l.logger.SetPrefix("")
}

func (l *Logger) Panic(args ...any) {
	l.logger.SetPrefix("[PANIC]: ")
	message := l.formatMessageWithCallStack(3, 5, args...)
	formatted := make([]string, 0)
	lines := strings.SplitSeq(message, "\n")
	for line := range lines {
		trimmed := strings.TrimSpace(line)
		formatted = append(formatted, trimmed)
	}

	l.logger.Panic(strings.Join(formatted, "\n"))
	l.logger.SetPrefix("")
}

func (l *Logger) Infof(format string, args ...any) {
	l.logger.SetPrefix("[INFO]: ")
	message := l.formatMessage(fmt.Sprintf(format, args...))
	l.logger.Println(message)
	l.logger.SetPrefix("")
}

func (l *Logger) Warnf(format string, args ...any) {
	l.logger.SetPrefix("[WARN]: ")
	message := l.formatMessage(fmt.Sprintf(format, args...))
	l.logger.Println(message)
	l.logger.SetPrefix("")
}

func (l *Logger) ErrorF(format string, args ...any) {
	l.logger.SetPrefix("[ERROR]: ")
	message := l.formatMessage(fmt.Sprintf(format, args...))
	l.logger.Println(message)
	l.logger.SetPrefix("")
}

func (l *Logger) Fatalf(format string, args ...any) {
	l.logger.SetPrefix("[ERROR]: ")
	message := l.formatMessage(fmt.Sprintf(format, args...))
	l.logger.Fatal(message)
	l.logger.SetPrefix("")
}

func (l *Logger) Panicf(format string, args ...any) {
	l.logger.SetPrefix("[PANIC]: ")
	message := l.formatMessageWithCallStack(3, 5, fmt.Sprintf(format, args...))
	formatted := make([]string, 0)
	lines := strings.SplitSeq(message, "\n")
	for line := range lines {
		trimmed := strings.TrimSpace(line)
		formatted = append(formatted, trimmed)
	}

	l.logger.Panic(strings.Join(formatted, "\n"))
	l.logger.SetPrefix("")
}

func (l *Logger) Close() {
	l.logger.SetPrefix("[LOGGER]: ")
	l.logger.Println("Closing logger...")
	l.logger.SetOutput(os.Stdout)
	l.file.Close()
}
