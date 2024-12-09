Exercise: Concurrent Task Scheduler

You are tasked with building a concurrent task scheduler using Go. The scheduler should allow users to register tasks, schedule them for execution, and execute them concurrently. Additionally, the scheduler should support higher-order functions and leverage functional programming concepts.

Requirements:
1. Create a **`Task`** interface. 
      - Try to think what are the main functions for this interface and what is their purpose.
2. Implement the following functionality:
   - The scheduler should support registering tasks. Tasks should implement the **`Task`** interface.
   - The scheduler should allow users to schedule tasks for execution.
   - The scheduler should execute tasks concurrently, using goroutines and channels.
   - Implement a higher-order function, **`WithRetry`**, that takes a task and a retry count as arguments and returns a new task that retries the original task a specified number of times if it fails.
   - Implement a higher-order function, **`WithName`**, that takes a task and a name as arguments and returns a new task with the provided name.
   - Implement a higher-order function, **`WithLogging`**, that takes a task and adds logging capabilities to it.
   - Implement a functional programming concept to operate on a collection of tasks.

3. Write unit tests to ensure the correctness of your implementation, including tests for concurrency, higher-order functions, and functional programming concepts.
4. Document your code and provide clear instructions on how to use the scheduler and higher-order functions.

Feel free to use any additional Go packages or frameworks that you're familiar with. When you're finished, test your scheduler using different tasks, including retryable tasks, tasks with names, and tasks with logging capabilities. Verify that the scheduler executes the tasks concurrently and applies the higher-order functions correctly.

Note: This exercise aims to challenge your knowledge and skills in working with concurrency, higher-order functions, functional programming, and interfaces. Feel free to enhance the exercise or add more complex features to demonstrate your abilities further.

Good luck with the exercise! If you have any specific questions or need further assistance, feel free to ask.