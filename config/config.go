package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config holds the shell configuration
type Config struct {
	// Default job settings
	DefaultPartition string `json:"default_partition"`
	DefaultNodes     int    `json:"default_nodes"`
	DefaultCPUs      int    `json:"default_cpus"`
	DefaultMemory    string `json:"default_memory"`
	DefaultTime      string `json:"default_time"`
	DefaultQoS       string `json:"default_qos"`
	DefaultAccount   string `json:"default_account"`
	
	// Shell settings
	Prompt         string            `json:"prompt"`
	HistorySize    int               `json:"history_size"`
	AutoComplete   bool              `json:"auto_complete"`
	ShowTimestamps bool              `json:"show_timestamps"`
	ColorOutput    bool              `json:"color_output"`
	
	// Aliases
	Aliases map[string]string `json:"aliases"`
	
	// Output settings
	DefaultOutputDir string `json:"default_output_dir"`
	JobNameTemplate  string `json:"job_name_template"`
	
	// Advanced settings
	CommandTimeout   int  `json:"command_timeout_seconds"`
	ConfirmDangerous bool `json:"confirm_dangerous_operations"`
	SaveJobHistory   bool `json:"save_job_history"`
}

// Default returns a configuration with sensible defaults
func Default() *Config {
	homeDir, _ := os.UserHomeDir()
	
	return &Config{
		// Job defaults
		DefaultPartition: "",
		DefaultNodes:     1,
		DefaultCPUs:      1,
		DefaultMemory:    "",
		DefaultTime:      "01:00:00",
		DefaultQoS:       "",
		DefaultAccount:   "",
		
		// Shell settings
		Prompt:         "slsh> ",
		HistorySize:    1000,
		AutoComplete:   true,
		ShowTimestamps: false,
		ColorOutput:    true,
		
		// Aliases
		Aliases: map[string]string{
			"q":  "queue",
			"j":  "jobs",
			"n":  "nodes",
			"h":  "help",
			"st": "status",
		},
		
		// Output settings
		DefaultOutputDir: filepath.Join(homeDir, "slurm_jobs"),
		JobNameTemplate:  "job_%j",
		
		// Advanced settings
		CommandTimeout:   30,
		ConfirmDangerous: true,
		SaveJobHistory:   true,
	}
}

// Load loads configuration from file or returns default
func Load() *Config {
	config := Default()
	
	configPath := GetConfigPath()
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// No config file exists, create one with defaults
		config.Save()
		return config
	}
	
	data, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Printf("Warning: Could not read config file: %v\n", err)
		return config
	}
	
	if err := json.Unmarshal(data, config); err != nil {
		fmt.Printf("Warning: Could not parse config file: %v\n", err)
		return Default()
	}
	
	return config
}

// Save saves the configuration to file
func (c *Config) Save() error {
	configPath := GetConfigPath()
	configDir := filepath.Dir(configPath)
	
	// Ensure config directory exists
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}
	
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}
	
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}
	
	return nil
}

// GetConfigPath returns the path to the configuration file
func GetConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "/tmp/.slshrc"
	}
	
	return filepath.Join(homeDir, ".config", "slsh", "config.json")
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.HistorySize < 0 {
		return fmt.Errorf("history_size cannot be negative")
	}
	
	if c.CommandTimeout < 1 {
		return fmt.Errorf("command_timeout_seconds must be at least 1")
	}
	
	// Validate time format if set
	if c.DefaultTime != "" && !isValidTimeFormat(c.DefaultTime) {
		return fmt.Errorf("invalid default_time format: %s", c.DefaultTime)
	}
	
	return nil
}

// isValidTimeFormat checks if time is in valid Slurm format
func isValidTimeFormat(timeStr string) bool {
	// Basic validation for HH:MM:SS, MM:SS, or minutes
	// This could be enhanced with more sophisticated parsing
	return len(timeStr) > 0
}

// SetAlias sets an alias
func (c *Config) SetAlias(name, command string) {
	if c.Aliases == nil {
		c.Aliases = make(map[string]string)
	}
	c.Aliases[name] = command
}

// RemoveAlias removes an alias
func (c *Config) RemoveAlias(name string) {
	if c.Aliases != nil {
		delete(c.Aliases, name)
	}
}

// GetAlias gets an alias
func (c *Config) GetAlias(name string) (string, bool) {
	if c.Aliases == nil {
		return "", false
	}
	alias, exists := c.Aliases[name]
	return alias, exists
}

// UpdateDefaults updates default job settings
func (c *Config) UpdateDefaults(partition string, nodes int, cpus int, memory string, time string) {
	if partition != "" {
		c.DefaultPartition = partition
	}
	if nodes > 0 {
		c.DefaultNodes = nodes
	}
	if cpus > 0 {
		c.DefaultCPUs = cpus
	}
	if memory != "" {
		c.DefaultMemory = memory
	}
	if time != "" {
		c.DefaultTime = time
	}
}

// Print prints the configuration in a readable format
func (c *Config) Print() {
	fmt.Println("=== Slurm Shell Configuration ===")
	fmt.Println()
	
	fmt.Println("Job Defaults:")
	fmt.Printf("  Partition: %s\n", c.DefaultPartition)
	fmt.Printf("  Nodes: %d\n", c.DefaultNodes)
	fmt.Printf("  CPUs: %d\n", c.DefaultCPUs)
	fmt.Printf("  Memory: %s\n", c.DefaultMemory)
	fmt.Printf("  Time: %s\n", c.DefaultTime)
	fmt.Printf("  QoS: %s\n", c.DefaultQoS)
	fmt.Printf("  Account: %s\n", c.DefaultAccount)
	fmt.Println()
	
	fmt.Println("Shell Settings:")
	fmt.Printf("  Prompt: %s\n", c.Prompt)
	fmt.Printf("  History Size: %d\n", c.HistorySize)
	fmt.Printf("  Auto Complete: %t\n", c.AutoComplete)
	fmt.Printf("  Show Timestamps: %t\n", c.ShowTimestamps)
	fmt.Printf("  Color Output: %t\n", c.ColorOutput)
	fmt.Println()
	
	if len(c.Aliases) > 0 {
		fmt.Println("Aliases:")
		for name, command := range c.Aliases {
			fmt.Printf("  %-10s = %s\n", name, command)
		}
		fmt.Println()
	}
	
	fmt.Println("Output Settings:")
	fmt.Printf("  Default Output Dir: %s\n", c.DefaultOutputDir)
	fmt.Printf("  Job Name Template: %s\n", c.JobNameTemplate)
	fmt.Println()
	
	fmt.Println("Advanced Settings:")
	fmt.Printf("  Command Timeout: %d seconds\n", c.CommandTimeout)
	fmt.Printf("  Confirm Dangerous Operations: %t\n", c.ConfirmDangerous)
	fmt.Printf("  Save Job History: %t\n", c.SaveJobHistory)
}