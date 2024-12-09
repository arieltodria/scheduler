package scheduler

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"strings"
	"sync"
	"testing"
	"time"
)

var defaultLogger = slog.Default()
var defaultTestNames = []string{"Cheddar", "Camembert"}

// SuccessfulTaskFunc simulates a successful task function (returns nil).
func SuccessfulTaskFunc() error {
	return nil
}

type TestCase struct {
	name           string // Name of the test case.
	tasks          []Task // List of tasks to register.
	wantCount      int    // Expected task count after registration.
	expectFailures []bool // Expected outcome for each task execution (true for failure, false for success).
}

// TestRegisterTask uses multiple test cases to verify all tasks are registered correctly to the scheduler's tasks slice.
func TestRegisterTask(t *testing.T) {
	testCases := []TestCase{
		{
			name:      "No Task",
			tasks:     []Task{},
			wantCount: 0,
		},
		{
			name: "Single Task",
			tasks: []Task{
				&BaseTask{Id: "1", CustomFunc: SuccessfulTaskFunc},
			},
			wantCount: 1,
		},
		{
			name: "Multiple Tasks",
			tasks: []Task{
				&BaseTask{Id: "1", CustomFunc: SuccessfulTaskFunc},
				&BaseTask{Id: "2", CustomFunc: SuccessfulTaskFunc},
			},
			wantCount: 2,
		},
		{
			name: "Logging Task",
			tasks: []Task{
				WithLogging(&BaseTask{Id: "1", CustomFunc: SuccessfulTaskFunc}, defaultLogger),
			},
			wantCount: 1,
		},
		{
			name: "Named Task",
			tasks: []Task{
				WithName(&BaseTask{Id: "1", CustomFunc: SuccessfulTaskFunc}, defaultTestNames[0]),
			},
			wantCount: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testScheduler := NewConcurrentTaskScheduler()

			// Register tasks and check if the task count is as expected.
			testScheduler.RegisterTasks(tc.tasks)
			if got := testScheduler.NumOfTasks(); got != tc.wantCount {
				t.Fatalf("Num of registered task is %d but expected %d", got, tc.wantCount)
			}
		})
	}
}

// FailingTaskFunc simulates a failing task function (returns an error).
func FailingTaskFunc() error {
	return errors.New("simulated task failure")
}

func TestRegisterAndExecuteTasks(t *testing.T) {
	testCases := []TestCase{
		{
			name:           "No Task",
			tasks:          []Task{},
			wantCount:      0,
			expectFailures: []bool{},
		},
		{
			name: "Single Task- Success",
			tasks: []Task{
				&BaseTask{Id: "1", CustomFunc: SuccessfulTaskFunc},
			},
			wantCount:      1,
			expectFailures: []bool{false},
		},
		{
			name: "Single Task- Failure",
			tasks: []Task{
				&BaseTask{Id: "1", CustomFunc: FailingTaskFunc},
			},
			wantCount:      1,
			expectFailures: []bool{true},
		},
		{
			name: "Multiple Tasks- All Success",
			tasks: []Task{
				&BaseTask{Id: "1", CustomFunc: SuccessfulTaskFunc},
				&BaseTask{Id: "2", CustomFunc: SuccessfulTaskFunc},
			},
			wantCount:      2,
			expectFailures: []bool{false, false},
		},
		{
			name: "Multiple Tasks- All Failure",
			tasks: []Task{
				&BaseTask{Id: "1", CustomFunc: FailingTaskFunc},
				&BaseTask{Id: "2", CustomFunc: FailingTaskFunc},
			},
			wantCount:      2,
			expectFailures: []bool{true, true},
		},
		{
			name: "Multiple Tasks- First Success, Then Failure",
			tasks: []Task{
				&BaseTask{Id: "1", CustomFunc: SuccessfulTaskFunc},
				&BaseTask{Id: "2", CustomFunc: FailingTaskFunc},
			},
			wantCount:      2,
			expectFailures: []bool{false, true},
		},
		{
			name: "Multiple Tasks- First Failure, Then Success",
			tasks: []Task{
				&BaseTask{Id: "1", CustomFunc: FailingTaskFunc},
				&BaseTask{Id: "2", CustomFunc: SuccessfulTaskFunc},
			},
			wantCount:      2,
			expectFailures: []bool{true, false},
		},
		{
			name: "Logged Task - Success",
			tasks: []Task{
				WithLogging(&BaseTask{Id: "1", CustomFunc: SuccessfulTaskFunc}, defaultLogger),
			},
			wantCount:      1,
			expectFailures: []bool{false},
		},
		{
			name: "Logged Task - Failure",
			tasks: []Task{
				WithLogging(&BaseTask{Id: "2", CustomFunc: FailingTaskFunc}, defaultLogger),
			},
			wantCount:      1,
			expectFailures: []bool{true},
		},
		{
			name: "Named Task - Success",
			tasks: []Task{
				WithName(&BaseTask{Id: "1", CustomFunc: SuccessfulTaskFunc}, defaultTestNames[0]),
			},
			wantCount:      1,
			expectFailures: []bool{false},
		},
		{
			name: "Named Task - Failure",
			tasks: []Task{
				WithName(&BaseTask{Id: "2", CustomFunc: FailingTaskFunc}, defaultTestNames[0]),
			},
			wantCount:      1,
			expectFailures: []bool{true},
		},
		{
			name: "Named Logging Task - Failure",
			tasks: []Task{
				WithLogging(
					WithName(&BaseTask{Id: "2", CustomFunc: FailingTaskFunc}, defaultTestNames[0]),
					defaultLogger,
				),
			},
			wantCount:      1,
			expectFailures: []bool{true},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testScheduler := NewConcurrentTaskScheduler()

			// Register the tasks.
			testScheduler.RegisterTasks(tc.tasks)
			if got := testScheduler.NumOfTasks(); got != tc.wantCount {
				t.Fatalf("Expected num of tasks %d but got %d", got, tc.wantCount)
			}

			// Create a slice to store task errors or nil results.
			results := make([]error, len(tc.tasks))

			// Override ScheduleTasks() method to capture errors for test verification.
			var wg sync.WaitGroup
			for i, task := range tc.tasks {
				wg.Add(1)
				go func(i int, task Task) {
					defer wg.Done()
					results[i] = task.Execute(context.Background())
				}(i, task)
			}
			wg.Wait()

			// Validate the execution result matches the expected outcome for each task.
			for i, result := range results {
				gotFailure := result != nil
				if gotFailure != tc.expectFailures[i] {
					t.Fatalf("Task %s execution result=%v; expected failure=%v", tc.tasks[i].GetId(), gotFailure, tc.expectFailures[i])
				}
			}
		})
	}
}

