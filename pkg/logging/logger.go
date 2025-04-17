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

func NewLogger(filename string) *Logger {
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		os.Mkdir("logs", 0755)
	}
	logfile, err := os.Create("logs/" + filename)
	if err != nil {
		panic(err)
	}

	return &Logger{
		logger: log.New(logfile, "", log.LstdFlags), //|log.Lshortfile),
		file:   logfile,
	}
}

func (l *Logger) formatMessage(args ...any) string {
	_, file, line, _ := runtime.Caller(2)
	shortFile := file
	if idx := strings.LastIndex(file, "/"); idx != -1 {
		shortFile = file[idx+1:]
	}
	if idx := strings.LastIndex(shortFile, "\\"); idx != -1 {
		shortFile = shortFile[idx+1:]
	}
	message := fmt.Sprintf("%s:%d -> %s", shortFile, line, fmt.Sprint(args...))
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
	l.file.Close()
}
