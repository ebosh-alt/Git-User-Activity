package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github_user_activity/internal/gh"
)

func newActivityCmd() *cobra.Command {
	var tail int
	var outType string

	cmd := &cobra.Command{
		Use:   "activity <username>",
		Short: "Показать последние публичные действия пользователя GitHub",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			username := strings.TrimSpace(args[0])
			if username == "" {
				return errors.New("укажите имя пользователя GitHub")
			}
			if tail <= 0 {
				tail = 10
			}
			switch strings.ToLower(outType) {
			case "", "txt":
				outType = "txt"
			case "json":
				outType = "json"
			default:
				return fmt.Errorf("неизвестный тип вывода: %q (доступно: txt, json)", outType)
			}
			client := gh.Client{}
			events, err := client.UserEvents(cmd.Context(), username, tail)
			if err != nil {
				return err
			}
			return renderEvents(os.Stdout, events, outType)
		},
	}

	cmd.Flags().IntVarP(&tail, "tail", "t", 10, "сколько событий показать (с пагинацией)")
	cmd.Flags().StringVarP(&outType, "type", "T", "", "тип вывода (json, text)")
	return cmd
}

func renderEvents(dst *os.File, events []gh.Event, typ string) error {
	switch typ {
	case "txt":
		for _, e := range events {
			_, err := fmt.Fprintln(dst, gh.Human(e))
			if err != nil {
				return err
			}
		}
		return nil
	case "json":
		enc := json.NewEncoder(dst)
		enc.SetIndent("", "  ")
		return enc.Encode(events)
	default:
		return fmt.Errorf("неизвестный тип вывода: %q", typ)
	}
}
