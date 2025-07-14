package slurm

import "time"

// Job represents a Slurm job
type Job struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	User        string    `json:"user"`
	State       string    `json:"state"`
	Partition   string    `json:"partition"`
	Nodes       int       `json:"nodes"`
	CPUs        int       `json:"cpus"`
	TimeLimit   string    `json:"time_limit"`
	SubmitTime  time.Time `json:"submit_time"`
	StartTime   time.Time `json:"start_time,omitempty"`
	EndTime     time.Time `json:"end_time,omitempty"`
	NodeList    string    `json:"node_list,omitempty"`
	WorkDir     string    `json:"work_dir,omitempty"`
	Command     string    `json:"command,omitempty"`
}

// Node represents a Slurm node
type Node struct {
	Name      string `json:"name"`
	State     string `json:"state"`
	CPUs      int    `json:"cpus"`
	Memory    int    `json:"memory"`
	Partition string `json:"partition"`
	Features  string `json:"features,omitempty"`
}

// Partition represents a Slurm partition
type Partition struct {
	Name        string   `json:"name"`
	State       string   `json:"state"`
	MaxTime     string   `json:"max_time"`
	MaxNodes    int      `json:"max_nodes"`
	DefaultTime string   `json:"default_time"`
	Nodes       []string `json:"nodes"`
}

// JobOptions represents options for job submission
type JobOptions struct {
	Name        string            `json:"name,omitempty"`
	Partition   string            `json:"partition,omitempty"`
	Nodes       int               `json:"nodes,omitempty"`
	CPUs        int               `json:"cpus,omitempty"`
	Memory      string            `json:"memory,omitempty"`
	Time        string            `json:"time,omitempty"`
	QoS         string            `json:"qos,omitempty"`
	Account     string            `json:"account,omitempty"`
	Output      string            `json:"output,omitempty"`
	Error       string            `json:"error,omitempty"`
	WorkDir     string            `json:"work_dir,omitempty"`
	Environment map[string]string `json:"environment,omitempty"`
	ExtraArgs   []string          `json:"extra_args,omitempty"`
}

// Command represents a parsed command
type Command struct {
	Name    string            `json:"name"`
	Args    []string          `json:"args"`
	Options map[string]string `json:"options"`
}

// CommandResult represents the result of a command execution
type CommandResult struct {
	Success    bool   `json:"success"`
	ExitCode   int    `json:"exit_code"`
	Output     string `json:"output"`
	Error      string `json:"error"`
	Duration   time.Duration `json:"duration"`
}

// JobState constants
const (
	JobStatePending    = "PENDING"
	JobStateRunning    = "RUNNING"
	JobStateCompleted  = "COMPLETED"
	JobStateFailed     = "FAILED"
	JobStateCancelled  = "CANCELLED"
	JobStateTimeout    = "TIMEOUT"
)

// Node state constants
const (
	NodeStateIdle     = "IDLE"
	NodeStateAlloc    = "ALLOC"
	NodeStateMixed    = "MIXED"
	NodeStateDown     = "DOWN"
	NodeStateDrain    = "DRAIN"
	NodeStateReserved = "RESERVED"
)