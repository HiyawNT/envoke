package cmd

import (
	"fmt"
	"os"
	"syscall"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	envFlag string
)

// secretCmd represents the secret command
var secretCmd = &cobra.Command{
	Use:   "secret",
	Short: "Manage secrets",
	Long:  `Set, get, list, and delete encrypted secrets within environments.`,
}

// secretSetCmd sets a secret value
var secretSetCmd = &cobra.Command{
	Use:   "set [key]",
	Short: "Set a secret value",
	Long: `Set or update a secret value in the specified environment.

The value will be prompted securely (hidden input). The secret is encrypted
before storage using NaCl secretbox (XSalsa20-Poly1305).

Examples:
  envoke secret set API_KEY --env dev
  envoke secret set DATABASE_URL --env prod`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]

		// Determine environment
		env := envFlag
		if env == "" {
			activeEnv, err := svc.GetActiveEnvironment()
			if err != nil {
				return fmt.Errorf("no active environment set. Use --env flag or run 'envoke env use <n>'")
			}
			env = activeEnv.Name
		}

		// Prompt for value
		fmt.Printf("Enter value for %s: ", key)
		value, err := term.ReadPassword(int(syscall.Stdin))
		fmt.Println()
		if err != nil {
			return fmt.Errorf("failed to read value: %w", err)
		}

		if len(value) == 0 {
			return fmt.Errorf("value cannot be empty")
		}

		// Confirm value
		fmt.Print("Confirm value: ")
		confirmValue, err := term.ReadPassword(int(syscall.Stdin))
		fmt.Println()
		if err != nil {
			return fmt.Errorf("failed to read value: %w", err)
		}

		if string(value) != string(confirmValue) {
			return fmt.Errorf("values do not match")
		}

		// Set secret
		if err := svc.SetSecret(env, key, string(value)); err != nil {
			return fmt.Errorf("failed to set secret: %w", err)
		}

		fmt.Printf(" Secret '%s' set in environment '%s'\n", key, env)
		return nil
	},
}

// secretGetCmd retrieves a secret value
var secretGetCmd = &cobra.Command{
	Use:   "get [key]",
	Short: "Get a secret value",
	Long: `Retrieve and decrypt a secret value from the specified environment.

Examples:
  envoke secret get API_KEY --env dev
  envoke secret get DATABASE_URL --env prod`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]

		// Determine environment
		env := envFlag
		if env == "" {
			activeEnv, err := svc.GetActiveEnvironment()
			if err != nil {
				return fmt.Errorf("no active environment set. Use --env flag or run 'envoke env use <n>'")
			}
			env = activeEnv.Name
		}

		// Get secret
		value, err := svc.GetSecret(env, key)
		if err != nil {
			return fmt.Errorf("failed to get secret: %w", err)
		}

		fmt.Println(value)
		return nil
	},
}

// secretListCmd lists all secrets in an environment
var secretListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all secrets",
	Aliases: []string{"ls"},
	Long: `List all secret keys in the specified environment (values are not shown).

Examples:
  envoke secret list --env dev
  envoke secret list --env prod`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Determine environment
		env := envFlag
		if env == "" {
			activeEnv, err := svc.GetActiveEnvironment()
			if err != nil {
				return fmt.Errorf("no active environment set. Use --env flag or run 'envoke env use <n>'")
			}
			env = activeEnv.Name
		}

		// List secrets
		secrets, err := svc.ListSecrets(env)
		if err != nil {
			return fmt.Errorf("failed to list secrets: %w", err)
		}

		if len(secrets) == 0 {
			fmt.Printf("No secrets in environment '%s'\n", env)
			return nil
		}

		fmt.Printf("Secrets in environment '%s':\n\n", env)

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "KEY\tCREATED\tUPDATED")

		for _, secret := range secrets {
			fmt.Fprintf(w, "%s\t%s\t%s\n", secret.Key, secret.CreatedAt, secret.UpdatedAt)
		}

		w.Flush()
		return nil
	},
}

// secretDeleteCmd deletes a secret
var secretDeleteCmd = &cobra.Command{
	Use:     "delete [key]",
	Short:   "Delete a secret",
	Aliases: []string{"rm"},
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]

		// Determine environment
		env := envFlag
		if env == "" {
			activeEnv, err := svc.GetActiveEnvironment()
			if err != nil {
				return fmt.Errorf("no active environment set. Use --env flag or run 'envoke env use <n>'")
			}
			env = activeEnv.Name
		}

		// Confirm deletion
		fmt.Printf("Delete secret '%s' from environment '%s'? (y/N): ", key, env)
		var response string
		fmt.Scanln(&response)

		if response != "y" && response != "Y" {
			fmt.Println("Cancelled")
			return nil
		}

		// Delete secret
		if err := svc.DeleteSecret(env, key); err != nil {
			return fmt.Errorf("failed to delete secret: %w", err)
		}

		fmt.Printf("✅ Secret '%s' deleted from environment '%s'\n", key, env)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(secretCmd)
	secretCmd.AddCommand(secretSetCmd)
	secretCmd.AddCommand(secretGetCmd)
	secretCmd.AddCommand(secretListCmd)
	secretCmd.AddCommand(secretDeleteCmd)

	// Add --env flag to all secret commands
	secretSetCmd.Flags().StringVarP(&envFlag, "env", "e", "", "Environment name")
	secretGetCmd.Flags().StringVarP(&envFlag, "env", "e", "", "Environment name")
	secretListCmd.Flags().StringVarP(&envFlag, "env", "e", "", "Environment name")
	secretDeleteCmd.Flags().StringVarP(&envFlag, "env", "e", "", "Environment name")
}
