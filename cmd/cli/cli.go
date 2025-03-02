package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"log"
	"os"
	database "todo/db"
	"todo/internal/repository"
	"todo/models"
)

func main() {
	// Initialize the database
	database.InitDatabase()

	// Define CLI subcommands
	getCmd := flag.NewFlagSet("get", flag.ExitOnError)
	createCmd := flag.NewFlagSet("create", flag.ExitOnError)
	updateCmd := flag.NewFlagSet("update", flag.ExitOnError)
	flag.NewFlagSet("list", flag.ExitOnError)
	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)

	// Define flags for `get`
	getID := getCmd.Int("i", 0, "Todo ID to retrieve")

	// Define flags for `create`
	createTitle := createCmd.String("t", "", "Todo title")
	createDescription := createCmd.String("d", "", "Todo description")

	// Define flags for `update`
	updateID := updateCmd.Int("i", 0, "Todo ID to update")
	updateTitle := updateCmd.String("t", "", "New title")
	updateDescription := updateCmd.String("d", "", "New description")
	updateCompleted := updateCmd.String("c", "", "Mark as completed (true/false)")

	// Define flags for `delete`
	deleteID := deleteCmd.Int("i", 0, "Todo ID to delete")

	// Ensure a command is provided
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// Handle subcommands
	switch os.Args[1] {
	case "get":
		err := getCmd.Parse(os.Args[2:])
		if err != nil {
			return
		}
		if *getID == 0 {
			log.Fatal("Error: Please provide a valid todo ID with -i")
		}
		getTodo(*getID)

	case "create":
		err := createCmd.Parse(os.Args[2:])
		if err != nil {
			return
		}
		if *createTitle == "" || *createDescription == "" {
			log.Fatal("Error: Title (-t) and description (-d) are required")
		}
		createTodo(*createTitle, *createDescription)

	case "update":
		err := updateCmd.Parse(os.Args[2:])
		if err != nil {
			return
		}
		if *updateID == 0 {
			log.Fatal("Error: Please provide a valid todo ID with -i")
		}
		updateTodo(*updateID, *updateTitle, *updateDescription, *updateCompleted)

	case "list":
		listTodos()

	case "delete":
		err := deleteCmd.Parse(os.Args[2:])
		if err != nil {
			return
		}
		if *deleteID == 0 {
			log.Fatal("Error: Please provide a valid todo ID with -i")
		}
		deleteTodo(*deleteID)

	default:
		fmt.Println("Unknown command:", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

// Function to GET a todo directly from SQLite
func getTodo(id int) {
	ctx := context.Background()
	todo, err := database.Database.Queries.GetTodo(ctx, int64(id))
	if err != nil {
		log.Fatalf("Error fetching todo: %v", err)
	}
	printTable([]models.TodoResponse{sanitiseTodo(todo)})
}

// Function to CREATE a todo directly in SQLite
func createTodo(title, description string) {
	params := repository.CreateTodoParams{
		Title:       title,
		Description: description,
	}
	ctx := context.Background()
	todo, err := database.Database.Queries.CreateTodo(ctx, params)
	if err != nil {
		log.Fatalf("Error creating todo: %v", err)
	}
	fmt.Println("Todo Created Successfully!")
	printTable([]models.TodoResponse{sanitiseTodo(todo)})
}

// Function to UPDATE a todo directly in SQLite
func updateTodo(id int, title, description string, completedStr string) {
	ctx := context.Background()

	existingTodo, err := database.Database.Queries.GetTodo(ctx, int64(id))
	if err != nil {
		log.Fatalf("Error fetching existing todo: %v", err)
	}

	// Keep existing title/description if not provided
	if title == "" {
		title = existingTodo.Title
	}
	if description == "" {
		description = existingTodo.Description
	}

	// Convert completed flag
	completed := existingTodo.Completed // Default to existing value
	if completedStr == "true" {
		completed = true
	} else if completedStr == "false" {
		completed = false
	}

	params := repository.UpdateTodoParams{
		ID:          int64(id),
		Title:       title,
		Description: description,
		Completed:   completed,
	}

	// Perform the update
	todo, err := database.Database.Queries.UpdateTodo(ctx, params)
	if err != nil {
		log.Fatalf("Error updating todo: %v", err)
	}

	fmt.Println("Todo Updated Successfully!")
	printTable([]models.TodoResponse{sanitiseTodo(todo)})
}

// Function to LIST all todos directly from SQLite
func listTodos() {
	ctx := context.Background()
	todos, err := database.Database.Queries.ListTodos(ctx)
	if err != nil {
		log.Fatalf("Error listing todos: %v", err)
	}
	printTable(sanitiseTodos(todos))
}

// Function to DELETE a todo directly from SQLite
func deleteTodo(id int) {
	ctx := context.Background()
	err := database.Database.Queries.DeleteTodo(ctx, int64(id))
	if err != nil {
		log.Fatalf("Error deleting todo: %v", err)
	}
	fmt.Printf("Todo with ID %d successfully deleted\n", id)
}

// Function to display todos in a table format
func printTable(todos []models.TodoResponse) {
	if len(todos) == 0 {
		fmt.Println("No todos found.")
		return
	}

	// Create a new table writer
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Title", "Description", "Completed"})

	// Populate table rows
	for _, todo := range todos {
		completedText := "No"
		if todo.Completed {
			completedText = "Yes"
		}

		row := []string{
			fmt.Sprintf("%d", todo.ID),
			todo.Title,
			todo.Description,
			completedText,
		}
		table.Append(row)
	}

	// Customize the table style
	table.SetBorder(true)
	table.SetRowLine(true)
	table.SetAutoWrapText(true)
	table.Render() // Render the table
}

// Function to sanitize database todo to response format
func sanitiseTodo(todo repository.Todo) models.TodoResponse {
	return models.TodoResponse{
		ID:          todo.ID,
		Title:       todo.Title,
		Description: todo.Description,
		Completed:   todo.Completed,
	}
}

func sanitiseTodos(todos []repository.Todo) []models.TodoResponse {
	response := make([]models.TodoResponse, len(todos))
	for i, todo := range todos {
		response[i] = sanitiseTodo(todo)
	}
	return response
}

// Function to display CLI usage instructions
func printUsage() {
	fmt.Println("Usage: todo [command] [options]")
	fmt.Println("\nCommands:")
	fmt.Println("  get      -i <ID>                  Get a specific todo")
	fmt.Println("  create   -t <title> -d <desc>      Create a new todo")
	fmt.Println("  update   -i <ID> -t <title> -d <desc> -c <true/false>  Update a todo")
	fmt.Println("  list                               List all todos")
	fmt.Println("  delete   -i <ID>                   Delete a todo")
}