// SleepSuccessTaskFunc simulates a task that takes time to succeed.
func SleepSuccessTaskFunc(duration time.Duration) TaskFunc {
	return func() error {
		time.Sleep(duration)
		return nil
	}
}

// SleepFailTaskFunc simulates a task that takes time and then fails.
func SleepFailTaskFunc(duration time.Duration) TaskFunc {
	return func() error {
		time.Sleep(duration)
		return errors.New("simulated task failure after sleep")
	}
}

func TestConcurrency(t *testing.T) {
	testCases := []TestCase{
		{
			name: "Concurrent Execution with Sleep",
			tasks: []Task{
				&BaseTask{Id: "1", CustomFunc: SleepSuccessTaskFunc(2 * time.Second)},
				&BaseTask{Id: "2", CustomFunc: SleepFailTaskFunc(1 * time.Second)},
			},
			wantCount:      2,
			expectFailures: []bool{false, true},
		},
		{
			name: "Success and Wait for Failing Task",
			tasks: []Task{
				&BaseTask{Id: "1", CustomFunc: SleepSuccessTaskFunc(1 * time.Second)},
				&BaseTask{Id: "2", CustomFunc: SleepFailTaskFunc(2 * time.Second)},
			},
			wantCount:      2,
			expectFailures: []bool{false, true},
		},
		{
			name: "Success Wait for Another Success",
			tasks: []Task{
				&BaseTask{Id: "1", CustomFunc: SleepSuccessTaskFunc(1 * time.Second)},
				&BaseTask{Id: "2", CustomFunc: SleepSuccessTaskFunc(2 * time.Second)},
			},
			wantCount:      2,
			expectFailures: []bool{false, false},
		},
		{
			name: "Fail and Wait for Another Fail",
			tasks: []Task{
				&BaseTask{Id: "1", CustomFunc: SleepFailTaskFunc(1 * time.Second)},
				&BaseTask{Id: "2", CustomFunc: SleepFailTaskFunc(2 * time.Second)},
			},
			wantCount:      2,
			expectFailures: []bool{true, true},
		},
		{
			name: "Fail and Wait for Success",
			tasks: []Task{
				&BaseTask{Id: "1", CustomFunc: SleepFailTaskFunc(1 * time.Second)},
				&BaseTask{Id: "2", CustomFunc: SleepSuccessTaskFunc(2 * time.Second)},
			},
			wantCount:      2,
			expectFailures: []bool{true, false},
		},
		{
			name: "Sleep Success and Fail",
			tasks: []Task{
				&BaseTask{Id: "1", CustomFunc: SleepSuccessTaskFunc(2 * time.Second)},
				&BaseTask{Id: "2", CustomFunc: FailingTaskFunc},
			},
			wantCount:      2,
			expectFailures: []bool{false, true},
		},
		{
			name: "Sleep Fail and Success",
			tasks: []Task{
				&BaseTask{Id: "1", CustomFunc: SleepFailTaskFunc(2 * time.Second)},
				&BaseTask{Id: "2", CustomFunc: SuccessfulTaskFunc},
			},
			wantCount:      2,
			expectFailures: []bool{true, false},
		},
		{
			name: "Logged Task Success and Failure",
			tasks: []Task{
				WithLogging(&BaseTask{Id: "1", CustomFunc: SleepSuccessTaskFunc(2 * time.Second)}, defaultLogger),
				WithLogging(&BaseTask{Id: "2", CustomFunc: SleepFailTaskFunc(1 * time.Second)}, defaultLogger),
			},
			wantCount:      2,
			expectFailures: []bool{false, true},
		},
		{
			name: "Named Logged Task Success and Failure",
			tasks: []Task{
				WithLogging(
					WithName(&BaseTask{Id: "1", CustomFunc: SleepSuccessTaskFunc(1 * time.Second)}, defaultTestNames[0]),
					defaultLogger),
				WithLogging(
					WithName(&BaseTask{Id: "2", CustomFunc: SleepFailTaskFunc(2 * time.Second)}, defaultTestNames[0]),
					defaultLogger),
			},
			wantCount:      2,
			expectFailures: []bool{false, true},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testScheduler := NewConcurrentTaskScheduler()

			// Register the tasks
			testScheduler.RegisterTasks(tc.tasks)
			if got := testScheduler.NumOfTasks(); got != tc.wantCount {
				t.Fatalf("Got %d tasks but expected %d", got, tc.wantCount)
			}

			// Create a slice to hold task errors or nil results
			results := make([]error, len(tc.tasks))

			var wg sync.WaitGroup
			for i, task := range tc.tasks {
				wg.Add(1)
				go func(i int, task Task) {
					defer wg.Done()
					results[i] = task.Execute(context.Background())
				}(i, task)
			}
			wg.Wait()

			// Validate the execution result matches the expected outcome for each task
			for i, result := range results {
				gotFailure := result != nil
				if gotFailure != tc.expectFailures[i] {
					t.Fatalf("Task %s execution=%v but expected failure=%v", tc.tasks[i].GetId(), gotFailure, tc.expectFailures[i])
				}
			}
		})
	}
}

