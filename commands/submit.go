package commands

import (
	"fmt"
	"slsh/config"
	"slsh/slurm"
)

type SubmitCommand struct {
	client *slurm.Client
	config *config.Config
}

func NewSubmitCommand(client *slurm.Client, cfg *config.Config) *SubmitCommand {
	return &SubmitCommand{client: client, config: cfg}
}

func (s *SubmitCommand) Execute(cmd *slurm.Command, shell ShellInterface) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("usage: submit <script>")
	}
	
	script := cmd.Args[0]
	jobOpts := parseJobOptions(cmd.Options)
	
	result, err := s.client.SubmitJob(script, jobOpts)
	if err != nil {
		return fmt.Errorf("failed to submit job: %v", err)
	}
	
	if result.Output != "" {
		fmt.Print(result.Output)
	}
	return nil
}

func (s *SubmitCommand) Description() string {
	return "Submit a batch job script"
}

func (s *SubmitCommand) Usage() string {
	return "submit <script> - Submit a job script using sbatch"
}