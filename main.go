package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

var keyFlag = flag.String("key", "", "OpenAI API key")
var unsaveKeyFlag = flag.String("unsave", "", "OpenAI API key (won't be saved)")
var promptFlag = flag.String("prompt", "", "Prompt for OpenAI API")
var dirFlag = flag.String("dir", "", "Directory for .key file")
var helpFlag = flag.Bool("help", false, "Show help message")

func main() {
	flag.Parse()

	if *helpFlag {
		printHelp()
		os.Exit(0)
	}

	dir := *dirFlag
	if dir == "" {
		var err error
		dir, err = os.UserHomeDir()
		if err != nil {
			fmt.Println("Error getting user home directory:", err)
			os.Exit(1)
		}
	}

	if *keyFlag != "" {
		saveAPIKey(*keyFlag, dir)
	}

	var key string
	if *unsaveKeyFlag != "" {
		key = *unsaveKeyFlag
	} else {
		key = loadAPIKey(dir)
	}

	if key == "" {
		fmt.Println("Please set an API key first.")
		os.Exit(1)
	}

	var prompt string
	if *promptFlag != "" {
		prompt = *promptFlag
	} else if len(flag.Args()) > 0 {
		prompt = flag.Args()[0]
	} else {
		fmt.Println("Please provide a prompt.")
		os.Exit(1)
	}

	for {
		cmd := getCommandFromAPI(prompt, key)

		cmd = strings.TrimPrefix(cmd, "$ ")
		fmt.Printf("Generated command: %s\n", cmd)
		fmt.Print("Do you want to accept, retry, or cancel? (A/R/C): ")

		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(response)

		response = strings.ToLower(response) // Convert the response to lowercase

		switch response {
		case "a":
			executeCommand(cmd)
			os.Exit(1)
		case "r":
			fmt.Print("Try rephrasing the prompt to make the task more explicit:\n")
			prompt, _ = reader.ReadString('\n')

			// Modify the command
			// Here, we just continue the loop, which will re-run getCommandFromAPI
			continue
		case "c":
			// Cancel
			fmt.Println("Canceled.")
			os.Exit(1)
		default:
			fmt.Println("Invalid response.")
		}
	}
}

func printHelp() {
	fmt.Println(`Usage: tao [OPTIONS] [PROMPT]

Options:
  -key KEY        Set the OpenAI API key and save it for future use.
  -unsave KEY     Set the OpenAI API key without saving it.
  -prompt PROMPT  Set the prompt for the OpenAI API.
  -dir DIR        Set the directory for the .key file.
  -help           Show this help message.

Examples:
  tao -key my_open_ai_key
  tao -unsave my_open_ai_key -prompt my_prompt
  tao -dir /tmp/ -key my_open_ai_key
  tao my_prompt`)
}

func saveAPIKey(key string, dir string) {
	file := dir + "/.tao"
	os.WriteFile(file, []byte(key), 0644)
}

func loadAPIKey(dir string) string {
	file := dir + "/.tao"
	key, _ := os.ReadFile(file)
	return string(key)
}
