package commands

import (
	"fmt"

	"slsh/slurm"
)

// StatusCommand implements the 'status' command
type StatusCommand struct {
	client *slurm.Client
}

// NewStatusCommand creates a new status command
func NewStatusCommand(client *slurm.Client) *StatusCommand {
	return &StatusCommand{
		client: client,
	}
}

// Execute executes the status command
func (s *StatusCommand) Execute(cmd *slurm.Command, shell ShellInterface) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("usage: status <job_id>")
	}
	
	jobID := cmd.Args[0]
	result, err := s.client.GetJobStatus(jobID)
	if err != nil {
		return fmt.Errorf("failed to get job status: %v", err)
	}
	
	if result.Output != "" {
		fmt.Print(result.Output)
	}
	if result.Error != "" {
		fmt.Print(result.Error)
	}
	
	return nil
}

// Description returns the command description
func (s *StatusCommand) Description() string {
	return "Show status of a specific job"
}

// Usage returns the command usage
func (s *StatusCommand) Usage() string {
	return `status <job_id>

Show detailed status information for a specific job.

Examples:
  status 12345    # Show status of job 12345`
}