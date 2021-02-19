package log

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	Time         = 1 << iota
	DirName
	FileName
	Line
	FuncName
)


func (l *Logger) outPut(msg string, level string) {
	var m = make(map[string]string, 5)

	if (l.flag & Time) != 0 {
		m["Time"] = "[" + time.Now().Format("2006-01-02 15:04:05") + "]"
	}
	if (l.flag & DirName) != 0 || (l.flag & FileName) != 0 ||
		(l.flag & Line) != 0 || (l.flag & FuncName) != 0 {
		funcName, dir, file, line := l.getInfo(3)	// just in sync mode
		if (l.flag & DirName) != 0 {
			m["DirName"] = dir
		}
		if (l.flag & FileName) != 0 {
			m["FileName"] = file
		}
		if (l.flag & Line) != 0 {
			m["Line"] = line
		}
		if (l.flag & FuncName) != 0 {
			m["FuncName"] = funcName
		}
	}
	if len(l.prefix) != 0 {
		m["MsgPrefix"] = l.prefix
	}
	m["msg"] = msg
	m["level"] = level

	// output json
	var content string
	if strings.ToLower(l.format) == "json" {
		jsonBytes, _ := json.Marshal(m)
		content = string(jsonBytes)
	}else {
		// output string
		content = fmt.Sprintf("%s %s %s %s %s %s %s %s",
			m["Time"], m["level"], m["DirName"], m["FileName"], m["Line"], m["FuncName"], m["MsgPrefix"], m["msg"])
	}

	if l.syn {
		_, _ = fmt.Fprintf(l.out, "%s\n", content)
	} else {
		// Need to lock when writing synchronously
		l.mu.Lock()
		_, _ = fmt.Fprintf(l.out, "%s\n", content)
		l.mu.Unlock()
	}
	f, ok := l.out.(*os.File)
	if ok {
		l.mu.Lock()
		l.out, _ = checkAndCut(f, l.maxSize)
		l.mu.Unlock()
	}
}

// get caller info
func (l *Logger) getInfo(skip int) (funcName string, dir string, file string, line string) {
	caller, file, iLine, ok := runtime.Caller(skip)
	if !ok {
		file = "unknown"
		line = "0"
	}
	funcName = runtime.FuncForPC(caller).Name()
	dir = path.Dir(file)
	file = path.Base(file)
	return funcName, dir, file, strconv.Itoa(iLine)
}

// log file segment
func checkAndCut(f *os.File, byteSize int64) (*os.File, error) {
	info, err := f.Stat()
	if err != nil {
		return f, err
	}

	if info.Size() < byteSize {
		return f, nil
	}

	longFileName := f.Name()
	dirName, fileName := path.Split(longFileName)
	suffix := strings.TrimPrefix(path.Ext(fileName), ".")
	fileNameWithoutExt := strings.TrimSuffix(fileName, path.Ext(fileName))
	newName := path.Join(dirName,
		strings.Join([]string{fileNameWithoutExt, time.Now().Format("20060102150405"), suffix}, "."))

	err = f.Close()
	if err != nil {
		return f, err
	}

	err = os.Rename(longFileName, newName)
	if err != nil {
		return f, err
	}

	f, err = os.Create(longFileName)
	if err != nil {
		return f, err
	}

	return f, err
}


