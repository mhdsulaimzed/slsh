package shell

import (
	"fmt"
	"strings"
	"unicode"

	"slsh/slurm"
)

// ParseCommand parses a command line into a Command struct
func ParseCommand(line string) (*slurm.Command, error) {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil, fmt.Errorf("empty command")
	}

	tokens, err := tokenize(line)
	if err != nil {
		return nil, fmt.Errorf("failed to tokenize command: %v", err)
	}

	if len(tokens) == 0 {
		return nil, fmt.Errorf("no tokens found")
	}

	cmd := &slurm.Command{
		Name:    tokens[0],
		Args:    []string{},
		Options: make(map[string]string),
	}

	// Parse tokens into args and options
	for i := 1; i < len(tokens); i++ {
		token := tokens[i]

		if strings.HasPrefix(token, "-") {
			// This is an option
			if i+1 < len(tokens) && !strings.HasPrefix(tokens[i+1], "-") {
				// Option has a value
				cmd.Options[token] = tokens[i+1]
				i++ // Skip the next token as it's the value
			} else {
				// Option is a flag (no value)
				cmd.Options[token] = ""
			}
		} else {
			// This is an argument
			cmd.Args = append(cmd.Args, token)
		}
	}

	return cmd, nil
}

// tokenize splits a command line into tokens, handling quotes and escapes
func tokenize(line string) ([]string, error) {
	var tokens []string
	var current strings.Builder
	var inQuotes bool
	var quoteChar rune
	var escaped bool

	for _, char := range line {
		if escaped {
			current.WriteRune(char)
			escaped = false
			continue
		}

		switch char {
		case '\\':
			escaped = true
		case '"', '\'':
			if !inQuotes {
				inQuotes = true
				quoteChar = char
			} else if char == quoteChar {
				inQuotes = false
			} else {
				current.WriteRune(char)
			}
		case ' ', '\t':
			if inQuotes {
				current.WriteRune(char)
			} else {
				if current.Len() > 0 {
					tokens = append(tokens, current.String())
					current.Reset()
				}
			}
		default:
			current.WriteRune(char)
		}
	}

	if inQuotes {
		return nil, fmt.Errorf("unclosed quote")
	}

	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}

	return tokens, nil
}

// ParseJobOptions parses command options into JobOptions struct
func ParseJobOptions(options map[string]string) *slurm.JobOptions {
	jobOpts := &slurm.JobOptions{
		Environment: make(map[string]string),
	}

	for opt, value := range options {
		switch opt {
		case "-J", "--job-name":
			jobOpts.Name = value
		case "-p", "--partition":
			jobOpts.Partition = value
		case "-N", "--nodes":
			if nodes := parseInt(value); nodes > 0 {
				jobOpts.Nodes = nodes
			}
		case "-c", "--cpus-per-task":
			if cpus := parseInt(value); cpus > 0 {
				jobOpts.CPUs = cpus
			}
		case "--mem":
			jobOpts.Memory = value
		case "-t", "--time":
			jobOpts.Time = value
		case "--qos":
			jobOpts.QoS = value
		case "-A", "--account":
			jobOpts.Account = value
		case "-o", "--output":
			jobOpts.Output = value
		case "-e", "--error":
			jobOpts.Error = value
		case "-D", "--chdir":
			jobOpts.WorkDir = value
		default:
			// Store unknown options as extra args
			if value != "" {
				jobOpts.ExtraArgs = append(jobOpts.ExtraArgs, opt, value)
			} else {
				jobOpts.ExtraArgs = append(jobOpts.ExtraArgs, opt)
			}
		}
	}

	return jobOpts
}

// parseInt safely parses a string to int
func parseInt(s string) int {
	var result int
	for _, char := range s {
		if unicode.IsDigit(char) {
			result = result*10 + int(char-'0')
		} else {
			return 0 // Invalid number
		}
	}
	return result
}

// ValidateCommand performs basic validation on a command
func ValidateCommand(cmd *slurm.Command) error {
	if cmd.Name == "" {
		return fmt.Errorf("command name cannot be empty")
	}

	// Check for conflicting options
	if _, hasNodes := cmd.Options["-N"]; hasNodes {
		if _, hasNodeList := cmd.Options["-w"]; hasNodeList {
			return fmt.Errorf("cannot specify both -N (nodes) and -w (nodelist)")
		}
	}

	// Validate time format if specified
	if timeLimit, exists := cmd.Options["-t"]; exists && timeLimit != "" {
		if !isValidTimeFormat(timeLimit) {
			return fmt.Errorf("invalid time format: %s (use format: HH:MM:SS or minutes)", timeLimit)
		}
	}

	return nil
}

// isValidTimeFormat checks if a time string is in valid Slurm format
func isValidTimeFormat(timeStr string) bool {
	// Simple validation - accepts HH:MM:SS, MM:SS, or just minutes
	if strings.Contains(timeStr, ":") {
		parts := strings.Split(timeStr, ":")
		if len(parts) < 2 || len(parts) > 3 {
			return false
		}
		// Could add more detailed validation here
		return true
	}

	// Check if it's just a number (minutes)
	for _, char := range timeStr {
		if !unicode.IsDigit(char) {
			return false
		}
	}
	return len(timeStr) > 0
}

// SplitCommandLine splits a command line respecting quotes and escapes
func SplitCommandLine(line string) []string {
	tokens, err := tokenize(line)
	if err != nil {
		// Fallback to simple split if tokenization fails
		return strings.Fields(line)
	}
	return tokens
}
