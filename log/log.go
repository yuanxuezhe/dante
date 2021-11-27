package log

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

// levels
const (
	debugLevel   = 0
	releaseLevel = 1
	errorLevel   = 2
	fatalLevel   = 3
)

const (
	printDebugLevel   = "[debug  ] "
	printReleaseLevel = "[release] "
	printErrorLevel   = "[error  ] "
	printFatalLevel   = "[fatal  ] "
)

type Logger struct {
	level      int
	printlevel int
	baseLogger *log.Logger
	baseFile   *os.File
}

func New(logLevel string, printLevel string, pathname string, flag int, console bool) (*Logger, error) {
	// logLevel
	var ilogLevel int
	switch strings.ToLower(logLevel) {
	case "debug":
		ilogLevel = debugLevel
	case "release":
		ilogLevel = releaseLevel
	case "error":
		ilogLevel = errorLevel
	case "fatal":
		ilogLevel = fatalLevel
	default:
		return nil, errors.New("unknown level: " + logLevel)
	}

	// printLevel
	var iprintLevel int
	switch strings.ToLower(printLevel) {
	case "debug":
		iprintLevel = debugLevel
	case "release":
		iprintLevel = releaseLevel
	case "error":
		iprintLevel = errorLevel
	case "fatal":
		iprintLevel = fatalLevel
	default:
		return nil, errors.New("unknown level: " + logLevel)
	}

	// logger
	var baseLogger *log.Logger
	var baseFile *os.File
	if pathname != "" {
		now := time.Now()

		_, err := os.Stat(pathname)
		fmt.Println(pathname)
		if os.IsNotExist(err) {
			// 创建文件夹
			err := os.Mkdir(pathname, os.ModePerm)
			if err != nil {
				return nil, err
			}
		}

		pathname = fmt.Sprintf("%s\\%d%02d%02d",
			pathname,
			now.Year(),
			now.Month(),
			now.Day())

		filename := fmt.Sprintf("Userlogs-%02d%02d%02d.log",
			now.Hour(),
			now.Minute(),
			now.Second())

		_, err = os.Stat(pathname)
		fmt.Println(pathname)
		if os.IsNotExist(err) {
			// 创建文件夹
			err := os.Mkdir(pathname, os.ModePerm)
			if err != nil {
				return nil, err
			}
		}
		file, err := os.Create(path.Join(pathname, filename))
		if err != nil {
			return nil, err
		}

		baseLogger = log.New(file, "", flag)
		baseFile = file
	} else {
		baseLogger = log.New(os.Stdout, "", flag)
	}

	// new
	logger := new(Logger)
	logger.level = ilogLevel
	logger.printlevel = iprintLevel
	logger.baseLogger = baseLogger
	logger.baseFile = baseFile

	return logger, nil
}

// It's dangerous to call the method on logging
func (logger *Logger) Close() {
	if logger.baseFile != nil {
		logger.baseFile.Close()
	}

	logger.baseLogger = nil
	logger.baseFile = nil
}

func (logger *Logger) doPrintf(level int, printLevel string, format string, a ...interface{}) {
	if logger.baseLogger == nil {
		panic("logger closed")
	}

	if level < logger.level && level < logger.printlevel {
		return
	}

	format = printLevel + format
	if level >= logger.level {
		logger.baseLogger.Output(3, fmt.Sprintf(format, a...))
	}

	if level >= logger.printlevel {
		fmt.Sprintf(format, a...)
	}

	if level == fatalLevel {
		os.Exit(1)
	}
}

func (logger *Logger) Debug(format string, a ...interface{}) {
	logger.doPrintf(debugLevel, printDebugLevel, format, a...)
}

func (logger *Logger) Release(format string, a ...interface{}) {
	logger.doPrintf(releaseLevel, printReleaseLevel, format, a...)
}

func (logger *Logger) Error(format string, a ...interface{}) {
	logger.doPrintf(errorLevel, printErrorLevel, format, a...)
}

func (logger *Logger) Fatal(format string, a ...interface{}) {
	logger.doPrintf(fatalLevel, printFatalLevel, format, a...)
}

var gLogger, _ = New("debug", "debug", "", log.LstdFlags, true)

// It's dangerous to call the method on logging
func Export(logger *Logger) {
	if logger != nil {
		gLogger = logger
	}
}

func Debug(format string, a ...interface{}) {
	gLogger.doPrintf(debugLevel, printDebugLevel, format, a...)
}

func Release(format string, a ...interface{}) {
	gLogger.doPrintf(releaseLevel, printReleaseLevel, format, a...)
}

func Error(format string, a ...interface{}) {
	gLogger.doPrintf(errorLevel, printErrorLevel, format, a...)
}

func Fatal(format string, a ...interface{}) {
	gLogger.doPrintf(fatalLevel, printFatalLevel, format, a...)
}

func Close() {
	gLogger.Close()
}
