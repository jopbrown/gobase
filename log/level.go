package log

import (
	"strconv"

	"log/slog"

	"github.com/jopbrown/gobase/errors"
)

type Level slog.Level

const (
	LevelDebug = Level(slog.LevelDebug)
	LevelInfo  = Level(slog.LevelInfo)
	LevelWarn  = Level(slog.LevelWarn)
	LevelError = Level(slog.LevelError)
	LevelPanic = Level(slog.LevelError + 1)
	LevelFatal = Level(slog.LevelError + 2)

	LevelNone = Level(-99)
	LevelAll  = Level(99)
)

var level2Name = map[Level]string{
	LevelNone:  "NONE",
	LevelDebug: "DEBUG",
	LevelInfo:  "INFO",
	LevelWarn:  "WARN",
	LevelError: "ERROR",
	LevelPanic: "PANIC",
	LevelFatal: "FATAL",
	LevelAll:   "ALL",
}

var name2Level = map[string]Level{
	"NONE":  LevelNone,
	"DEBUG": LevelDebug,
	"INFO":  LevelInfo,
	"WARN":  LevelWarn,
	"ERROR": LevelError,
	"PANIC": LevelPanic,
	"FATAL": LevelFatal,
	"ALL":   LevelAll,
}

func (l Level) Level() slog.Level { return slog.Level(l) }

func (l Level) String() string {
	return level2Name[l]
}

func (l *Level) parse(s string) (err error) {
	ll, ok := name2Level[s]
	if !ok {
		return errors.Errorf("unable to parse level string: %q", s)
	}

	*l = ll
	return nil
}

func (l Level) MarshalJSON() ([]byte, error) {
	return strconv.AppendQuote(nil, l.String()), nil
}

func (l *Level) UnmarshalJSON(data []byte) error {
	s, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}
	return l.parse(s)
}

func (l Level) MarshalText() ([]byte, error) {
	return []byte(l.String()), nil
}

func (l *Level) UnmarshalText(data []byte) error {
	return l.parse(string(data))
}
