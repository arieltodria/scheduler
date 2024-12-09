package scheduler

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
)

// Scheduler is a struct that manages and executes task concurrently.
// It holds a channel for tasks, a list of registered tasks, and a wait group
// to synchronize goroutines.
type Scheduler struct {
	tasks  []Task
	Report *Report
}

// NewConcurrentTaskScheduler creates a new Scheduler with an initialized Report.
func NewConcurrentTaskScheduler() *Scheduler {
	return &Scheduler{
		Report: EmptyReport(),
	}
}

// NumOfTasks returns the num of registered tasks.
func (s *Scheduler) NumOfTasks() int {
	return len(s.tasks)
}

// RegisterTasks adds multiple tasks to the scheduler.
func (s *Scheduler) RegisterTasks(tasks []Task) {
	slog.Info("Registering Tasks")
	for _, t := range tasks {
		s.RegisterTask(t)
	}
}

// RegisterTask adds a task to the scheduler.
func (s *Scheduler) RegisterTask(t Task) {
	slog.Info(fmt.Sprintf("Registering ExampleTask %s", t.GetId()))
	s.tasks = append(s.tasks, t)
}

// ScheduleTasksWithInputContext schedules tasks with context from input
func (s *Scheduler) ScheduleTasksWithInputContext(ctx context.Context) {
	if ctx.Err() != nil {
		slog.Error("Context is nil! Can't schedule tasks")
		s.Report.FailWithContext(s.tasks)
		return
	}
	s.ScheduleTasks(ctx)
}

// ScheduleTasksWithoutInputContext schedules tasks with generated local context
func (s *Scheduler) ScheduleTasksWithoutInputContext() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if ctx.Err() != nil {
		slog.Error("Context is nil! Can't schedule tasks")
		s.Report.FailWithContext(s.tasks)
		return
	}
	s.ScheduleTasks(ctx)
}

// ScheduleTasks schedules all registered tasks for execution.
func (s *Scheduler) ScheduleTasks(ctx context.Context) {
	var wg sync.WaitGroup
	wg.Add(len(s.tasks))

	for _, t := range s.tasks {
		go func(t Task) {
			defer wg.Done()

			// Execute task.
			if err := t.Execute(ctx); err != nil {
				slog.Error("failed to execute", "task", t.GetId())
				s.Report.AddResult(t.GetId(), false, err)
			} else {
				slog.Info("Task completed", "task", t.GetId())
				s.Report.AddResult(t.GetId(), true, nil)
			}
		}(t)
	}

	// Wait for all tasks to complete.
	wg.Wait()
	// Print report.
	s.Report.String()
}
