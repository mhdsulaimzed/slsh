package commands

import (
	"fmt"
	"sort"

	"slsh/slurm"
)

// CommandHandler represents a command handler function
type CommandHandler interface {
	Execute(cmd *slurm.Command, shell ShellInterface) error
	Description() string
	Usage() string
}

// ShellInterface defines the interface that commands can use to interact with the shell
type ShellInterface interface {
	GetConfig() interface{}
	GetHistory() interface{}
	GetClient() *slurm.Client
	Stop()
	AddAlias(name, command string)
	RemoveAlias(name string)
	GetAliases() map[string]string
}

// Registry manages command registration and execution
type Registry struct {
	commands map[string]CommandHandler
}

// NewRegistry creates a new command registry
func NewRegistry() *Registry {
	return &Registry{
		commands: make(map[string]CommandHandler),
	}
}

// Register registers a command handler
func (r *Registry) Register(name string, handler CommandHandler) {
	r.commands[name] = handler
}

// Execute executes a command
func (r *Registry) Execute(cmd *slurm.Command, shell ShellInterface) error {
	handler, exists := r.commands[cmd.Name]
	if !exists {
		// Try to execute as a system/Slurm command
		return r.executeSystemCommand(cmd, shell)
	}
	
	return handler.Execute(cmd, shell)
}

// GetCommand returns a command handler by name
func (r *Registry) GetCommand(name string) (CommandHandler, bool) {
	handler, exists := r.commands[name]
	return handler, exists
}

// GetCommandNames returns a sorted list of all command names
func (r *Registry) GetCommandNames() []string {
	names := make([]string, 0, len(r.commands))
	for name := range r.commands {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// GetCommands returns all registered commands
func (r *Registry) GetCommands() map[string]CommandHandler {
	return r.commands
}

// executeSystemCommand executes a system command or Slurm command
func (r *Registry) executeSystemCommand(cmd *slurm.Command, shell ShellInterface) error {
	client := shell.GetClient()
	
	// Check if it's a known Slurm command
	slurmCommands := []string{
		"srun", "sbatch", "scancel", "squeue", "sinfo", "sacct", 
		"scontrol", "sstat", "sprio", "sshare", "sreport", 
		"salloc", "sattach", "sacctmgr",
	}
	
	isSlurmCommand := false
	for _, slurmCmd := range slurmCommands {
		if cmd.Name == slurmCmd {
			isSlurmCommand = true
			break
		}
	}
	
	if isSlurmCommand {
		// Execute as Slurm command
		args := buildArgs(cmd)
		result, err := client.Execute(cmd.Name, args...)
		if err != nil {
			return err
		}
		
		if result.Output != "" {
			fmt.Print(result.Output)
		}
		if result.Error != "" {
			fmt.Print(result.Error)
		}
		
		return nil
	}
	
	// Execute as interactive system command
	args := buildArgs(cmd)
	return client.ExecuteInteractive(cmd.Name, args...)
}

// buildArgs builds command arguments from a Command struct
func buildArgs(cmd *slurm.Command) []string {
	var args []string
	
	// Add options
	for opt, value := range cmd.Options {
		args = append(args, opt)
		if value != "" {
			args = append(args, value)
		}
	}
	
	// Add positional arguments
	args = append(args, cmd.Args...)
	
	return args
}