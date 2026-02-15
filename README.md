# 🔐 envoke

**envoke** (env + invoke) - A production-grade, local-first encrypted secret manager for the terminal.

Store and manage secrets across multiple environments with military-grade encryption. All secrets are encrypted using NaCl secretbox (XSalsa20-Poly1305) before storage.

## ✨ Features

- 🔒 **Military-grade encryption**: NaCl secretbox (XSalsa20-Poly1305)
- 🏠 **Local-first**: All secrets stored encrypted on your machine
- 🎨 **Beautiful TUI**: Full-screen terminal interface with Bubble Tea
- 🌍 **Multi-environment**: Separate secrets for dev, staging, prod, etc.
- 🚀 **Zero trust**: Master passphrase never stored, keys derived with Argon2id
- 🔧 **CLI-first**: Scriptable commands for automation
- 📦 **Portable**: Single binary, no dependencies

## 🚀 Quick Start

### Installation

```bash
# Clone the repository
git clone https://github.com/HiyawNT/envoke.git
cd envoke

# Build
go build -o envoke .

# Install globally (optional)
go install
```

### Initialize

```bash
# Initialize envoke with a master passphrase
envoke init
```

**Important**: Your master passphrase is never stored. Keep it safe - without it, you cannot decrypt your secrets.

### Create an Environment

```bash
# Create environments for different stages
envoke env create dev
envoke env create staging
envoke env create prod

# Set the active environment
envoke env use dev
```

### Manage Secrets

```bash
# Add a secret (prompts for value)
envoke secret set API_KEY --env dev
envoke secret set DATABASE_URL --env dev

# Get a secret
envoke secret get API_KEY --env dev

# List all secrets (keys only, values hidden)
envoke secret list --env dev

# Delete a secret
envoke secret delete API_KEY --env dev
```

### Run Commands with Secrets

```bash
# Inject secrets as environment variables
envoke run --env dev -- npm start
envoke run --env prod -- ./my-app
```

### Use the TUI

```bash
# Launch the interactive terminal UI
envoke tui
```

#### TUI Keyboard Shortcuts

- `↑/↓` or `j/k` - Navigate
- `←/→` or `h/l` - Switch environments
- `enter` - Select/confirm
- `a` - Add new secret
- `t` - Toggle secret visibility
- `e` - Edit secret (coming soon)
- `d` - Delete secret (coming soon)
- `esc` - Go back
- `q` - Quit

## 🔒 Security Architecture

### Encryption

- **Algorithm**: NaCl secretbox (XSalsa20-Poly1305)
- **Key Derivation**: Argon2id (memory-hard, side-channel resistant)
- **Key Size**: 256 bits (32 bytes)
- **Nonce**: 192 bits (24 bytes), randomly generated per secret
- **Authentication**: Poly1305 MAC (prevents tampering)

### Key Derivation Parameters

Argon2id is configured for ~100ms on modern hardware:
- **Time**: 2 iterations
- **Memory**: 64 MB
- **Parallelism**: 4 threads
- **Output**: 32 bytes (256 bits)

### Storage

- **Database**: SQLite with encrypted blobs
- **Location**: `~/.config/envoke/envoke.db`
- **Schema**: Environments, Secrets (encrypted), Metadata
- **Permissions**: Config directory created with `0700` (owner-only)

### What is Stored

-  Encrypted secret values (NaCl secretbox output)
-  Nonces (24 bytes per secret)
-  Salt for key derivation (16 bytes, random)
-  Master passphrase (never stored)
-  Encryption keys (derived on-demand)
-  Plaintext secrets (encrypted before storage)

## 📂 Project Structure

```
envoke/
├── cmd/                    # CLI commands
│   ├── root.go            # Root command + auth
│   ├── init.go            # Initialize envoke
│   ├── env.go             # Environment management
│   ├── secret.go          # Secret operations
│   ├── run.go             # Run commands with secrets
│   ├── tui.go             # Launch TUI
│   └── login.go           # Cloud sync (stub)
├── internal/
│   ├── crypto/            # Encryption/decryption
│   ├── storage/           # SQLite storage
│   ├── models/            # Database models
│   ├── tui/               # Terminal UI
│   ├── config/            # Configuration management
│   └── service/           # Business logic
├── main.go
├── go.mod
└── README.md
```

## 🛠️ Development

### Prerequisites

- Go 1.21+
- SQLite3

### Build

```bash
go build -o envoke .
```

### Run Tests

```bash
go test ./...
```

### Development Mode

```bash
# Run without installing
go run . init
go run . env create dev
go run . tui
```

## 🗺️ Roadmap

### Phase 1 (Current)

- ✅ Core encryption (NaCl secretbox)
- ✅ Local storage (SQLite)
- ✅ CLI commands
- ✅ Terminal UI (Bubble Tea)
- ✅ Multi-environment support

### Phase 2 (Next)

- [ ] `.env` import/export
- [ ] Secret rotation
- [ ] Audit logging
- [ ] Secret templates
- [ ] Bulk operations

### Phase 3 (Future)

- [ ] S3 sync (encrypted backup)
- [ ] Team sharing (with key exchange)
- [ ] Secret versioning
- [ ] RBAC (role-based access control)
- [ ] Plugin system

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

### Security Issues

If you discover a security vulnerability, please email security@example.com instead of using the issue tracker.

## 📄 License

MIT License - see LICENSE file for details.

## 🙏 Acknowledgments

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - Terminal UI framework
- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [NaCl](https://nacl.cr.yp.to/) - Cryptography library
- [GORM](https://gorm.io/) - ORM for Go


## ⚠️ Disclaimer

This is a production-quality starter template. Before using in production:

1. **Audit the cryptography**: Have a security expert review the crypto implementation
2. **Test thoroughly**: Add comprehensive test coverage
3. **Backup your secrets**: Keep secure backups of your master passphrase and config
4. **Monitor for vulnerabilities**: Keep dependencies updated

---

Made with ❤️ using Go and Bubble Tea
