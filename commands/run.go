package commands

import (
	"fmt"
	"strings"

	"slsh/config"
	"slsh/slurm"
	"slsh/utils"
)

// RunCommand implements the 'run' command
type RunCommand struct {
	client *slurm.Client
	config *config.Config
}

// NewRunCommand creates a new run command
func NewRunCommand(client *slurm.Client, cfg *config.Config) *RunCommand {
	return &RunCommand{
		client: client,
		config: cfg,
	}
}

// Execute executes the run command
func (r *RunCommand) Execute(cmd *slurm.Command, shell ShellInterface) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("usage: run <command> [arguments...]")
	}
	
	// Parse job options from command
	jobOpts := parseJobOptions(cmd.Options)
	
	// Apply defaults from config
	r.applyDefaults(jobOpts)
	
	// Build the command to execute
	command := strings.Join(cmd.Args, " ")
	
	// Show what we're about to execute
	fmt.Printf("Running: %s\n", command)
	if jobOpts.Partition != "" {
		fmt.Printf("Partition: %s\n", jobOpts.Partition)
	}
	if jobOpts.Nodes > 0 {
		fmt.Printf("Nodes: %d\n", jobOpts.Nodes)
	}
	if jobOpts.Time != "" {
		fmt.Printf("Time limit: %s\n", jobOpts.Time)
	}
	fmt.Println()
	
	// Execute the job
	result, err := r.client.RunJob(command, jobOpts)
	if err != nil {
		return fmt.Errorf("failed to run job: %v", err)
	}
	
	// Display output
	if result.Output != "" {
		fmt.Print(result.Output)
	}
	if result.Error != "" {
		fmt.Print(result.Error)
	}
	
	// Show completion status
	if result.Success {
		fmt.Printf(utils.FormatSuccess("Job completed successfully", r.config.ColorOutput))
	} else {
		fmt.Printf(utils.FormatError("Job failed with exit code %d", r.config.ColorOutput), result.ExitCode)
	}
	fmt.Printf(" (Duration: %s)\n", utils.FormatDuration(result.Duration))
	
	return nil
}

// applyDefaults applies default configuration to job options
func (r *RunCommand) applyDefaults(opts *slurm.JobOptions) {
	if opts.Partition == "" && r.config.DefaultPartition != "" {
		opts.Partition = r.config.DefaultPartition
	}
	
	if opts.Nodes == 0 && r.config.DefaultNodes > 0 {
		opts.Nodes = r.config.DefaultNodes
	}
	
	if opts.CPUs == 0 && r.config.DefaultCPUs > 0 {
		opts.CPUs = r.config.DefaultCPUs
	}
	
	if opts.Memory == "" && r.config.DefaultMemory != "" {
		opts.Memory = r.config.DefaultMemory
	}
	
	if opts.Time == "" && r.config.DefaultTime != "" {
		opts.Time = r.config.DefaultTime
	}
	
	if opts.QoS == "" && r.config.DefaultQoS != "" {
		opts.QoS = r.config.DefaultQoS
	}
	
	if opts.Account == "" && r.config.DefaultAccount != "" {
		opts.Account = r.config.DefaultAccount
	}
}

// Description returns the command description
func (r *RunCommand) Description() string {
	return "Execute a command using srun with configured defaults"
}

// Usage returns the command usage
func (r *RunCommand) Usage() string {
	return `run [OPTIONS] <command> [arguments...]

Execute a command on the cluster using srun. This command applies
your configured defaults and provides a simplified interface.

Examples:
  run hostname                    # Run hostname on default resources
  run -N 2 hostname               # Run on 2 nodes
  run -p gpu nvidia-smi           # Run on GPU partition
  run -t 30:00 ./my_simulation    # Run with 30 minute time limit

Options:
  -J, --job-name <name>           Job name
  -p, --partition <partition>     Partition to use
  -N, --nodes <count>             Number of nodes
  -c, --cpus-per-task <count>     CPUs per task
  --mem <memory>                  Memory per node
  -t, --time <time>               Time limit (HH:MM:SS)
  --qos <qos>                     Quality of Service
  -A, --account <account>         Account to charge
  -o, --output <file>             Output file
  -e, --error <file>              Error file

The command will use your configured defaults for any options not specified.`
}

// parseJobOptions parses command options into JobOptions struct
func parseJobOptions(options map[string]string) *slurm.JobOptions {
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
		if char >= '0' && char <= '9' {
			result = result*10 + int(char-'0')
		} else {
			return 0 // Invalid number
		}
	}
	return result
}