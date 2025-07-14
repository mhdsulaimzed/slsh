package commands

import (
	"fmt"
	"os"
	"slsh/slurm"
)

type JobsCommand struct {
	client *slurm.Client
}

func NewJobsCommand(client *slurm.Client) *JobsCommand {
	return &JobsCommand{client: client}
}

func (j *JobsCommand) Execute(cmd *slurm.Command, shell ShellInterface) error {
	user := os.Getenv("USER")
	result, err := j.client.GetQueue(user)
	if err != nil {
		return fmt.Errorf("failed to get jobs: %v", err)
	}
	
	if result.Output != "" {
		fmt.Print(result.Output)
	}
	return nil
}

func (j *JobsCommand) Description() string {
	return "Show your jobs"
}

func (j *JobsCommand) Usage() string {
	return "jobs - Show all your jobs"
}
