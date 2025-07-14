package utils

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// Prompt handles shell prompt display and formatting
type Prompt struct {
	template string
	showTime bool
	showUser bool
	showHost bool
	showCwd  bool
}

// NewPrompt creates a new prompt with the given template
func NewPrompt(template string) *Prompt {
	if template == "" {
		template = "slsh> "
	}
	
	return &Prompt{
		template: template,
		showTime: strings.Contains(template, "%t"),
		showUser: strings.Contains(template, "%u"),
		showHost: strings.Contains(template, "%h"),
		showCwd:  strings.Contains(template, "%w"),
	}
}

// Show displays the prompt
func (p *Prompt) Show() {
	fmt.Print(p.Format())
}

// Format formats the prompt string
func (p *Prompt) Format() string {
	prompt := p.template
	
	// Replace template variables
	if p.showTime {
		timeStr := time.Now().Format("15:04:05")
		prompt = strings.ReplaceAll(prompt, "%t", timeStr)
	}
	
	if p.showUser {
		user := os.Getenv("USER")
		if user == "" {
			user = "unknown"
		}
		prompt = strings.ReplaceAll(prompt, "%u", user)
	}
	
	if p.showHost {
		hostname, err := os.Hostname()
		if err != nil {
			hostname = "unknown"
		}
		prompt = strings.ReplaceAll(prompt, "%h", hostname)
	}
	
	if p.showCwd {
		cwd, err := os.Getwd()
		if err != nil {
			cwd = "unknown"
		}
		// Show only the basename for brevity
		parts := strings.Split(cwd, "/")
		if len(parts) > 0 {
			cwd = parts[len(parts)-1]
		}
		prompt = strings.ReplaceAll(prompt, "%w", cwd)
	}
	
	return prompt
}

// SetPrompt sets a new prompt template
func (p *Prompt) SetPrompt(template string) {
	p.template = template
	p.showTime = strings.Contains(template, "%t")
	p.showUser = strings.Contains(template, "%u")
	p.showHost = strings.Contains(template, "%h")
	p.showCwd = strings.Contains(template, "%w")
}

// GetTemplate returns the current prompt template
func (p *Prompt) GetTemplate() string {
	return p.template
}

// Static prompt formats for common use cases
var (
	SimplePrompt    = "slsh> "
	TimedPrompt     = "[%t] slsh> "
	UserPrompt      = "%u@slsh> "
	FullPrompt      = "[%t] %u@%h:%w slsh> "
	MinimalPrompt   = "> "
	ColorPrompt     = "\033[32mslsh\033[0m> "
)

// SetBuiltinPrompt sets a predefined prompt style
func (p *Prompt) SetBuiltinPrompt(style string) {
	switch style {
	case "simple":
		p.SetPrompt(SimplePrompt)
	case "timed":
		p.SetPrompt(TimedPrompt)
	case "user":
		p.SetPrompt(UserPrompt)
	case "full":
		p.SetPrompt(FullPrompt)
	case "minimal":
		p.SetPrompt(MinimalPrompt)
	case "color":
		p.SetPrompt(ColorPrompt)
	default:
		p.SetPrompt(SimplePrompt)
	}
}