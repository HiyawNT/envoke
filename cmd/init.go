package cmd

import (
	"fmt"
	"syscall"

	"github.com/HiyawNT/envoke/internal/config"
	"github.com/HiyawNT/envoke/internal/crypto"
	"github.com/HiyawNT/envoke/internal/storage"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize envoke with a master passphrase",
	Long: `Initialize envoke by creating a configuration file and setting up encryption.

This command will:
- Prompt for a master passphrase
- Generate a random salt for key derivation
- Create the configuration directory and database
- Mark the installation as initialized

The master passphrase is used to derive an encryption key using Argon2id.
This key encrypts all secrets before storage. The passphrase is never stored.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.New()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		if cfg.IsInitialized() {
			return fmt.Errorf("envoke is already initialized")
		}

		fmt.Println("🔐 Initializing envoke...")
		fmt.Println()

		// Get passphrase
		fmt.Print("Enter master passphrase: ")
		passphrase, err := term.ReadPassword(int(syscall.Stdin))
		fmt.Println()
		if err != nil {
			return fmt.Errorf("failed to read passphrase: %w", err)
		}

		if len(passphrase) < 8 {
			return fmt.Errorf("passphrase must be at least 8 characters")
		}

		// Confirm passphrase
		fmt.Print("Confirm passphrase: ")
		confirmPassphrase, err := term.ReadPassword(int(syscall.Stdin))
		fmt.Println()
		if err != nil {
			return fmt.Errorf("failed to read passphrase: %w", err)
		}

		if string(passphrase) != string(confirmPassphrase) {
			return fmt.Errorf("passphrases do not match")
		}

		// Generate salt
		salt, err := crypto.GenerateSalt()
		if err != nil {
			return err
		}

		// Store salt in config
		if err := cfg.SetSalt(crypto.EncodeSalt(salt)); err != nil {
			return fmt.Errorf("failed to save salt: %w", err)
		}

		// Initialize database
		dbPath, err := config.GetDBPath()
		if err != nil {
			return fmt.Errorf("failed to get database path: %w", err)
		}

		store, err := storage.NewSQLiteStorage(dbPath)
		if err != nil {
			return fmt.Errorf("failed to initialize database: %w", err)
		}
		defer store.Close()

		// Mark as initialized
		if err := cfg.SetInitialized(true); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		fmt.Println()
		fmt.Println("✅ envoke initialized successfully!")
		fmt.Println()
		fmt.Println("Next steps:")
		fmt.Println("  1. Create an environment: envoke env create dev")
		fmt.Println("  2. Add a secret: envoke secret set API_KEY")
		fmt.Println("  3. View secrets: envoke tui")
		fmt.Println()
		fmt.Println("⚠️  IMPORTANT: Keep your master passphrase safe!")
		fmt.Println("   Without it, you cannot decrypt your secrets.")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
