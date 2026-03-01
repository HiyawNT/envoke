package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to cloud sync (coming soon)",
	Long: `Login to enable cloud sync with S3-compatible storage.

This feature is under development and will allow you to:
- Sync secrets across multiple machines
- Backup secrets to S3/MinIO
- Collaborate with team members (future)

Currently, envoke operates in local-only mode.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(" Cloud sync is not yet implemented")
		fmt.Println()
		fmt.Println("envoke currently operates in local-only mode.")
		fmt.Println("All secrets are stored encrypted in: ~/.config/envoke/envoke.db")
		fmt.Println()
		fmt.Println("Cloud sync with S3 will be available in a future release.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
