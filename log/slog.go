package log

import (
	"context"

	"github.com/jopbrown/gobase/errors"
	"golang.org/x/exp/slog"
)

type sLoggerHandler struct {
	l *Logger
	slog.Handler
}

func newSLoggerHandler(l *Logger, json bool) *sLoggerHandler {
	opts := slog.HandlerOptions{
		AddSource: l.format.AddSource,
		Level:     l.minLevel,
	}
	opts.ReplaceAttr = func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey {
			if !l.format.AddDateTime {
				return slog.Attr{}
			}
		}
		return a
	}
	var h slog.Handler
	if json {
		h = opts.NewJSONHandler(l.out)
	} else {
		h = opts.NewTextHandler(l.out)
	}

	attrs := make([]slog.Attr, 0, 2)
	if l.verbose > 0 {
		attrs = append(attrs, slog.Int("verbose", l.verbose))
	}

	if l.prefix != "" {
		attrs = append(attrs, slog.String("prefix", l.prefix))
	}
	if len(attrs) > 0 {
		h = h.WithAttrs(attrs)
	}

	sh := &sLoggerHandler{}
	sh.l = l
	sh.Handler = h

	return sh
}

func (h *sLoggerHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.l.enabled(Level(level))
}

func (h *sLoggerHandler) clone() *sLoggerHandler {
	newH := &sLoggerHandler{}
	newH.l = h.l
	newH.Handler = h.Handler
	return newH
}

func (h *sLoggerHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newH := h.clone()
	newH.Handler = h.Handler.WithAttrs(attrs)
	return newH
}

func (h *sLoggerHandler) WithGroup(name string) slog.Handler {
	newH := h.clone()
	newH.Handler = h.Handler.WithGroup(name)
	return newH
}

type sTeeLoggerHandler struct {
	tee *TeeLogger
	hs  []slog.Handler
}

func newSTeeLoggerHandler(tee *TeeLogger, json bool) *sTeeLoggerHandler {
	sh := &sTeeLoggerHandler{}
	sh.tee = tee
	sh.hs = make([]slog.Handler, 0, len(tee.loggers))

	for _, l := range tee.loggers {
		switch v := l.(type) {
		case *Logger:
			sh.hs = append(sh.hs, newSLoggerHandler(v, json))
		case *TeeLogger:
			sh.hs = append(sh.hs, newSTeeLoggerHandler(v, json))
		}

	}

	return sh
}

func (th *sTeeLoggerHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return th.tee.enabled(Level(level))
}

func (th *sTeeLoggerHandler) Handle(ctx context.Context, r slog.Record) error {
	var err error
	for _, h := range th.hs {
		if !h.Enabled(ctx, r.Level) {
			continue
		}

		err = errors.Join(err, h.Handle(ctx, r))
	}

	return err
}

func (th *sTeeLoggerHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newH := &sTeeLoggerHandler{}
	newH.tee = th.tee
	newH.hs = make([]slog.Handler, 0, len(th.hs))
	for _, h := range th.hs {
		newH.hs = append(newH.hs, h.WithAttrs(attrs))
	}

	return newH
}

func (th *sTeeLoggerHandler) WithGroup(name string) slog.Handler {
	newH := &sTeeLoggerHandler{}
	newH.tee = th.tee
	newH.hs = make([]slog.Handler, 0, len(th.hs))
	for _, h := range th.hs {
		newH.hs = append(newH.hs, h.WithGroup(name))
	}

	return newH
}
