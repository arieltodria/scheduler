package main

import (
	"fmt"
	"log/slog"
	"schedulertask/schedulertask/scheduler"
)

func main() {
	// Create a concurrent task scheduler
	testScheduler := scheduler.NewConcurrentTaskScheduler()

	// Create example tasks and register them in the scheduler
	task1 := &scheduler.BaseTask{Id: "1", CustomFunc: func() error {
		return fmt.Errorf("test failure")
	}}
	task2 := &scheduler.BaseTask{Id: "2"}
	task3 := &scheduler.BaseTask{Id: "3"}
	task4 := &scheduler.BaseTask{Id: "4"}

	// Create example tasks and register them in the scheduler
	task5 := scheduler.WithRetry(scheduler.WithName(scheduler.WithLogging(task1, slog.Default()), "cheddar"), 3)
	task6 := scheduler.WithRetry(scheduler.WithName(scheduler.WithLogging(task2, slog.Default()), "cheese"), 3)
	task7 := scheduler.WithRetry(scheduler.WithName(scheduler.WithLogging(task3, slog.Default()), "is"), 3)
	task8 := scheduler.WithRetry(scheduler.WithName(scheduler.WithLogging(task4, slog.Default()), "awesome"), 3)

	testScheduler.RegisterTask(task5)
	testScheduler.RegisterTask(task6)
	testScheduler.RegisterTask(task7)
	testScheduler.RegisterTask(task8)

	// Schedule and execute the tasks concurrently
	testScheduler.ScheduleTasksWithoutInputContext()
}
