package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("$ ")
		if !scanner.Scan() {
			break
		}
		input := scanner.Text()

		// Parse multiple commands separated by semicolons
		commands := strings.Split(input, ";")
		for _, command := range commands {
			executeCommand(strings.TrimSpace(command))
		}
	}
}

func executeCommand(input string) {
	// Parse pipe operations
	pipes := strings.Split(input, "|")

	if len(pipes) > 1 {
		executePipedCommands(pipes)
		return
	}

	// Parse command and arguments
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return
	}
	cmd := parts[0]
	args := parts[1:]

	switch cmd {
	case "exit":
		os.Exit(0)
	case "cd":
		if len(args) == 0 {
			fmt.Println("cd: missing directory")
			return
		}
		err := os.Chdir(args[0])
		if err != nil {
			fmt.Println("cd:", err)
		}
	case "history":
		// Implement history functionality
		fmt.Println("History")
	case "ls":
		dir := "."
		if len(args) > 0 {
			dir = args[0]
		}
		files, err := os.ReadDir(dir)
		if err != nil {
			fmt.Println("ls:", err)
			return
		}
		for _, file := range files {
			fmt.Println(file.Name())
		}
	case "pwd":
		wd, err := os.Getwd()
		if err != nil {
			fmt.Println("pwd:", err)
			return
		}
		fmt.Println(wd)
	case "mkdir":
		if len(args) == 0 {
			fmt.Println("mkdir: missing directory name")
			return
		}
		err := os.Mkdir(args[0], 0755)
		if err != nil {
			fmt.Println("mkdir:", err)
		}
	case "rm":
		if len(args) == 0 {
			fmt.Println("rm: missing file name")
			return
		}
		for _, file := range args {
			err := os.Remove(file)
			if err != nil {
				fmt.Println("rm:", err)
			}
		}
	case "cat":
		if len(args) == 0 {
			fmt.Println("cat: missing file name")
			return
		}
		for _, file := range args {
			data, err := os.ReadFile(file)
			if err != nil {
				fmt.Println("cat:", err)
				continue
			}
			fmt.Println(string(data))
		}
	case "echo":
		fmt.Println(strings.Join(args, " "))
	case "date":
		fmt.Println(time.Now().Format("Mon Jan _2 15:04:05 MST 2006"))
	case "whoami":
		fmt.Println(os.Getenv("USER"))
	case "env":
		for _, env := range os.Environ() {
			fmt.Println(env)
		}
	case "clear":
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	case "help":
		fmt.Println("Available commands:")
		fmt.Println("  exit   - Exit the shell")
		fmt.Println("  cd     - Change directory")
		fmt.Println("  history- Show command history")
		fmt.Println("  ls     - List files in directory")
		fmt.Println("  pwd    - Print current directory")
		fmt.Println("  mkdir  - Create a directory")
		fmt.Println("  rm     - Remove file(s)")
		fmt.Println("  cat    - Concatenate and display file(s)")
		fmt.Println("  echo   - Display message")
		fmt.Println("  date   - Print current date and time")
		fmt.Println("  whoami - Print current user")
		fmt.Println("  env    - Print environment variables")
		fmt.Println("  clear  - Clear the screen")
		fmt.Println("  help   - Display this help message")
	default:
		// Execute external command
		cmd := exec.Command(cmd, args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		err := cmd.Run()
		if err != nil {
			fmt.Println("Error:", err)
		}
	}
}

func executePipedCommands(pipes []string) {
	var cmd *exec.Cmd
	var err error

	for _, pipe := range pipes {
		parts := strings.Fields(pipe)
		if len(parts) == 0 {
			continue
		}
		cmd = exec.Command(parts[0], parts[1:]...)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		if err = cmd.Start(); err != nil {
			fmt.Println("Error:", err)
			return
		}
	}
	if err = cmd.Wait(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if status, ok := exitError.Sys().(syscall.WaitStatus); ok {
				fmt.Printf("Error: Command exited with status %d\n", status.ExitStatus())
			}
		} else {
			fmt.Println("Error:", err)
		}
	}
}
