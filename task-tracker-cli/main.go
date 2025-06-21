package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Fadhelalfin/task-tracker-cli/task" // <<< PASTIKAN INI SESUAI DENGAN go.mod Anda
)

func main() {
	// Definisikan perintah (subkomando)
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "add":
		handleAdd()
	case "list":
		handleList()
	case "update":
		handleUpdate()
	case "progress":
		handleProgress()
	case "complete":
		handleComplete()
	case "delete":
		handleDelete()
	case "help":
		printUsage()
	default:
		fmt.Printf("Error: Unknown command '%s'\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: task-tracker <command> [arguments]")
	fmt.Println("\nCommands:")
	fmt.Println("  add <description>           Add a new task")
	fmt.Println("  list                        List all tasks")
	fmt.Println("  update <id> <new_description> Update task description")
	fmt.Println("  progress <id>               Mark task as in-progress")
	fmt.Println("  complete <id>               Mark task as completed")
	fmt.Println("  delete <id>                 Delete a task")
	fmt.Println("  help                        Show this help message")
	fmt.Println("\nExamples:")
	fmt.Println("  task-tracker add 'Buy groceries'")
	fmt.Println("  task-tracker list")
	fmt.Println("  task-tracker update 1 'Buy milk and eggs'")
	fmt.Println("  task-tracker progress 2")
	fmt.Println("  task-tracker complete 3")
	fmt.Println("  task-tracker delete 1")
}

func handleAdd() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: task-tracker add <description>")
		os.Exit(1)
	}
	description := strings.Join(os.Args[2:], " ")

	tasks, err := task.LoadTasks()
	if err != nil {
		log.Fatalf("Error loading tasks: %v", err)
	}

	newTask := task.Task{
		ID:          task.GetNextID(tasks),
		Description: description,
		Status:      task.StatusPending,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	tasks = append(tasks, newTask)

	if err := task.SaveTasks(tasks); err != nil {
		log.Fatalf("Error saving tasks: %v", err)
	}
	fmt.Printf("Task added successfully! ID: %d\n", newTask.ID)
}

func handleList() {
	tasks, err := task.LoadTasks()
	if err != nil {
		log.Fatalf("Error loading tasks: %v", err)
	}

	if len(tasks) == 0 {
		fmt.Println("No tasks found. Add one with 'task-tracker add <description>'")
		return
	}

	fmt.Println("Your Tasks:")
	for _, t := range tasks {
		fmt.Printf("ID: %d | Status: %s | Description: %s | Created: %s\n",
			t.ID, strings.ToUpper(t.Status), t.Description, t.CreatedAt.Format("2006-01-02 15:04"))
	}
}

func getTaskIDFromArgs() (int, error) {
	if len(os.Args) < 3 {
		return 0, fmt.Errorf("missing task ID")
	}
	id, err := strconv.Atoi(os.Args[2])
	if err != nil {
		return 0, fmt.Errorf("invalid task ID: %v", err)
	}
	return id, nil
}

func findTaskByID(tasks []task.Task, id int) *task.Task {
	for i := range tasks {
		if tasks[i].ID == id {
			return &tasks[i] // Mengembalikan pointer agar bisa diupdate
		}
	}
	return nil
}

func handleUpdate() {
	id, err := getTaskIDFromArgs()
	if err != nil {
		fmt.Printf("Usage: task-tracker update <id> <new_description>\nError: %v\n", err)
		os.Exit(1)
	}
	if len(os.Args) < 4 {
		fmt.Println("Usage: task-tracker update <id> <new_description>")
		os.Exit(1)
	}
	newDescription := strings.Join(os.Args[3:], " ")

	tasks, err := task.LoadTasks()
	if err != nil {
		log.Fatalf("Error loading tasks: %v", err)
	}

	foundTask := findTaskByID(tasks, id)
	if foundTask == nil {
		fmt.Printf("Task with ID %d not found.\n", id)
		os.Exit(1)
	}

	foundTask.Description = newDescription
	foundTask.UpdatedAt = time.Now()

	if err := task.SaveTasks(tasks); err != nil {
		log.Fatalf("Error saving tasks: %v", err)
	}
	fmt.Printf("Task %d updated successfully.\n", id)
}

func handleProgress() {
	id, err := getTaskIDFromArgs()
	if err != nil {
		fmt.Printf("Usage: task-tracker progress <id>\nError: %v\n", err)
		os.Exit(1)
	}

	tasks, err := task.LoadTasks()
	if err != nil {
		log.Fatalf("Error loading tasks: %v", err)
	}

	foundTask := findTaskByID(tasks, id)
	if foundTask == nil {
		fmt.Printf("Task with ID %d not found.\n", id)
		os.Exit(1)
	}

	if foundTask.Status == task.StatusCompleted {
		fmt.Printf("Task %d is already completed. Cannot set to in-progress.\n", id)
		os.Exit(1)
	}
	foundTask.Status = task.StatusInProgress
	foundTask.UpdatedAt = time.Now()

	if err := task.SaveTasks(tasks); err != nil {
		log.Fatalf("Error saving tasks: %v", err)
	}
	fmt.Printf("Task %d marked as 'in-progress'.\n", id)
}

func handleComplete() {
	id, err := getTaskIDFromArgs()
	if err != nil {
		fmt.Printf("Usage: task-tracker complete <id>\nError: %v\n", err)
		os.Exit(1)
	}

	tasks, err := task.LoadTasks()
	if err != nil {
		log.Fatalf("Error loading tasks: %v", err)
	}

	foundTask := findTaskByID(tasks, id)
	if foundTask == nil {
		fmt.Printf("Task with ID %d not found.\n", id)
		os.Exit(1)
	}

	if foundTask.Status == task.StatusCompleted {
		fmt.Printf("Task %d is already completed.\n", id)
		os.Exit(0)
	}
	foundTask.Status = task.StatusCompleted
	foundTask.UpdatedAt = time.Now()

	if err := task.SaveTasks(tasks); err != nil {
		log.Fatalf("Error saving tasks: %v", err)
	}
	fmt.Printf("Task %d marked as 'completed'.\n", id)
}

func handleDelete() {
	id, err := getTaskIDFromArgs()
	if err != nil {
		fmt.Printf("Usage: task-tracker delete <id>\nError: %v\n", err)
		os.Exit(1)
	}

	tasks, err := task.LoadTasks()
	if err != nil {
		log.Fatalf("Error loading tasks: %v", err)
	}

	var updatedTasks []task.Task
	found := false
	for _, t := range tasks {
		if t.ID == id {
			found = true
		} else {
			updatedTasks = append(updatedTasks, t)
		}
	}

	if !found {
		fmt.Printf("Task with ID %d not found.\n", id)
		os.Exit(1)
	}

	if err := task.SaveTasks(updatedTasks); err != nil {
		log.Fatalf("Error saving tasks: %v", err)
	}
	fmt.Printf("Task %d deleted successfully.\n", id)
}