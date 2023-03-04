package main

import (
	"bufio"
	"fmt"
	"github.com/snowflake-server-go/src/server"
	"log"
	"os"
	"strings"
)

func main() {
	s, err := server.NewServer(":10000")
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		if err := s.Start(); err != nil {
			log.Fatal(err)
		}
	}()
	fmt.Println("Server started on port 10000. Type 'help' for a list of commands. Press enter to exit.")

	validCommands := map[string]string{
		"help": "Display this help message.",
		"exit": "Exit the program.",
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("$ ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		if text == "" {
			continue
		}
		if err := handleCommand(text, validCommands, reader); err != nil {
			fmt.Println(err)
		}
		if text == "exit" {
			if confirmExit(reader) {
				break
			}
		}
	}
}

func handleCommand(cmd string, validCommands map[string]string, reader *bufio.Reader) error {
	switch cmd {
	case "help":
		fmt.Println("Valid commands:")
		for cmd, usage := range validCommands {
			fmt.Printf("%-10s%s\n", cmd, usage)
		}
	case "exit":
		return nil
	default:
		if usage, ok := validCommands[cmd]; ok {
			fmt.Printf("%s (%s)\n", cmd, usage)
		} else {
			for name, _ := range validCommands {
				if strings.HasPrefix(name, cmd) {
					fmt.Printf("Did you mean '%s'? (%s)\n", name, validCommands[name])
					return nil
				}
			}
			return fmt.Errorf("unknown command: %s", cmd)
		}
	}
	return nil
}

func confirmExit(reader *bufio.Reader) bool {
	fmt.Print("Are you sure you want to exit? (yes/no) ")
	answer, _ := reader.ReadString('\n')
	answer = strings.TrimSpace(strings.ToLower(answer))
	return answer == "yes" || answer == "y" || answer == "ok"
}
