package scheduler

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

// LoggingTask struct adds logging capability to a Task
type LoggingTask struct {
	Task   Task
	Logger *slog.Logger
}

// WithLogging is a higher-order function that adds logging to a Task
func WithLogging(task Task, logger *slog.Logger) Task {
	return &LoggingTask{
		Task:   task,
		Logger: logger,
	}
}

// Execute runs the Task with logging
func (l *LoggingTask) Execute(ctx context.Context) error {
	l.Logger.Info(fmt.Sprintf("Executing LoggingTask %s", l.Task.GetId()))
	start := time.Now()
	err := l.Task.Execute(ctx)
	timeElapsed := time.Since(start).Seconds()
	if err != nil {
		l.Logger.Info(fmt.Sprintf("Task %s failed after %f seconds with error: %v", l.Task.GetId(), timeElapsed, err))
	} else {
		l.Logger.Info(fmt.Sprintf("Task %s completed successfully after %f seconds", l.Task.GetId(), timeElapsed))
	}
	return err
}

// Error is a placeholder method that returns the error from the task execution.
func (l *LoggingTask) Error() error {
	return l.Task.Error()
}

// String returns a formatted string representation of the LoggingTask, including its ID.
func (l *LoggingTask) String() string {
	return fmt.Sprintf("LoggingTask %s", l.Task.GetId())
}

// Print logs the string representation of the LoggingTask for debugging purposes.
func (l *LoggingTask) Print() {
	slog.Info(l.String())
}

// GetId returns the Id of the task
func (l *LoggingTask) GetId() string {
	return l.Task.GetId()
}
