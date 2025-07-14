package commands

import (
	"fmt"
	"strings"

	"slsh/slurm"
	"slsh/utils"
)

// HelpCommand implements the 'help' command
type HelpCommand struct {
	registry *Registry
}

// NewHelpCommand creates a new help command
func NewHelpCommand(registry *Registry) *HelpCommand {
	return &HelpCommand{
		registry: registry,
	}
}

// Execute executes the help command
func (h *HelpCommand) Execute(cmd *slurm.Command, shell ShellInterface) error {
	if len(cmd.Args) == 0 {
		// Show general help
		h.showGeneralHelp()
		return nil
	}
	
	// Show help for specific command
	commandName := cmd.Args[0]
	if handler, exists := h.registry.GetCommand(commandName); exists {
		h.showCommandHelp(commandName, handler)
	} else {
		fmt.Printf("Unknown command: %s\n", commandName)
		fmt.Println("Use 'help' to see available commands.")
	}
	
	return nil
}

// showGeneralHelp displays the general help message
func (h *HelpCommand) showGeneralHelp() {
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║                    Slurm Shell (slsh) Help                  ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
	fmt.Println()
	
	fmt.Println("slsh is a specialized shell for Slurm HPC environments that provides")
	fmt.Println("simplified commands and intelligent defaults for job management.")
	fmt.Println()
	
	// Built-in commands
	fmt.Println("Built-in Commands:")
	fmt.Println("==================")
	
	table := utils.NewTable([]string{"Command", "Description"}, true)
	
	commands := h.registry.GetCommands()
	commandNames := h.registry.GetCommandNames()
	
	for _, name := range commandNames {
		if handler, exists := commands[name]; exists {
			// Skip single-letter aliases for cleaner display
			if len(name) == 1 {
				continue
			}
			table.AddRow([]string{name, handler.Description()})
		}
	}
	
	table.Print()
	fmt.Println()
	
	// Aliases
	fmt.Println("Quick Aliases:")
	fmt.Println("==============")
	aliasTable := utils.NewTable([]string{"Alias", "Command"}, true)
	aliasTable.AddRow([]string{"q", "queue"})
	aliasTable.AddRow([]string{"j", "jobs"})
	aliasTable.AddRow([]string{"n", "nodes"})
	aliasTable.AddRow([]string{"h", "help"})
	aliasTable.Print()
	fmt.Println()
	
	// Slurm commands
	fmt.Println("Slurm Commands:")
	fmt.Println("===============")
	fmt.Println("All standard Slurm commands are available:")
	fmt.Println("  srun, sbatch, scancel, squeue, sinfo, sacct, scontrol, etc.")
	fmt.Println()
	
	// Usage examples
	fmt.Println("Usage Examples:")
	fmt.Println("===============")
	fmt.Println("  run hostname                    # Execute hostname on cluster")
	fmt.Println("  run -N 2 -p gpu nvidia-smi     # Run on 2 GPU nodes")
	fmt.Println("  submit my_job.sh               # Submit batch job")
	fmt.Println("  queue                          # Show job queue")
	fmt.Println("  status 12345                   # Check job 12345 status")
	fmt.Println("  cancel 12345                   # Cancel job 12345")
	fmt.Println("  nodes                          # Show node information")
	fmt.Println("  config                         # Show configuration")
	fmt.Println("  alias myrun \"run -N 4 -p gpu\"   # Create custom alias")
	fmt.Println()
	
	fmt.Println("For detailed help on a specific command, use: help <command>")
	fmt.Println("For configuration options, use: help config")
	fmt.Println()
}

// showCommandHelp displays help for a specific command
func (h *HelpCommand) showCommandHelp(name string, handler CommandHandler) {
	fmt.Printf("Command: %s\n", name)
	fmt.Printf("Description: %s\n\n", handler.Description())
	
	usage := handler.Usage()
	if usage != "" {
		fmt.Println("Usage:")
		fmt.Println(strings.ReplaceAll(usage, "\n", "\n"))
		fmt.Println()
	}
}

// Description returns the command description
func (h *HelpCommand) Description() string {
	return "Show help information for commands"
}

// Usage returns the command usage
func (h *HelpCommand) Usage() string {
	return `help [command]

Show help information. Without arguments, shows general help.
With a command name, shows detailed help for that command.

Examples:
  help            # Show general help
  help run        # Show help for 'run' command
  help config     # Show configuration help`
}