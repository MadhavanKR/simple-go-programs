package commandLineInterface

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "todo <sub-commands>",
	Short: "lets you manage your todos",
	Long:  "provides commands to add, remove, update your todo list",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("hello world from: ", cmd.Use)
	},
}

var addTodoCommand = &cobra.Command{
	Use:   "add <todo task>",
	Short: "creates a todo task",
	Long:  "creates a todo task and adds to your todo list",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		addTodo(args[0])
	},
}

var getTodosCommand = &cobra.Command{
	Use:   "list",
	Short: "lists all your todo tasks",
	Long:  "lists all your todo tasks",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		todosMap, getTodoErr := getTodos()
		if getTodoErr != nil {
			fmt.Println("error while fetching your todo list: ", getTodoErr)
		} else {
			fmt.Println("your todos are: ")
			for key, value := range todosMap {
				fmt.Printf("todos that are %s\n", key)
				for _, todo := range value {
					fmt.Printf("%d: %s\n", todo.Id, todo.Task)
				}
			}
		}
	},
}

var completeTodosCommand = &cobra.Command{
	Use:   "complete",
	Short: "mark your todo completed",
	Long:  "mark your todo completed",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		todoId, convErr := strconv.Atoi(args[0])
		if convErr != nil {
			fmt.Println("make sure you are providing the right todo id.")
			return
		}
		completeTodoErr := completeTodos(uint64(todoId))
		if completeTodoErr != nil {
			fmt.Println("error while completing your todo: ", completeTodoErr)
		} else {
			fmt.Println("tada! your todo is marked complete")
		}
	},
}

func Start() {
	rootCmd.AddCommand(addTodoCommand)
	rootCmd.AddCommand(getTodosCommand)
	rootCmd.AddCommand(completeTodosCommand)
	rootCmd.Execute()
}
