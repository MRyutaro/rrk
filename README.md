# rrk

[![GitHub release](https://img.shields.io/github/release/MRyutaro/rrk.svg)](https://github.com/MRyutaro/rrk/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**rrk** (rireki) is a Go-based single-binary CLI tool that displays shell history in directory tree format, making it easy to see which commands were executed in each directory.

> 📖 **日本語ドキュメント** - [README.ja.md](./README.ja.md)

## Features

- 🌳 **Tree visualization** - Display command history in directory tree format
- 📁 **Directory-based** organization - See which commands were run in each directory
- 🎯 **Focused view** - View history for specific directories
- 🚀 **Single binary** - no dependencies
- 💾 **Lightweight** - file-based storage, no database required
- 🐚 **Shell integration** - supports bash/zsh with automatic setup
- 🔄 **Auto-update** - built-in update mechanism with GitHub releases
- 🗑️ **Easy removal** - clean uninstall with data preservation options

## Installation

### Quick Install (Recommended)

```bash
curl -LsSf https://raw.githubusercontent.com/MRyutaro/rrk/main/install.sh | sh
```

This script will:
- Download the appropriate binary for your system
- Install it to `~/.local/bin` (or `$INSTALL_DIR` if set)
- Automatically set up shell integration (bash/zsh)
- Add the installation directory to your PATH if needed

### Build from Source

```bash
git clone https://github.com/MRyutaro/rrk.git
cd rrk
make build
sudo mv rrk /usr/local/bin/

# After building from source, set up shell integration:
rrk setup
```

## Usage

### Tree Visualization

```bash
# Display all command history in tree format
rrk

# Display history for a specific directory
rrk /path/to/directory

# Limit the number of commands shown per directory
rrk -n 5
```

### Update rrk

```bash
# Update to the latest version
rrk update
```

The update command will:
- Download the latest version from GitHub releases
- Replace the current binary
- Verify the installation
- Clear update notification cache

### Version Information

```bash
# Show version (includes update notification if available)
rrk version
rrk -v
rrk --version
```

### Uninstall

```bash
# Remove shell integration and all data
rrk uninstall

# Remove shell integration and all data without confirmation
rrk uninstall -y
```

The uninstall command will:
- Remove shell integration from `~/.bashrc`/`~/.zshrc`
- Delete all rrk data from `~/.rrk/`
- Provide instructions for removing the binary

## Example Usage

```bash
# Display command history in tree format
$ rrk
/home/user
├── project/
│   ├── git status
│   ├── git add .
│   └── git commit -m "fix bug"
├── scripts/
│   ├── ./deploy.sh
│   └── python backup.py
└── .config/
    └── vim init.vim

/var
└── log/
    ├── tail -f syslog
    ├── grep ERROR *.log
    └── journalctl -u nginx

# View history for a specific directory
$ rrk /home/user/project
├── git status
├── git add .
└── git commit -m "fix bug"

# Limit commands per directory
$ rrk -n 2
/home/user
├── project/
│   ├── git add .
│   └── git commit -m "fix bug"
└── scripts/
    └── python backup.py
```

## Data Storage

- History data is stored in `~/.rrk/history.jsonl` (JSONL format)
- Session information is stored in `~/.rrk/current_session`
- Shell integration script is stored in `~/.rrk/hook.sh`
- Version cache is stored in `~/.rrk/.rrk_version_cache`
- No external database required

## Advanced Usage

### Manual Shell Integration

If you prefer manual setup or need custom configuration:

```bash
# Generate shell integration script
rrk hook init bash > ~/.rrk_integration.sh
rrk hook init zsh > ~/.rrk_integration.sh

# Source it in your shell configuration
echo "source ~/.rrk_integration.sh" >> ~/.bashrc  # or ~/.zshrc
```

### Manual History Recording

```bash
# Record a command manually
rrk hook record "your command here"

# Initialize a new session
rrk hook session-init
```

## CI/CD Integration

rrk includes automated release management:

- **Pull Request Merging**: Automatically creates patch releases
- **Manual Tagging**: Triggers release builds for all platforms
- **Multi-platform Builds**: Linux, macOS, Windows (AMD64/ARM64)
- **Automatic Updates**: Built-in update notification and installation

## For Developers

See [`docs/DEVELOPERS.md`](./docs/DEVELOPERS.md) and [`docs/REQ.md`](./docs/REQ.md) for detailed information.

### Contributing

- Pull requests merged to the `main` branch automatically trigger a patch version release
- The CI/CD pipeline handles version bumping and GitHub releases automatically
- Use `make patch`, `make minor`, or `make major` for local version management

## License

MIT License