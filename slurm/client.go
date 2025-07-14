package slurm

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// Client handles Slurm command execution
type Client struct {
	timeout time.Duration
}

// NewClient creates a new Slurm client
func NewClient() *Client {
	return &Client{
		timeout: 30 * time.Second,
	}
}

// Execute executes a Slurm command with the given arguments
func (c *Client) Execute(command string, args ...string) (*CommandResult, error) {
	start := time.Now()
	
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	
	cmd := exec.CommandContext(ctx, command, args...)
	
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	err := cmd.Run()
	
	result := &CommandResult{
		Success:  err == nil,
		Output:   stdout.String(),
		Error:    stderr.String(),
		Duration: time.Since(start),
	}
	
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
		} else {
			result.ExitCode = -1
		}
		return result, fmt.Errorf("command failed: %v", err)
	}
	
	return result, nil
}

// RunJob submits and runs a job using srun
func (c *Client) RunJob(command string, options *JobOptions) (*CommandResult, error) {
	args := []string{}
	
	// Add job options
	if options != nil {
		args = append(args, c.buildJobArgs(options)...)
	}
	
	// Add the command to execute
	if command != "" {
		args = append(args, command)
	}
	
	return c.Execute("srun", args...)
}

// SubmitJob submits a job using sbatch
func (c *Client) SubmitJob(scriptPath string, options *JobOptions) (*CommandResult, error) {
	args := []string{}
	
	// Add job options
	if options != nil {
		args = append(args, c.buildJobArgs(options)...)
	}
	
	// Add script path
	args = append(args, scriptPath)
	
	return c.Execute("sbatch", args...)
}

// CancelJob cancels a job using scancel
func (c *Client) CancelJob(jobID string) (*CommandResult, error) {
	return c.Execute("scancel", jobID)
}

// GetJobStatus gets status of a specific job
func (c *Client) GetJobStatus(jobID string) (*CommandResult, error) {
	return c.Execute("squeue", "-j", jobID, "--format=%i,%T,%P,%u,%M,%N,%r")
}

// GetQueue gets the job queue
func (c *Client) GetQueue(user string) (*CommandResult, error) {
	args := []string{"--format=%i,%T,%P,%u,%M,%N,%j"}
	if user != "" {
		args = append(args, "-u", user)
	}
	return c.Execute("squeue", args...)
}

// GetNodes gets node information
func (c *Client) GetNodes() (*CommandResult, error) {
	return c.Execute("sinfo", "-N", "--format=%N,%T,%P,%C,%m,%f")
}

// GetPartitions gets partition information
func (c *Client) GetPartitions() (*CommandResult, error) {
	return c.Execute("sinfo", "--format=%P,%a,%l,%D,%N")
}

// GetAccountInfo gets account information for a user
func (c *Client) GetAccountInfo(user string) (*CommandResult, error) {
	if user == "" {
		user = os.Getenv("USER")
	}
	return c.Execute("sacctmgr", "show", "user", user, "-s")
}

// GetClusterInfo gets basic cluster information
func (c *Client) GetClusterInfo() string {
	result, err := c.Execute("scontrol", "show", "config")
	if err != nil {
		return ""
	}
	
	// Parse cluster name from output
	lines := strings.Split(result.Output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "ClusterName") {
			parts := strings.Split(line, "=")
			if len(parts) >= 2 {
				return strings.TrimSpace(parts[1])
			}
		}
	}
	
	return "Unknown"
}

// CheckSlurmAvailable checks if Slurm commands are available
func (c *Client) CheckSlurmAvailable() error {
	_, err := exec.LookPath("srun")
	if err != nil {
		return fmt.Errorf("srun not found in PATH - is Slurm installed?")
	}
	
	_, err = exec.LookPath("squeue")
	if err != nil {
		return fmt.Errorf("squeue not found in PATH - is Slurm installed?")
	}
	
	return nil
}

// buildJobArgs builds command line arguments from JobOptions
func (c *Client) buildJobArgs(options *JobOptions) []string {
	var args []string
	
	if options.Name != "" {
		args = append(args, "--job-name="+options.Name)
	}
	
	if options.Partition != "" {
		args = append(args, "--partition="+options.Partition)
	}
	
	if options.Nodes > 0 {
		args = append(args, fmt.Sprintf("--nodes=%d", options.Nodes))
	}
	
	if options.CPUs > 0 {
		args = append(args, fmt.Sprintf("--cpus-per-task=%d", options.CPUs))
	}
	
	if options.Memory != "" {
		args = append(args, "--mem="+options.Memory)
	}
	
	if options.Time != "" {
		args = append(args, "--time="+options.Time)
	}
	
	if options.QoS != "" {
		args = append(args, "--qos="+options.QoS)
	}
	
	if options.Account != "" {
		args = append(args, "--account="+options.Account)
	}
	
	if options.Output != "" {
		args = append(args, "--output="+options.Output)
	}
	
	if options.Error != "" {
		args = append(args, "--error="+options.Error)
	}
	
	if options.WorkDir != "" {
		args = append(args, "--chdir="+options.WorkDir)
	}
	
	// Add environment variables
	for key, value := range options.Environment {
		args = append(args, "--export="+key+"="+value)
	}
	
	// Add extra arguments
	args = append(args, options.ExtraArgs...)
	
	return args
}

// ExecuteInteractive runs a command interactively (with stdin/stdout)
func (c *Client) ExecuteInteractive(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	return cmd.Run()
}

// SetTimeout sets the command execution timeout
func (c *Client) SetTimeout(timeout time.Duration) {
	c.timeout = timeout
}