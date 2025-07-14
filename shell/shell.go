package shell

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"slsh/commands"
	"slsh/config"
	"slsh/slurm"
	"slsh/utils"
)

// Shell represents the main shell instance
type Shell struct {
	config   *config.Config
	history  *History
	client   *slurm.Client
	commands *commands.Registry
	prompt   *utils.Prompt
	running  bool
}

// New creates a new shell instance
func New() *Shell {
	cfg := config.Load()
	
	return &Shell{
		config:   cfg,
		history:  NewHistory(cfg.HistorySize),
		client:   slurm.NewClient(),
		commands: commands.NewRegistry(),
		prompt:   utils.NewPrompt(cfg.Prompt),
		running:  false,
	}
}

// Run starts the main shell loop
func (s *Shell) Run() error {
	// Load history
	if err := s.history.Load(); err != nil {
		fmt.Printf("Warning: Failed to load history: %v\n", err)
	}

	// Show welcome message
	s.showWelcome()

	// Register built-in commands
	s.registerBuiltinCommands()

	// Main REPL loop
	s.running = true
	scanner := bufio.NewScanner(os.Stdin)
	
	for s.running {
		// Show prompt
		s.prompt.Show()
		
		// Read input
		if !scanner.Scan() {
			break
		}
		
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		
		// Execute command
		s.executeCommand(line)
	}
	
	// Save history before exit
	if err := s.history.Save(); err != nil {
		fmt.Printf("Warning: Failed to save history: %v\n", err)
	}
	
	return scanner.Err()
}

// executeCommand executes a single command
func (s *Shell) executeCommand(line string) {
	startTime := time.Now()
	success := true
	
	// Parse command
	cmd, err := ParseCommand(line)
	if err != nil {
		fmt.Printf("Error parsing command: %v\n", err)
		success = false
		s.history.Add(line, success, time.Since(startTime))
		return
	}
	
	// Validate command
	if err := ValidateCommand(cmd); err != nil {
		fmt.Printf("Invalid command: %v\n", err)
		success = false
		s.history.Add(line, success, time.Since(startTime))
		return
	}
	
	// Check for aliases
	if alias, exists := s.config.Aliases[cmd.Name]; exists {
		// Replace command with alias
		aliasCmd, err := ParseCommand(alias + " " + strings.Join(cmd.Args, " "))
		if err != nil {
			fmt.Printf("Error parsing alias: %v\n", err)
			success = false
			s.history.Add(line, success, time.Since(startTime))
			return
		}
		cmd = aliasCmd
	}
	
	// Execute command
	err = s.commands.Execute(cmd, s)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		success = false
	}
	
	// Add to history
	duration := time.Since(startTime)
	s.history.Add(line, success, duration)
}

// registerBuiltinCommands registers all built-in commands
func (s *Shell) registerBuiltinCommands() {
	// Job execution commands
	s.commands.Register("run", commands.NewRunCommand(s.client, s.config))
	s.commands.Register("submit", commands.NewSubmitCommand(s.client, s.config))
	
	// Job management commands
	s.commands.Register("status", commands.NewStatusCommand(s.client))
	s.commands.Register("cancel", commands.NewCancelCommand(s.client))
	s.commands.Register("queue", commands.NewQueueCommand(s.client))
	s.commands.Register("jobs", commands.NewJobsCommand(s.client))
	
	// Node information commands
	s.commands.Register("nodes", commands.NewNodesCommand(s.client))
	s.commands.Register("partitions", commands.NewPartitionsCommand(s.client))
	
	// Shell management commands
	s.commands.Register("history", commands.NewHistoryCommand(s.history))
	s.commands.Register("alias", commands.NewAliasCommand(s.config))
	s.commands.Register("config", commands.NewConfigCommand(s.config))
	s.commands.Register("help", commands.NewHelpCommand(s.commands))
	s.commands.Register("exit", commands.NewExitCommand(s))
	s.commands.Register("quit", commands.NewExitCommand(s))
	
	// Shortcuts
	s.commands.Register("q", commands.NewQueueCommand(s.client))
	s.commands.Register("j", commands.NewJobsCommand(s.client))
	s.commands.Register("n", commands.NewNodesCommand(s.client))
	s.commands.Register("h", commands.NewHelpCommand(s.commands))
}

// showWelcome displays the welcome message
func (s *Shell) showWelcome() {
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║                    Slurm Shell (slsh) v1.0.0                ║")
	fmt.Println("║                High Performance Computing Shell              ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("Welcome to slsh - A specialized shell for Slurm HPC environments!")
	fmt.Println("Type 'help' for available commands or 'exit' to quit.")
	fmt.Println()
	
	// Show current cluster info if available
	if info := s.client.GetClusterInfo(); info != "" {
		fmt.Printf("Connected to cluster: %s\n", info)
		fmt.Println()
	}
}

// Stop stops the shell
func (s *Shell) Stop() {
	s.running = false
}

// GetConfig returns the shell configuration
func (s *Shell) GetConfig() *config.Config {
	return s.config
}

// GetHistory returns the shell history
func (s *Shell) GetHistory() *History {
	return s.history
}

// GetClient returns the Slurm client
func (s *Shell) GetClient() *slurm.Client {
	return s.client
}

// UpdatePrompt updates the shell prompt
func (s *Shell) UpdatePrompt(newPrompt string) {
	s.prompt.SetPrompt(newPrompt)
}

// ExecuteDirectCommand executes a command directly (for testing or API use)
func (s *Shell) ExecuteDirectCommand(command string) error {
	cmd, err := ParseCommand(command)
	if err != nil {
		return fmt.Errorf("failed to parse command: %v", err)
	}
	
	if err := ValidateCommand(cmd); err != nil {
		return fmt.Errorf("invalid command: %v", err)
	}
	
	return s.commands.Execute(cmd, s)
}

// GetAvailableCommands returns a list of available commands
func (s *Shell) GetAvailableCommands() []string {
	return s.commands.GetCommandNames()
}

// IsRunning returns whether the shell is currently running
func (s *Shell) IsRunning() bool {
	return s.running
}

// AddAlias adds a new alias
func (s *Shell) AddAlias(name, command string) {
	if s.config.Aliases == nil {
		s.config.Aliases = make(map[string]string)
	}
	s.config.Aliases[name] = command
}

// RemoveAlias removes an alias
func (s *Shell) RemoveAlias(name string) {
	if s.config.Aliases != nil {
		delete(s.config.Aliases, name)
	}
}

// GetAliases returns all aliases
func (s *Shell) GetAliases() map[string]string {
	return s.config.Aliases
}