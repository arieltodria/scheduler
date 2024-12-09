package scheduler

import (
	"context"
	"fmt"
	"log/slog"
)

// RetryTask represents a task structure with an Task and retry counter
type RetryTask struct {
	Task       Task
	RetryCount int
}

// WithRetry is a higher-order function that adds a retry counter to a Task
func WithRetry(task Task, retryCount int) Task {
	return &RetryTask{
		Task:       task,
		RetryCount: retryCount,
	}
}

// Execute retries the Task up to RetryCount times if it fails
func (r *RetryTask) Execute(ctx context.Context) error {
	slog.Info("Executing RetryTask", "Name", r.Task.GetId())
	var err error
	for i := 0; i < r.RetryCount; i++ {
		err = r.Task.Execute(ctx)
		if err == nil {
			return nil
		}
		slog.Error(fmt.Sprintf("Retry %d/%d for Task %s failed with error: %v",
			i+1, r.RetryCount, r.Task.GetId(), err))
	}
	return err
}

// Error is a placeholder method that returns the error from the task execution.
func (r *RetryTask) Error() error {
	return r.Task.Error()
}

// String returns a formatted string representation of the RetryTask, including its ID.
func (r *RetryTask) String() string {
	return fmt.Sprintf("RetryTask %s count=%d", r.Task.GetId(), r.RetryCount)
}

// Print logs the string representation of the RetryTask for debugging purposes.
func (r *RetryTask) Print() {
	slog.Info(r.String())
}

// GetId returns the Id of the task
func (r *RetryTask) GetId() string {
	return r.Task.GetId()
}
