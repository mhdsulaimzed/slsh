package utils

import (
	"fmt"
	"strings"
	"time"
)

// Colors for terminal output
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
	ColorBold   = "\033[1m"
)

// FormatJobState colorizes job states
func FormatJobState(state string, useColor bool) string {
	if !useColor {
		return state
	}
	
	switch strings.ToUpper(state) {
	case "RUNNING":
		return ColorGreen + state + ColorReset
	case "PENDING":
		return ColorYellow + state + ColorReset
	case "COMPLETED":
		return ColorBlue + state + ColorReset
	case "FAILED", "CANCELLED", "TIMEOUT":
		return ColorRed + state + ColorReset
	default:
		return state
	}
}

// FormatNodeState colorizes node states
func FormatNodeState(state string, useColor bool) string {
	if !useColor {
		return state
	}
	
	switch strings.ToUpper(state) {
	case "IDLE":
		return ColorGreen + state + ColorReset
	case "ALLOC", "MIXED":
		return ColorYellow + state + ColorReset
	case "DOWN", "DRAIN", "FAIL":
		return ColorRed + state + ColorReset
	default:
		return state
	}
}

// FormatDuration formats a duration in a human-readable way
func FormatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%.1fs", d.Seconds())
	} else if d < time.Hour {
		return fmt.Sprintf("%.1fm", d.Minutes())
	} else if d < 24*time.Hour {
		return fmt.Sprintf("%.1fh", d.Hours())
	} else {
		days := int(d.Hours() / 24)
		hours := d.Hours() - float64(days*24)
		return fmt.Sprintf("%dd%.1fh", days, hours)
	}
}

// FormatMemory formats memory sizes
func FormatMemory(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// Table represents a simple table for output formatting
type Table struct {
	Headers []string
	Rows    [][]string
	useColor bool
}

// NewTable creates a new table
func NewTable(headers []string, useColor bool) *Table {
	return &Table{
		Headers:  headers,
		Rows:     make([][]string, 0),
		useColor: useColor,
	}
}

// AddRow adds a row to the table
func (t *Table) AddRow(row []string) {
	t.Rows = append(t.Rows, row)
}

// Print prints the table
func (t *Table) Print() {
	if len(t.Headers) == 0 {
		return
	}
	
	// Calculate column widths
	widths := make([]int, len(t.Headers))
	for i, header := range t.Headers {
		widths[i] = len(header)
	}
	
	for _, row := range t.Rows {
		for i, cell := range row {
			if i < len(widths) {
				cellLen := len(stripAnsiCodes(cell))
				if cellLen > widths[i] {
					widths[i] = cellLen
				}
			}
		}
	}
	
	// Print header
	if t.useColor {
		fmt.Print(ColorBold)
	}
	for i, header := range t.Headers {
		fmt.Printf("%-*s", widths[i]+2, header)
	}
	if t.useColor {
		fmt.Print(ColorReset)
	}
	fmt.Println()
	
	// Print separator
	for i := range t.Headers {
		fmt.Print(strings.Repeat("-", widths[i]+2))
	}
	fmt.Println()
	
	// Print rows
	for _, row := range t.Rows {
		for i, cell := range row {
			if i < len(widths) {
				// Account for ANSI color codes when padding
				padding := widths[i] + 2 - (len(cell) - len(stripAnsiCodes(cell)))
				fmt.Printf("%-*s", padding, cell)
			}
		}
		fmt.Println()
	}
}

// stripAnsiCodes removes ANSI color codes from a string for length calculation
func stripAnsiCodes(s string) string {
	// Simple regex replacement would be better, but avoiding external deps
	result := ""
	inEscape := false
	
	for _, r := range s {
		if r == '\033' {
			inEscape = true
		} else if inEscape && r == 'm' {
			inEscape = false
		} else if !inEscape {
			result += string(r)
		}
	}
	
	return result
}

// FormatSuccess formats success/error messages
func FormatSuccess(msg string, useColor bool) string {
	if useColor {
		return ColorGreen + "✓ " + msg + ColorReset
	}
	return "✓ " + msg
}

// FormatError formats error messages
func FormatError(msg string, useColor bool) string {
	if useColor {
		return ColorRed + "✗ " + msg + ColorReset
	}
	return "✗ " + msg
}

// FormatWarning formats warning messages
func FormatWarning(msg string, useColor bool) string {
	if useColor {
		return ColorYellow + "⚠ " + msg + ColorReset
	}
	return "⚠ " + msg
}

// FormatInfo formats info messages
func FormatInfo(msg string, useColor bool) string {
	if useColor {
		return ColorBlue + "ℹ " + msg + ColorReset
	}
	return "ℹ " + msg
}