func TestReportAddResult(t *testing.T) {
	report := EmptyReport()

	// Mock tasks
	tasks := []Task{
		&BaseTask{Id: "1", CustomFunc: SuccessfulTaskFunc},
		&BaseTask{Id: "2", CustomFunc: FailingTaskFunc},
		&BaseTask{Id: "3", CustomFunc: SuccessfulTaskFunc},
	}
	expectedTotalTasks := len(tasks) // Total tasks processed
	expectedFailures := 1            // One task is expected to fail

	// Process each task and add the result to the report.
	for _, task := range tasks {
		err := task.Execute(context.Background())
		report.AddResult(task.GetId(), err == nil, err)
	}

	// 1. Verify that Counter holds the total number of processed tasks.
	if report.Count != expectedTotalTasks {
		t.Fatalf("Expected total count of processed tasks to be %d but got %d", expectedTotalTasks, report.Count)
	}

	// 2. Verify the length of the Errors map matches the total number of processed tasks.
	if len(report.Errors) != expectedTotalTasks {
		t.Fatalf("Expected errors map length of %d but got %d", expectedTotalTasks, len(report.Errors))
	}

	// 3. Verify that ErrorsCount holds the total number of failed tasks.
	if report.ErrorsCount != expectedFailures {
		t.Fatalf("Expected total errors of %d but got %d", expectedFailures, report.ErrorsCount)
	}
}

// CaptureLogs sets up a logger that writes to a buffer and returns the logger and a function to retrieve logs.
func CaptureLogs() (*slog.Logger, func() string) {
	var buf bytes.Buffer
	handler := slog.NewJSONHandler(&buf, nil)
	logger := slog.New(handler)

	// Function to retrieve log output
	getLogs := func() string {
		return buf.String()
	}

	return logger, getLogs
}

