package scheduler

import (
	"context"
	"fmt"
	"log/slog"
)

// NamedTask struct adds a custom Name to a Task
type NamedTask struct {
	Task     Task
	TaskName string
}

// WithName is a higher-order function that adds a custom Name to a Task
func WithName(task Task, name string) Task {
	return &NamedTask{
		Task:     task,
		TaskName: name,
	}
}

func (n *NamedTask) Error() error {
	return n.Task.Error()
}

// Execute runs the Task's Execute method
func (n *NamedTask) Execute(ctx context.Context) error {
	slog.Info("Executing NamedTask", "Name", n.TaskName)
	return n.Task.Execute(ctx)
}

// Name returns the custom Name for the Task
func (n *NamedTask) Name() string {
	return n.TaskName
}

// String returns a formatted string representation of the NamedTask, including its ID.
func (n *NamedTask) String() string {
	return fmt.Sprintf("NamedTask name=%s", n.TaskName)
}

// Print logs the string representation of the NamedTask for debugging purposes.
func (n *NamedTask) Print() {
	slog.Info(n.String())
}

// GetId returns the Id of the task
func (n *NamedTask) GetId() string {
	return n.TaskName
}
