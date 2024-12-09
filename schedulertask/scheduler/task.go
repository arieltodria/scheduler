package scheduler

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
)

type Task interface {
	Execute(ctx context.Context) error // Runs the task and returns an error if it fails
	Error() error                      // Returns the last error, if any
	String() string                    // Return pretty string
	Print()                            // Print pretty
	GetId() string                     // GetId returns the Id of the task
}

// BaseTask represents a basic task structure with an ID and an optional custom function for execution.
type BaseTask struct {
	Id         string
	CustomFunc TaskFunc // Optional custom function for task execution
}

// TaskFunc represents a custom function that can be used for task execution.
// It returns an error if the execution fails.
type TaskFunc func() error

// Execute runs the task's custom function if provided, handling context cancellation and deadline exceedance.
// It also simulates failure based on the TEST_MODE environment variable.
func (t *BaseTask) Execute(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		if errors.Is(err, context.Canceled) {
			return fmt.Errorf("context canceled before executing task %s", t.Id)
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return fmt.Errorf("context deadline exceeded before executing task %s", t.Id)
		}
	}
	slog.Info("Executing BaseTask:", "BaseTask", t.Id)
	if t.CustomFunc != nil {
		err := t.CustomFunc()
		if err != nil {
			return err
		}
	}
	// Env vars that will be used in later unitests to stimulate task failures.
	if err := ctx.Err(); err != nil {
		if errors.Is(err, context.Canceled) {
			return fmt.Errorf("context canceled before executing task %s", t.Id)
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return fmt.Errorf("context deadline exceeded before executing task %s", t.Id)
		}
	}

	if os.Getenv("TEST_MODE") == "fail" {
		return fmt.Errorf("simulated failure for task %s", t.Id)
	}
	if os.Getenv("TEST_MODE") == "random" {
		randomInt := rand.Intn(2)
		if randomInt == 1 {
			return nil
		} else {
			return fmt.Errorf("simulated random failure for task %s", t.Id)
		}
	}

	return nil
}

// Error is a placeholder method that returns the error from the task execution.
func (t *BaseTask) Error() error {
	return t.Error()
}

// String returns a formatted string representation of the BaseTask, including its ID.
func (t *BaseTask) String() string {
	return fmt.Sprintf("BaseTask %s", t.Id)
}

// Print logs the string representation of the BaseTask for debugging purposes.
func (t *BaseTask) Print() {
	slog.Info(t.String())
}

// GetId returns the Id of the task
func (t *BaseTask) GetId() string {
	return t.Id
}
