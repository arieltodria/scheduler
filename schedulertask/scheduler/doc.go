package scheduler

// Package scheduler is a library for registering and executing tasks with concurrency.

// A task is any kind of function the user wants to execute (without any input parameters).
// The Scheduler runs all tasks according to registered order. In case tasks is finished successfully it returns.

// Report struct is used to hold the tasks execution results with its errors using maps.
// The map keys are strings of the task id, and they map to the result (boolean) and to the errors (error type)
