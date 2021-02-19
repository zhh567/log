package log

import (
	"context"
	"io"
	"os"
	"sync"
)

type Level uint16

// log level
const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
	FATAL
	OFF
)

type Logger struct {
	level  Level
	mu     sync.Mutex	// writ file concurrently
	syn	   bool
	prefix string
	out    io.Writer
	format string
	flag   int
	maxSize   int64
}

type msgOrder struct {
	msg	string
	order string
}

var ch chan msgOrder

func NewLogger(level Level) *Logger {
	return &Logger{
		level: level,
		out: os.Stdout,
		flag: Time,
		syn: false,
		maxSize: 1024 * 1024 * 10,
	}
}

// The log is written asynchronously in this way,
// but the calling line cannot be determined.
func NewLoggerSync(ctx context.Context, level Level) *Logger {
	l := &Logger{
		level: level,
		out: os.Stdout,
		flag: Time,
		syn: true,
		maxSize: 1024 * 1024 * 10,
	}
	ch = make(chan msgOrder, 20)
	go worker(ctx , l, ch)
	return l
}


func (l *Logger) SetPrefix(prefix string)  {
	l.prefix = prefix
}

func (l *Logger) SetJsonFormat()  {
	l.format = "json"
}

func (l *Logger) SetFlag(flag int)  {
	l.flag = flag
}

func (l *Logger) SetOutput(output io.Writer)  {
	l.out = output
}

func (l *Logger) SetMaxSize(maxSize int64)  {
	l.maxSize = maxSize
}

func (l *Logger) Debug(msg string) {
	if l.level <= DEBUG {
		if l.syn {
			ch <- msgOrder{msg: msg, order: "[DEBUG]"}
		} else {
			l.outPut(msg, "[DEBUG]")
		}
	}
}

func (l *Logger) Info(msg string) {
	if l.level <= INFO {
		if l.syn {
			ch <- msgOrder{msg: msg, order: "[INFO]"}
		} else {
			l.outPut(msg, "[INFO]")
		}
	}
}

func (l *Logger) Warn(msg string) {
	if l.level <= WARN {
		if l.syn {
			ch <- msgOrder{msg: msg, order: "[WARN]"}
		} else {
			l.outPut(msg, "[WARN]")
		}
	}
}

func (l *Logger) Error(msg string)  {
	if l.level <= ERROR {
		if l.syn {
			ch <- msgOrder{msg: msg, order: "[ERROR]"}
		} else {
			l.outPut(msg, "[ERROR]")
		}
	}
}

func (l *Logger) FATAL(msg string)  {
	if l.level <= FATAL {
		if l.syn {
			ch <- msgOrder{msg: msg, order: "[FATAL]"}
		} else {
			l.outPut(msg, "[FATAL]")
		}
	}
}
