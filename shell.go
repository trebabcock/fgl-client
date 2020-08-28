package main

import (
	"fmt"
	"strings"

	"github.com/c-bata/go-prompt"
	. "github.com/logrusorgru/aurora"
)

func shell() {
	clear()
	for {
		t := prompt.Input(makeShellPrompt(), completer, prompt.OptionPrefixTextColor(prompt.Red))
		history.Add(t)
		execInput(t)
	}
}

func execInput(input string) {
	input = strings.TrimSuffix(input, "\n")
	args := strings.Split(input, " ")

	switch args[0] {
	case "cd":
		if len(args) < 2 {
			fmt.Println("cd: path required")
		}
		// not implemented yet
	case "exit":
		mainMenu("")
	case "clear":
		clear()
	case "pwd":
		fmt.Println(currentDir.Name)
	case "ls":
		for _, dir := range currentDir.Children {
			fmt.Println(Blue(dir.Name))
		}
		for _, file := range currentDir.Files {
			fmt.Println(Cyan(file.Name))
		}
	case "cat":
		// not implemented yet
	case "":
		// nothing
	case "echo":
		if len(args) < 2 {
			fmt.Println("echo: incorrect usage")
		} else {
			for _, arg := range args[1:] {
				fmt.Print(arg + " ")
			}
			fmt.Println("")
		}
	default:
		fmt.Println(args[0] + ": command not found")
	}
}
