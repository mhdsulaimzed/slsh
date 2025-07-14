package commands

import (
	"fmt"
	"os"

	"slsh/slurm"
)

// QueueCommand implements the 'queue' command
type QueueCommand struct {
	client *slurm.Client
}

// NewQueueCommand creates a new queue command
func NewQueueCommand(client *slurm.Client) *QueueCommand {
	return &QueueCommand{
		client: client,
	}
}

// Execute executes the queue command
func (q *QueueCommand) Execute(cmd *slurm.Command, shell ShellInterface) error {
	var user string
	
	// Check if user is specified
	if len(cmd.Args) > 0 {
		user = cmd.Args[0]
	} else {
		// Default to current user
		user = os.Getenv("USER")
	}
	
	result, err := q.client.GetQueue(user)
	if err != nil {
		return fmt.Errorf("failed to get queue: %v", err)
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
func (q *QueueCommand) Description() string {
	return "Show the job queue"
}

// Usage returns the command usage
func (q *QueueCommand) Usage() string {
	return `queue [user]

Show the job queue. Without arguments, shows jobs for current user.
With a username, shows jobs for that user (if you have permission).

Examples:
  queue           # Show your jobs
  queue alice     # Show alice's jobs
  queue --all     # Show all jobs (if supported)`
}