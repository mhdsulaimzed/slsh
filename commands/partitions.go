package commands
import (
	"fmt"
	"slsh/slurm"
)

type PartitionsCommand struct {
	client *slurm.Client
}

func NewPartitionsCommand(client *slurm.Client) *PartitionsCommand {
	return &PartitionsCommand{client: client}
}

func (p *PartitionsCommand) Execute(cmd *slurm.Command, shell ShellInterface) error {
	result, err := p.client.GetPartitions()
	if err != nil {
		return fmt.Errorf("failed to get partitions: %v", err)
	}
	
	if result.Output != "" {
		fmt.Print(result.Output)
	}
	return nil
}

func (p *PartitionsCommand) Description() string {
	return "Show partition information"
}

func (p *PartitionsCommand) Usage() string {
	return "partitions - Show cluster partition information"
}