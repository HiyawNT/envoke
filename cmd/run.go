package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run -- [command]",
	Short: "Run a command with secrets as environment variables",
	Long: `Run a command with secrets injected as environment variables.

All secrets from the specified environment will be loaded and made available
to the command. This is useful for running applications that read configuration
from environment variables.

Examples:
  envoke run --env dev -- npm start
  envoke run --env prod -- ./my-app
  envoke run --env staging -- python app.py`,
	Args: cobra.MinimumNArgs(1),
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

		// Get all secrets
		secrets, err := svc.ListSecretsWithValues(env)
		if err != nil {
			return fmt.Errorf("failed to load secrets: %w", err)
		}

		if len(secrets) == 0 {
			fmt.Printf("  No secrets found in environment '%s'\n", env)
		} else {
			fmt.Printf(" Loaded %d secrets from '%s'\n", len(secrets), env)
		}

		// Prepare command
		cmdName := args[0]
		cmdArgs := args[1:]

		command := exec.Command(cmdName, cmdArgs...)

		// Set up environment
		command.Env = os.Environ()
		for key, value := range secrets {
			command.Env = append(command.Env, fmt.Sprintf("%s=%s", key, value))
		}

		// Connect stdio
		command.Stdin = os.Stdin
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr

		// Run command
		if err := command.Run(); err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				os.Exit(exitErr.ExitCode())
			}
			return fmt.Errorf("failed to run command: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringVarP(&envFlag, "env", "e", "", "Environment name")

	// Custom usage to show the -- separator
	runCmd.SetUsageTemplate(strings.ReplaceAll(runCmd.UsageTemplate(),
		"{{.UseLine}}",
		"{{.CommandPath}} [flags] -- [command]"))
}