func TestWithLogging(t *testing.T) {
	// Start capturing logs
	testLogger, getLogs := CaptureLogs()

	// Define a task that will fail
	mockTask := &BaseTask{Id: "10", CustomFunc: FailingTaskFunc}

	// Wrap the task with logging using the test logger
	loggedTask := WithLogging(mockTask, testLogger)

	// Execute the task
	err := loggedTask.Execute(context.Background())
	if err == nil {
		t.Fatal("Expected an error, but got none")
	}

	// Get the captured log output
	logOutput := getLogs()

	if logOutput == "" {
		t.Fatalf("Failed to capture logged task output")
	}
	if !strings.Contains(logOutput, loggedTask.GetId()) {
		t.Fatalf("Expected task ID in log")
	}
	if !strings.Contains(logOutput, "after") || !strings.Contains(logOutput, "seconds") {
		t.Fatalf("Expected elapsed time to be printed in the task logs")
	}
	if !strings.Contains(logOutput, "Executing LoggingTask") {
		t.Fatalf("Expected 'Executing task' log entry, but it was missing")
	}
	if !strings.Contains(logOutput, "failed") {
		t.Fatalf("Expected 'failed' log entry, but it was missing")
	}
}

func TestWithNameWrapper(t *testing.T) {
	// Initialize a new Scheduler
	testScheduler := NewConcurrentTaskScheduler()

	// Register tasks
	// Define named tasks
	namedTask := WithName(&BaseTask{Id: "42", CustomFunc: SuccessfulTaskFunc}, defaultTestNames[0]) // Expected to succeed
	failTask := WithName(&BaseTask{Id: "43", CustomFunc: FailingTaskFunc}, defaultTestNames[1])     // Expected to fail

	// Register tasks
	testScheduler.RegisterTask(namedTask)
	testScheduler.RegisterTask(failTask)

	testScheduler.ScheduleTasksWithoutInputContext()

	// Check if the Report has the correct count of processed tasks
	if testScheduler.Report.Count != 2 {
		t.Fatalf("Expected 2 total processed tasks, got %d", testScheduler.Report.Count)
	}

	// Verify the number of errors
	if testScheduler.Report.ErrorsCount != 1 {
		t.Fatalf("Expected 1 failed task, got %d", testScheduler.Report.ErrorsCount)
	}

	// Expected results to be stored with their respective task names
	expectedResultsKeys := map[string]bool{
		defaultTestNames[0]: true,
		defaultTestNames[1]: true,
	}

	// Check that all expected task names are present in the Errors map, regardless of success or failure
	for key := range testScheduler.Report.Errors {
		if !expectedResultsKeys[key] {
			t.Fatalf("Unexpected key in report: %s", key)
		}
	}

	// Check that the failed task is recorded properly
	if _, found := testScheduler.Report.Errors[defaultTestNames[1]]; !found {
		t.Fatalf("Expected error entry for task named %s", defaultTestNames[1])
	}
	if err := testScheduler.Report.Errors[defaultTestNames[0]]; err != nil {
		t.Fatalf("Expected no error for task named %s", defaultTestNames[0])
	}
}

func TestContextHandling(t *testing.T) {
	testCases := []struct {
		name           string
		tasks          []Task
		contextAction  func(ctx context.Context, cancel context.CancelFunc)
		wantExecutions int // how many tasks should successfully execute
	}{
		{
			name: "Context Cancellation",
			tasks: []Task{
				&BaseTask{Id: "1", CustomFunc: SleepSuccessTaskFunc(2 * time.Second)},
				&BaseTask{Id: "2", CustomFunc: SleepFailTaskFunc(2 * time.Second)},
			},
			contextAction: func(ctx context.Context, cancel context.CancelFunc) {
				time.Sleep(1 * time.Second)
				cancel()
			},
			wantExecutions: 0,
		},
		{
			name: "Context Deadline",
			tasks: []Task{
				&BaseTask{Id: "1", CustomFunc: SleepSuccessTaskFunc(3 * time.Second)},
			},
			contextAction: func(ctx context.Context, cancel context.CancelFunc) {
				// Deadline is ahead of task sleep period
			},
			wantExecutions: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testScheduler := NewConcurrentTaskScheduler()
			testScheduler.RegisterTasks(tc.tasks)

			ctx, cancel := context.WithCancel(context.Background())
			if tc.name == "Context Deadline" {
				ctx, cancel = context.WithTimeout(context.Background(), 1*time.Second)
			}
			defer cancel()

			// simulate context action (cancel or wait for deadline)
			if tc.contextAction != nil {
				go tc.contextAction(ctx, cancel)
			}

			// Schedule Tasks using input context
			testScheduler.ScheduleTasksWithInputContext(ctx)

			// Check results
			executionCount := 0
			for _, err := range testScheduler.Report.Errors {
				if err == nil {
					executionCount++
				}
			}
			if executionCount != tc.wantExecutions {
				t.Errorf("Expected %d successful executions but got %d", tc.wantExecutions, executionCount)
			}
		})
	}
}
