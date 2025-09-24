package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/mattn/go-shellwords"
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Interactive mode. Type 'help' or 'exit'.")
			return startREPL()
		},
	}

	//root.CompletionOptions.HiddenDefaultCmd = true
	root.AddCommand(newGitCLICmd())

	return root
}

func Execute() {
	if err := NewRootCmd().Execute(); err != nil {
		_, err := fmt.Fprintln(os.Stderr, err)
		if err != nil {
			return
		}
		os.Exit(1)
	}
}

func startREPL() error {
	reader := bufio.NewScanner(os.Stdin)
	parser := shellwords.NewParser() // понимает кавычки, экранирование и т.п.

	for {
		fmt.Print("> ")
		if !reader.Scan() {
			return reader.Err()
		}
		line := strings.TrimSpace(reader.Text())
		if line == "" {
			continue
		}
		if line == "exit" || line == "quit" {
			fmt.Println("bye!")
			return nil
		}

		args, err := parser.Parse(line)
		if err != nil {
			_, err := fmt.Fprintln(os.Stderr, "parse error:", err)
			if err != nil {
				return err
			}
			continue
		}

		rc := NewRootCmd()
		rc.SetArgs(args)
		if err := rc.Execute(); err != nil {
			_, err := fmt.Fprintln(os.Stderr, err)
			if err != nil {
				return err
			}
		}
	}
}
