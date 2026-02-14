package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

// envCmd represents the env command
var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Manage environments",
	Long:  `Create, list, and manage environments for organizing secrets.`,
}

// envCreateCmd creates a new environment
var envCreateCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new environment",
	Long: `Create a new environment for organizing secrets.

Examples:
  envoke env create dev
  envoke env create staging
  envoke env create prod`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		envName := args[0]

		if err := svc.CreateEnvironment(envName); err != nil {
			return fmt.Errorf("failed to create environment: %w", err)
		}

		fmt.Printf("✅ Environment '%s' created successfully\n", envName)
		return nil
	},
}

// envListCmd lists all environments
var envListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all environments",
	Aliases: []string{"ls"},
	RunE: func(cmd *cobra.Command, args []string) error {
		envs, err := svc.ListEnvironments()
		if err != nil {
			return fmt.Errorf("failed to list environments: %w", err)
		}

		if len(envs) == 0 {
			fmt.Println("No environments found. Create one with: envoke env create <name>")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "NAME\tACTIVE\tCREATED")

		for _, env := range envs {
			active := ""
			if env.IsActive {
				active = "✓"
			}
			fmt.Fprintf(w, "%s\t%s\t%s\n", env.Name, active, env.CreatedAt.Format("2006-01-02 15:04"))
		}

		w.Flush()
		return nil
	},
}

// envUseCmd sets the active environment
var envUseCmd = &cobra.Command{
	Use:   "use [name]",
	Short: "Set the active environment",
	Long:  `Set an environment as the active one for commands that don't specify an environment.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		envName := args[0]

		if err := svc.SetActiveEnvironment(envName); err != nil {
			return fmt.Errorf("failed to set active environment: %w", err)
		}

		fmt.Printf("✅ Active environment set to '%s'\n", envName)
		return nil
	},
}

// envDeleteCmd deletes an environment
var envDeleteCmd = &cobra.Command{
	Use:     "delete [name]",
	Short:   "Delete an environment",
	Aliases: []string{"rm"},
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		envName := args[0]

		// Confirm deletion
		fmt.Printf("⚠️  This will delete environment '%s' and all its secrets. Continue? (y/N): ", envName)
		var response string
		fmt.Scanln(&response)

		if response != "y" && response != "Y" {
			fmt.Println("Cancelled")
			return nil
		}

		if err := svc.DeleteEnvironment(envName); err != nil {
			return fmt.Errorf("failed to delete environment: %w", err)
		}

		fmt.Printf("✅ Environment '%s' deleted\n", envName)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(envCmd)
	envCmd.AddCommand(envCreateCmd)
	envCmd.AddCommand(envListCmd)
	envCmd.AddCommand(envUseCmd)
	envCmd.AddCommand(envDeleteCmd)
}
