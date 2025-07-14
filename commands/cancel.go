package commands

import (
	"fmt"
	"slsh/slurm"
)

type CancelCommand struct {
	client *slurm.Client
}

func NewCancelCommand(client *slurm.Client) *CancelCommand {
	return &CancelCommand{client: client}
}

func (c *CancelCommand) Execute(cmd *slurm.Command, shell ShellInterface) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("usage: cancel <job_id>")
	}
	
	jobID := cmd.Args[0]
	_, err := c.client.CancelJob(jobID)
	if err != nil {
		return fmt.Errorf("failed to cancel job: %v", err)
	}
	
	fmt.Printf("Job %s cancelled\n", jobID)
	return nil
}

func (c *CancelCommand) Description() string {
	return "Cancel a job"
}

func (c *CancelCommand) Usage() string {
	return "cancel <job_id> - Cancel a running or pending job"
}