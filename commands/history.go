package commands

import (

	"slsh/shell"
	"slsh/slurm"
	
)

type HistoryCommand struct {
	history *shell.History
}

func NewHistoryCommand(history *shell.History) *HistoryCommand {
	return &HistoryCommand{history: history}
}

func (h *HistoryCommand) Execute(cmd *slurm.Command, shell ShellInterface) error {
	showTime := false
	showDuration := false
	
	// Parse options
	for opt := range cmd.Options {
		switch opt {
		case "-t", "--time":
			showTime = true
		case "-d", "--duration":
			showDuration = true
		}
	}
	
	h.history.PrintHistory(showTime, showDuration)
	return nil
}

func (h *HistoryCommand) Description() string {
	return "Show command history"
}

func (h *HistoryCommand) Usage() string {
	return `history [-t] [-d]

Show command history.

Options:
  -t, --time      Show timestamps
  -d, --duration  Show execution duration`
}
