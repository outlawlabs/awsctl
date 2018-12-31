package logger

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
)

type Logger func(format string, a ...interface{})

const (
	AlwaysLabel   = "✈"
	CriticalLabel = "✖"
	DebugLabel    = "▶"
	InfoLabel     = "ℹ"
	SuccessLabel  = "✔"
	WarningLabel  = "!"
)

var (
	Level = 3
	Color = true
)

func Log(format string, a ...interface{}) {
	a, w := extractLoggerArgs(format, a...)
	fmt.Fprintf(w, format, a...)
}

func Always(format string, a ...interface{}) {
	a, w := extractLoggerArgs(format, a...)
	s := fmt.Sprintf(label(format, AlwaysLabel), a...)

	if Color {
		w = color.Output
		s = color.GreenString(s)
	}

	fmt.Fprintf(w, s)
}

func Critical(format string, a ...interface{}) {
	if Level >= 1 {
		a, w := extractLoggerArgs(format, a...)
		s := fmt.Sprintf(label(format, CriticalLabel), a...)

		if Color {
			w = color.Output
			s = color.RedString(s)
		}

		fmt.Fprintf(w, s)
	}
}

func Info(format string, a ...interface{}) {
	if Level >= 3 {
		a, w := extractLoggerArgs(format, a...)
		s := fmt.Sprintf(label(format, InfoLabel), a...)

		if Color {
			w = color.Output
			s = color.CyanString(s)
		}

		fmt.Fprintf(w, s)
	}
}

func Success(format string, a ...interface{}) {
	if Level >= 3 {
		a, w := extractLoggerArgs(format, a...)
		s := fmt.Sprintf(label(format, SuccessLabel), a...)

		if Color {
			w = color.Output
			s = color.CyanString(s)
		}

		fmt.Fprintf(w, s)
	}
}

func Debug(format string, a ...interface{}) {
	if Level >= 4 {
		a, w := extractLoggerArgs(format, a...)
		s := fmt.Sprintf(label(format, DebugLabel), a...)

		fmt.Fprintf(w, s)
	}
}

func Warning(format string, a ...interface{}) {
	if Level >= 2 {
		a, w := extractLoggerArgs(format, a...)
		s := fmt.Sprintf(label(format, WarningLabel), a...)

		if Color {
			w = color.Output
			s = color.GreenString(s)
		}

		fmt.Fprintf(w, s)
	}
}

func extractLoggerArgs(format string, a ...interface{}) ([]interface{}, io.Writer) {
	var w io.Writer = os.Stdout

	if n := len(a); n > 0 {
		// extract an io.Writer at the end of a
		if value, ok := a[n-1].(io.Writer); ok {
			w = value
			a = a[0 : n-1]
		}
	}

	return a, w
}

func label(format, label string) string {
	t := time.Now()
	rfct := t.Format(time.RFC3339)
	if !strings.Contains(format, "\n") {
		format = fmt.Sprintf("%s%s", format, "\n")
	}
	return fmt.Sprintf("%s [%s]  %s", rfct, label, format)
}
