package shell

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// HistoryEntry represents a single history entry
type HistoryEntry struct {
	Command   string    `json:"command"`
	Timestamp time.Time `json:"timestamp"`
	Success   bool      `json:"success"`
	Duration  time.Duration `json:"duration"`
}

// History manages command history
type History struct {
	entries  []HistoryEntry
	maxSize  int
	filePath string
}

// NewHistory creates a new history manager
func NewHistory(maxSize int) *History {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "/tmp"
	}
	
	return &History{
		entries:  make([]HistoryEntry, 0),
		maxSize:  maxSize,
		filePath: filepath.Join(homeDir, ".slsh_history"),
	}
}

// Add adds a command to history
func (h *History) Add(command string, success bool, duration time.Duration) {
	entry := HistoryEntry{
		Command:   strings.TrimSpace(command),
		Timestamp: time.Now(),
		Success:   success,
		Duration:  duration,
	}

	// Skip empty commands and duplicates
	if entry.Command == "" {
		return
	}
	
	// Skip if same as last command
	if len(h.entries) > 0 && h.entries[len(h.entries)-1].Command == entry.Command {
		return
	}

	h.entries = append(h.entries, entry)

	// Maintain max size
	if len(h.entries) > h.maxSize {
		h.entries = h.entries[len(h.entries)-h.maxSize:]
	}
}

// GetAll returns all history entries
func (h *History) GetAll() []HistoryEntry {
	return h.entries
}

// GetLast returns the last n entries
func (h *History) GetLast(n int) []HistoryEntry {
	if n <= 0 || len(h.entries) == 0 {
		return []HistoryEntry{}
	}
	
	start := len(h.entries) - n
	if start < 0 {
		start = 0
	}
	
	return h.entries[start:]
}

// Search searches for commands containing the given string
func (h *History) Search(query string) []HistoryEntry {
	var results []HistoryEntry
	query = strings.ToLower(query)
	
	for _, entry := range h.entries {
		if strings.Contains(strings.ToLower(entry.Command), query) {
			results = append(results, entry)
		}
	}
	
	return results
}

// GetByIndex returns a command by its index (1-based)
func (h *History) GetByIndex(index int) (string, error) {
	if index < 1 || index > len(h.entries) {
		return "", fmt.Errorf("history index %d out of range (1-%d)", index, len(h.entries))
	}
	
	return h.entries[index-1].Command, nil
}

// Clear clears all history
func (h *History) Clear() {
	h.entries = make([]HistoryEntry, 0)
}

// Save saves history to file
func (h *History) Save() error {
	file, err := os.Create(h.filePath)
	if err != nil {
		return fmt.Errorf("failed to create history file: %v", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	for _, entry := range h.entries {
		// Format: timestamp|success|duration|command
		line := fmt.Sprintf("%d|%t|%d|%s\n", 
			entry.Timestamp.Unix(), 
			entry.Success, 
			entry.Duration.Nanoseconds(),
			entry.Command)
		
		if _, err := writer.WriteString(line); err != nil {
			return fmt.Errorf("failed to write history entry: %v", err)
		}
	}

	return nil
}

// Load loads history from file
func (h *History) Load() error {
	file, err := os.Open(h.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // No history file exists, that's fine
		}
		return fmt.Errorf("failed to open history file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	entries := make([]HistoryEntry, 0)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		entry, err := parseHistoryLine(line)
		if err != nil {
			continue // Skip invalid lines
		}

		entries = append(entries, entry)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to read history file: %v", err)
	}

	// Keep only the last maxSize entries
	if len(entries) > h.maxSize {
		entries = entries[len(entries)-h.maxSize:]
	}

	h.entries = entries
	return nil
}

// parseHistoryLine parses a history line from the file
func parseHistoryLine(line string) (HistoryEntry, error) {
	parts := strings.SplitN(line, "|", 4)
	if len(parts) != 4 {
		return HistoryEntry{}, fmt.Errorf("invalid history line format")
	}

	var entry HistoryEntry
	var err error

	// Parse timestamp
	var timestamp int64
	if _, err = fmt.Sscanf(parts[0], "%d", &timestamp); err != nil {
		return HistoryEntry{}, fmt.Errorf("invalid timestamp: %v", err)
	}
	entry.Timestamp = time.Unix(timestamp, 0)

	// Parse success
	if _, err = fmt.Sscanf(parts[1], "%t", &entry.Success); err != nil {
		return HistoryEntry{}, fmt.Errorf("invalid success flag: %v", err)
	}

	// Parse duration
	var duration int64
	if _, err = fmt.Sscanf(parts[2], "%d", &duration); err != nil {
		return HistoryEntry{}, fmt.Errorf("invalid duration: %v", err)
	}
	entry.Duration = time.Duration(duration)

	// Command is the rest
	entry.Command = parts[3]

	return entry, nil
}

// PrintHistory prints history entries in a formatted way
func (h *History) PrintHistory(showTimestamp bool, showDuration bool) {
	if len(h.entries) == 0 {
		fmt.Println("No history entries")
		return
	}

	for i, entry := range h.entries {
		index := fmt.Sprintf("%4d", i+1)
		
		var prefix string
		if entry.Success {
			prefix = "✓"
		} else {
			prefix = "✗"
		}
		
		var timeStr string
		if showTimestamp {
			timeStr = entry.Timestamp.Format("15:04:05")
		}
		
		var durationStr string
		if showDuration && entry.Duration > 0 {
			durationStr = fmt.Sprintf("(%v)", entry.Duration.Truncate(time.Millisecond))
		}
		
		// Build output line
		var parts []string
		if showTimestamp {
			parts = append(parts, timeStr)
		}
		parts = append(parts, index, prefix, entry.Command)
		if showDuration {
			parts = append(parts, durationStr)
		}
		
		fmt.Println(strings.Join(parts, " "))
	}
}