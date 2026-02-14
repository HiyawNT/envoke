package cmd

import (
	"fmt"

	"github.com/HiyawNT/envoke/internal/tui"
	"github.com/spf13/cobra"
)

// tuiCmd represents the tui command
var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Launch the terminal UI",
	Long: `Launch the interactive terminal user interface for managing secrets.

The TUI provides a full-screen interface for:
- Browsing environments
- Viewing secrets (with toggle for visibility)
- Adding and editing secrets
- Navigation with keyboard shortcuts

Keyboard shortcuts:
  ↑/↓ or j/k    Navigate
  ←/→ or h/l    Switch environments
  enter         Select/confirm
  a             Add new secret
  t             Toggle secret visibility
  esc           Go back
  q             Quit`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := tui.Run(svc); err != nil {
			return fmt.Errorf("TUI error: %w", err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}
