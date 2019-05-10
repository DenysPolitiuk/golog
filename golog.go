package golog

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	fileExtension   = ".log"
	combineFilename = "combine"
	timeLayout      = "02-01-2006 15:04:05.0000"
	dateRune        = "@DATE"
	severityRune    = "@SEVERITY"
	msgRune         = "@MSG"
)

var (
	Location      = "."
	Application   = ""
	IsMultiLog    = true
	MessageFormat = "[" + dateRune + "][" + severityRune + "] : " + msgRune
)

type CustomError string

func (c CustomError) Error() string {
	return string(c)
}

type Severity string

const (
	ERROR   Severity = "error"
	DEBUG   Severity = "debug"
	INFO    Severity = "info"
	combine Severity = combineFilename
)

func (s Severity) String() string {
	return string(s)
}

func LogAny(msg, severity string) (string, error) {
	return Log(msg, Severity(severity))
}

func Log(msg string, severity Severity) (string, error) {
	timestamp := fmt.Sprint(time.Now().Format(timeLayout))
	entry := strings.Replace(MessageFormat, dateRune, timestamp, -1)
	entry = strings.Replace(entry, severityRune, strings.ToUpper(string(severity)), -1)
	entry = strings.Replace(entry, msgRune, msg, -1)
	var err error
	if IsMultiLog {
		severityFile, err := fileSetup(Location, makeFileName(string(severity)))
		if err != nil {
			return "", err
		}
		err = writeToFile(severityFile, entry)
		if err != nil {
			return "", err
		}
	}
	combineFile, err := fileSetup(Location, makeFileName(string(combine)))
	if err != nil {
		return "", err
	}
	err = writeToFile(combineFile, entry)
	if err != nil {
		return "", err
	}
	return entry, nil
}

func writeToFile(filename string, msg string) error {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.WriteString(msg + "\n"); err != nil {
		return err
	}
	return nil
}

func makeFileName(severity string) string {
	var name string
	if Application != "" {
		name = fmt.Sprintf("%v_", Application)
	}
	return fmt.Sprintf("%v%v%v", name, severity, fileExtension)
}

func fileSetup(folder, file string) (string, error) {
	fol, err := os.Stat(folder)
	if err != nil {
		return "", err
	}
	if !fol.IsDir() {
		return "", CustomError("fileSetup target must be a directory")
	}
	fullname := filepath.Join(folder, file)
	// check if file already exists
	if _, err := os.Stat(fullname); err == nil {
		return fullname, nil
	}
	f, err := os.Create(fullname)
	if err != nil {
		return "", err
	}
	defer f.Close()
	return fullname, nil
}
