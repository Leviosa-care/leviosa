package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"reflect"
	"strings"
	"sync"
)

type DevHandler struct {
	level slog.Leveler
	group string
	attrs []slog.Attr
	mu    *sync.Mutex
	w     io.Writer
}

func NewDevHandler(w io.Writer, level slog.Leveler) *DevHandler {
	if level == nil || reflect.TypeOf(level).Kind() == reflect.Ptr && reflect.ValueOf(level).IsNil() {
		level = &slog.LevelVar{}
	}
	return &DevHandler{
		level: level,
		mu:    new(sync.Mutex),
		w:     w,
	}
}

func (h *DevHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.level.Level()
}

func (h *DevHandler) Handle(ctx context.Context, r slog.Record) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Use strings.Builder for efficient string building
	var attrs strings.Builder

	// Format handler attributes
	for _, a := range h.attrs {
		if !a.Equal(slog.Attr{}) {
			attrs.WriteString(" ")
			if h.group != "" {
				attrs.WriteString(h.group)
				attrs.WriteString(".")
			}
			attrs.WriteString(a.Key)
			attrs.WriteString(": ")
			attrs.WriteString(a.Value.String())
			attrs.WriteString("\n")
		}
	}

	// Format record attributes
	r.Attrs(func(a slog.Attr) bool {
		if !a.Equal(slog.Attr{}) {
			attrs.WriteString(" ")
			if h.group != "" {
				attrs.WriteString(h.group)
				attrs.WriteString(".")
			}
			attrs.WriteString(a.Key)
			attrs.WriteString(": ")
			attrs.WriteString(a.Value.String())
			attrs.WriteString("\n")
		}
		return true
	})

	// Get level indicator and color
	levelIndicator := getLevelIndicator(r.Level)
	color := getLevelColor(r.Level)
	resetColor := "\033[0m"

	// Build final output
	attrsStr := attrs.String()
	attrsStr = strings.TrimRight(attrsStr, "\n")

	var newlines string
	if attrsStr != "" {
		newlines = "\n\n"
	}

	fmt.Fprintf(h.w, "%s[%s] %s %s%s\n%s%s",
		color,
		r.Time.Format("15:04:05"),
		levelIndicator,
		r.Message,
		resetColor,
		attrsStr,
		newlines)

	return nil
}

// getLevelIndicator returns visual indicator for log level
func getLevelIndicator(level slog.Level) string {
	switch {
	case level >= slog.LevelError:
		return "🔴 ERROR"
	case level >= slog.LevelWarn:
		return "⚠️  WARN "
	case level >= slog.LevelInfo:
		return "🔵 INFO "
	default:
		return "🔍 DEBUG"
	}
}

// getLevelColor returns ANSI color code for log level
func getLevelColor(level slog.Level) string {
	switch {
	case level >= slog.LevelError:
		return "\033[31m" // Red
	case level >= slog.LevelWarn:
		return "\033[33m" // Yellow
	case level >= slog.LevelInfo:
		return "\033[36m" // Cyan
	default:
		return "\033[37m" // Light gray
	}
}

func (h *DevHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &DevHandler{
		level: h.level,
		group: h.group,
		attrs: append(h.attrs, attrs...),
		mu:    h.mu,
		w:     h.w,
	}
}

func (h *DevHandler) WithGroup(name string) slog.Handler {
	return &DevHandler{
		level: h.level,
		group: strings.TrimSuffix(name+"."+h.group, "."),
		attrs: h.attrs,
		mu:    h.mu,
		w:     h.w,
	}
}
