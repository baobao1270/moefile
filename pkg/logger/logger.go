package logger

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

const (
	CRed     = LogColor("\033[01;31m")
	CGreen   = LogColor("\033[01;32m")
	CYellow  = LogColor("\033[01;33m")
	CBlue    = LogColor("\033[01;34m")
	CMagenta = LogColor("\033[01;35m")
	CCyan    = LogColor("\033[01;36m")
	CWhite   = LogColor("\033[01;37m")
	CReset   = LogColor("\033[0m")

	LDbg = LogLevel(0x00)
	LInf = LogLevel(0x01)
	LWrn = LogLevel(0x02)
	LErr = LogLevel(0x03)
)

var (
	LevelNameMap = map[LogLevel]string{
		LDbg: "DBG",
		LInf: "INF",
		LWrn: "WRN",
		LErr: "ERR",
	}

	LevelColorMap = map[LogLevel]LogColor{
		LDbg: CBlue,
		LInf: CGreen,
		LWrn: CYellow,
		LErr: CRed,
	}
)

type Logger struct {
	Writer     io.Writer
	TimeFormat string
	TagColor   map[string]LogColor
	MinLevel   LogLevel
}

type Tag struct {
	Name   string
	Color  LogColor
	Logger *Logger
}

type LogColor string
type LogLevel uint8

func New(w io.Writer) *Logger {
	return &Logger{
		Writer:     w,
		TimeFormat: time.RFC3339,
	}
}

func NewStdout() *Logger {
	l := New(os.Stdout)
	l.MinLevel = LInf
	return l
}

func (l *Logger) AddTagColor(prefix string, color LogColor) *Logger {
	l.TagColor[prefix] = color
	return l
}

func (l *Logger) SetTagColor(m map[string]LogColor) *Logger {
	l.TagColor = m
	return l
}

func (l *Logger) SetTimeFormat(format string) *Logger {
	l.TimeFormat = format
	return l
}

func (l *Logger) GetFormatedTime() string {
	return time.Now().Local().Format(l.TimeFormat)
}

func (l *Logger) Logf(tag string, level LogLevel, format string, a ...any) {
	l.Tag(tag).Logf(level, format, a...)
}

func (l *Logger) Tag(tag string) *Tag {
	color := CReset
	for prefix, c := range l.TagColor {
		if strings.HasPrefix(tag, prefix) {
			color = c
		}
	}
	return &Tag{Name: tag, Color: color, Logger: l}
}

func (tag *Tag) ColoredName() string {
	return fmt.Sprintf("%s%s%s", tag.Color, tag.Name, CReset)
}

func (tag *Tag) Logf(level LogLevel, format string, a ...any) {
	if level < tag.Logger.MinLevel {
		return
	}
	lines := strings.Split(fmt.Sprintf(format, a...), "\n")
	for _, line := range lines {
		c := fmt.Sprintf("%s %s - [%s] %s\n",
			tag.Logger.GetFormatedTime(),
			level.Colored(),
			tag.ColoredName(),
			line,
		)
		//nolint:errcheck
		tag.Logger.Writer.Write([]byte(c))
	}
}

func (tag *Tag) LogWriter(level LogLevel) io.Writer {
	r, w := io.Pipe()
	go func() {
		defer w.Close()
		buf2 := []byte{}
		buf1 := make([]byte, 1)
		for {
			_, err := r.Read(buf1)
			if err != nil {
				tag.Logger.Tag("log").Errf("LogWriter read error: %v", err)
				break
			}
			if buf1[0] == '\n' {
				tag.Logf(level, "%s", buf2)
				buf2 = []byte{}
			} else {
				buf2 = append(buf2, buf1[0])
			}
		}
	}()
	return w
}

func (tag *Tag) Dbgf(format string, a ...any) {
	tag.Logf(LDbg, format, a...)
}

func (tag *Tag) Inff(format string, a ...any) {
	tag.Logf(LInf, format, a...)
}

func (tag *Tag) Wrnf(format string, a ...any) {
	tag.Logf(LWrn, format, a...)
}

func (tag *Tag) Errf(format string, a ...any) {
	tag.Logf(LErr, format, a...)
}

func (ll *LogLevel) Colored() string {
	return fmt.Sprintf("%s%s%s", LevelColorMap[*ll], LevelNameMap[*ll], CReset)
}
