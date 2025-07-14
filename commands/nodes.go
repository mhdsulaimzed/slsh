package commands

import (
	"fmt"
	"slsh/slurm"
)

type NodesCommand struct {
	client *slurm.Client
}

func NewNodesCommand(client *slurm.Client) *NodesCommand {
	return &NodesCommand{client: client}
}

func (n *NodesCommand) Execute(cmd *slurm.Command, shell ShellInterface) error {
	result, err := n.client.GetNodes()
	if err != nil {
		return fmt.Errorf("failed to get nodes: %v", err)
	}
	
	if result.Output != "" {
		fmt.Print(result.Output)
	}
	return nil
}

func (n *NodesCommand) Description() string {
	return "Show node information"
}

func (n *NodesCommand) Usage() string {
	return "nodes - Show cluster node information"
}