package cmd

import (
	"fmt"
	"os"
	"syscall"

	"github.com/HiyawNT/envoke/internal/config"
	"github.com/HiyawNT/envoke/internal/crypto"
	"github.com/HiyawNT/envoke/internal/service"
	"github.com/HiyawNT/envoke/internal/storage"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	cfg     *config.Config
	svc     *service.SecretService
	cfgFile string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "envoke",
	Short: "Encrypted secret manager for the terminal",
	Long: `envoke (env + invoke) - A local-first encrypted secret manager
	
Store and manage secrets across multiple environments with military-grade encryption.
All secrets are encrypted using NaCl secretbox (XSalsa20-Poly1305) before storage.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Skip initialization check for init and help commands
		if cmd.Name() == "init" || cmd.Name() == "help" || cmd.Parent() == nil {
			return nil
		}

		var err error
		cfg, err = config.New()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		if !cfg.IsInitialized() {
			return fmt.Errorf("envoke is not initialized. Run 'envoke init' first")
		}

		// Initialize service for all commands except init
		if err := initService(); err != nil {
			return err
		}

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/envoke/config.yaml)")
}

// initService initializes the secret service with encryption
func initService() error {
	// Get passphrase from user
	fmt.Print("Enter master passphrase: ")
	passphrase, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		return fmt.Errorf("failed to read passphrase: %w", err)
	}

	if len(passphrase) == 0 {
		return fmt.Errorf("passphrase cannot be empty")
	}

	// Get salt from config
	saltEncoded := cfg.GetSalt()
	if saltEncoded == "" {
		return fmt.Errorf("salt not found in config. Run 'envoke init' first")
	}

	salt, err := crypto.DecodeSalt(saltEncoded)
	if err != nil {
		return fmt.Errorf("failed to decode salt: %w", err)
	}

	// Get verification Hash from config
	verificationHash := cfg.GetVerificationHash()
	if verificationHash == "" {

		return fmt.Errorf("verification hash not found in config. Run 'envoke init' first")
	}

	// verify passphrase
	if !crypto.VerifyPassphrase(string(passphrase), salt, verificationHash) {

		return fmt.Errorf("Incorrect Passphrase")
	}

	// Derive key from passphrase (only after verification passed)
	key := crypto.DeriveKey(string(passphrase), salt)

	// Create encryptor
	encryptor, err := crypto.NewEncryptor(key)
	if err != nil {
		return fmt.Errorf("failed to create encryptor: %w", err)
	}

	// Initialize storage
	dbPath, err := config.GetDBPath()
	if err != nil {
		return fmt.Errorf("failed to get database path: %w", err)
	}

	store, err := storage.NewSQLiteStorage(dbPath)
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}

	// Create service
	svc = service.NewSecretService(store, encryptor)

	return nil
}
