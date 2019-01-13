package logger

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

const (
	// AlwaysLabel displays an airplane on "always" logs.
	AlwaysLabel = "✈"
	// CriticalLabel displays an X on "critical" logs.
	CriticalLabel = "✖"
	// DebugLabel displays a play button on "debug" logs.
	DebugLabel = "▶"
	// InfoLabel displays an i icon on "info" logs.
	InfoLabel = "ℹ"
	// SuccessLabel displays a check mark on "success" logs.
	SuccessLabel = "✔"
	// WarningLabel displays an ! on "warning" logs.
	WarningLabel = "!"
	// AskLabel displays a ? on "ask" logs.
	AskLabel = "?"
)

var (
	// Level defines the default log level.
	Level = 3
	// Color toggles output colorization.
	Color = true
	// Timestamps toggles timestamps on output logs.
	Timestamps = false
)

// Log will print a formatted generic statement to standard output.
func Log(format string, a ...interface{}) {
	a, w := extractLoggerArgs(format, a...)
	fmt.Fprintf(w, label(format, ""), a...)
}

// Table will print headers and data in a pretty formatted ASCII table.
func Table(headers []string, data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)
	for _, v := range data {
		table.Append(v)
	}
	table.Render()
}

// Ask will always print a formatted generic statement to standard output
// formatted to be a question.
func Ask(format string, a ...interface{}) {
	a, w := extractLoggerArgs(format, a...)
	s := fmt.Sprintf(label(format, AskLabel), a...)

	if Color {
		w = color.Output
		s = color.YellowString(s)
	}

	fmt.Fprintf(w, s)
}

// Always will "always" print a formatted generic statement to standard output.
func Always(format string, a ...interface{}) {
	a, w := extractLoggerArgs(format, a...)
	s := fmt.Sprintf(label(format, AlwaysLabel), a...)

	if Color {
		w = color.Output
		s = color.BlueString(s)
	}

	fmt.Fprintf(w, s)
}

// Critical will print a formatted generic statement to standard output if the
// global level is set to 1 or higher. The print statement will have a color of
// red if color is enabled.
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

// Info will print a formatted generic statement to standard output if the
// global level is set to 3 or higher. The print statement will have a color of
// cyan if color is enabled.
func Info(format string, a ...interface{}) {
	if Level >= 3 {
		a, w := extractLoggerArgs(format, a...)
		s := fmt.Sprintf(label(format, InfoLabel), a...)

		if Color {
			w = color.Output
			s = color.MagentaString(s)
		}

		fmt.Fprintf(w, s)
	}
}

// Success will print a formatted generic statement to standard output if the
// global level is set to 3 or higher. The print statement will have a color of
// cyan if color is enabled.
func Success(format string, a ...interface{}) {
	if Level >= 3 {
		a, w := extractLoggerArgs(format, a...)
		s := fmt.Sprintf(label(format, SuccessLabel), a...)

		if Color {
			w = color.Output
			s = color.GreenString(s)
		}

		fmt.Fprintf(w, s)
	}
}

// Debug will print a formatted generic statement to standard output if the
// global level is set to 4 or higher.
func Debug(format string, a ...interface{}) {
	if Level >= 4 {
		a, w := extractLoggerArgs(format, a...)
		s := fmt.Sprintf(label(format, DebugLabel), a...)

		fmt.Fprintf(w, s)
	}
}

// Warning will print a formatted generic statement to standard output if the
// global level is set to 2 or higher. The print statement will have a color of
// yellow if color is enabled.
func Warning(format string, a ...interface{}) {
	if Level >= 2 {
		a, w := extractLoggerArgs(format, a...)
		s := fmt.Sprintf(label(format, WarningLabel), a...)

		if Color {
			w = color.Output
			s = color.YellowString(s)
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
	if !strings.Contains(format, "\n") {
		format = fmt.Sprintf("%s%s", format, "\n")
	}
	if label != "" {
		format = fmt.Sprintf("[%s]  %s", label, format)
	}
	if Timestamps {
		t := time.Now()
		rfct := t.Format(time.RFC3339)
		format = fmt.Sprintf("%s %s", rfct, format)
	}
	return format
}
