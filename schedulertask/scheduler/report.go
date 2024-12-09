package scheduler

import (
	"errors"
	"fmt"
	"log/slog"
	"sync"
)

// Report represents the outcome of scheduled tasks.
// It holds the count of successful tasks, results indicating success or failure,
// and any errors encountered during execution.
type Report struct {
	Count       int              // Counter of total processed tasks
	ErrorsCount int              // Counter of total errors
	Errors      map[string]error // Task ID -> error (nil if successful)
	mu          sync.Mutex       // Mutex for map concurrent safety
}

// EmptyReport creates and returns a new Report instance with initialized maps for Results and Errors.
func EmptyReport() *Report {
	return &Report{
		Errors: make(map[string]error),
	}
}

// AddResult adds a task result to the report with thread safety.
// taskId is a string but can be either the Id or the name of the task.
// AddResult adds a task result to the report with thread safety.
func (r *Report) AddResult(taskId string, success bool, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.Errors[taskId] = err
	r.Count++
	if err != nil {
		r.ErrorsCount++
	}
}

// FailWithContext adds a fail error for tasks due too missing context.
func (r *Report) FailWithContext(tasks []Task) {
	for _, task := range tasks {
		r.AddResult(fmt.Sprintf("%s", task.GetId()), false, errors.New("context is nil! Can't schedule tasks"))
	}
}

// String prints the report details in a readable format.
func (r *Report) String() {
	// Acquire the lock to ensure thread-safe access to shared resources.
	r.mu.Lock()
	defer r.mu.Unlock()

	slog.Info("Scheduler Report:")
	slog.Info(fmt.Sprintf("Total Tasks: %d", r.Count))

	if len(r.Errors) == 0 {
		slog.Info("All tasks finished successfully.")
	} else {
		slog.Info(fmt.Sprintf("Total Errors: %d", len(r.Errors)))
	}

	for taskName, err := range r.Errors {
		if err == nil {
			slog.Info(fmt.Sprintf("Task %s completed successfully", taskName))
		} else {
			slog.Info(fmt.Sprintf("Task %s failed with error: %v", taskName, r.Errors[taskName]))
		}
	}
}
