package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

const (
	FormatJSON   = "json"
	FormatText   = "text"
	FormatCustom = "custom"
)

func SetupLogger(format string, opts *slog.HandlerOptions) {
	var sHandler slog.Handler
	switch format {
	case FormatJSON:
		sHandler = slog.NewJSONHandler(os.Stdout, opts)
	case FormatText:
		sHandler = slog.NewTextHandler(os.Stdout, opts)
	case FormatCustom:
		sHandler = NewCustomHandler(os.Stdout, opts)
	default:
		slog.Error("Unknown log format", slog.String("format", format))
		os.Exit(-1)
	}

	slog.SetDefault(slog.New(sHandler))
}

// ANSI colors for log levels
var levelColors = map[slog.Level]string{
	slog.LevelDebug: "\033[36m", // cyan
	slog.LevelInfo:  "\033[32m", // green
	slog.LevelWarn:  "\033[33m", // yellow
	slog.LevelError: "\033[31m", // red
}

const reset = "\033[0m"
const grey = "\033[30m"
const bold = "\033[1m"

type CustomHandler struct {
	w      *os.File
	opts   *slog.HandlerOptions
	attrs  []slog.Attr
	groups []string
}

func NewCustomHandler(w *os.File, opts *slog.HandlerOptions) *CustomHandler {
	return &CustomHandler{w: w, opts: opts, attrs: []slog.Attr{}, groups: []string{}}
}

func (h *CustomHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.opts.Level.Level()
}

func (h *CustomHandler) Handle(_ context.Context, r slog.Record) error {
	timestamp := r.Time.Format("15:04:05")
	color := levelColors[r.Level]
	level := r.Level.String()

	src := ""
	if r.PC != 0 && h.opts.AddSource {
		fs := r.Source()
		src += " " + filepath.Base(fs.File)
	}

	allAttrs := []slog.Attr{}
	allAttrs = append(allAttrs, h.attrs...)
	r.Attrs(func(a slog.Attr) bool {
		allAttrs = append(allAttrs, a)
		return true
	})

	line := fmt.Sprintf("%s%s%s %s%s%s %s%s\n%s", grey, timestamp, color, level, src, reset, bold, r.Message, reset)

	// Render attrs in YAML-like format under the line
	if len(allAttrs) > 0 {
		for _, a := range allAttrs {
			key := strings.Join(append(h.groups, a.Key), ".")
			line += fmt.Sprintf("%s%s %s: %s%s%v\n", bold, grey, key, reset, grey, a.Value)
		}
	}

	_, err := h.w.Write([]byte(line))
	return err
}

func (h *CustomHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newH := *h
	newH.attrs = append(append([]slog.Attr{}, h.attrs...), attrs...)
	return &newH
}

func (h *CustomHandler) WithGroup(name string) slog.Handler {
	newH := *h
	newH.groups = append(append([]string{}, h.groups...), name)
	return &newH
}
