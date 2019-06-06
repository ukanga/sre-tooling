package main

import (
	"fmt"
	"os"

	"github.com/onaio/sre-tooling/billing"
)

const helpSubCommand string = "help"

func main() {
	cli(os.Args)
}

func cli(args []string) {
	commandName := args[0]
	if len(args) > 1 {
		switch args[1] {
		case billing.Command:
			billing.Cli(commandName, helpSubCommand, args[2:])
		case helpSubCommand:
			fmt.Println(help(commandName))
		default:
			fmt.Println(help(commandName))
			os.Exit(1)
		}
		os.Exit(0)
	} else {
		fmt.Println(help(commandName))
		os.Exit(1)
	}
}

func help(commandName string) string {
	text := `
Usage: %s [command]

Common commands:
	%s		Billing specific commands
	%s		Prints this help message
`
	return fmt.Sprintf(text, commandName, billing.Command, helpSubCommand)
}
