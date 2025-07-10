# rrk

[![GitHub release](https://img.shields.io/github/release/MRyutaro/rrk.svg)](https://github.com/MRyutaro/rrk/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**rrk** (rireki) is a Go-based single-binary CLI tool that manages bash/zsh shell history by **session** and **directory**, enabling logical grouping and easy re-execution of past commands.

> 📖 **日本語ドキュメント** - [README.ja.md](./README.ja.md)

## Features

- 📁 **Directory-based** history management
- 🪟 **Session-based** history management  
- 🔄 **One-command re-execution** - `rrk rerun <ID>`
- 🚀 **Single binary** - no dependencies
- 💾 **Lightweight** - file-based storage, no database required
- 🐚 **Shell integration** - supports bash/zsh

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

### List All History

```bash
# Show all command history
rrk list

# Show latest 20 entries only
rrk list -n 20
```

### Session Management

```bash
# List all sessions
rrk session list

# Show current session history
rrk session show
rrk s show

# Show specific session history
rrk session show <SESSION_ID>
```

### Directory Management

```bash
# Show current directory history
rrk dir show
rrk d show

# Show specific directory history
rrk dir show /path/to/directory

# List directories with history (with numeric IDs)
rrk dir list

# Show directory history by ID
rrk dir show <ID>
rrk d show <ID>
```

### Command Re-execution

```bash
# Re-execute command by history ID
rrk rerun <HISTORY_ID>

# Example: re-execute command with ID=1
rrk rerun 1
```

### Update rrk

```bash
# Update to the latest version
rrk update
```

### Version Information

```bash
# Show version (includes latest GitHub release info)
rrk -v
rrk --version
```

> **Note**: Version comparison correctly handles development builds and update notifications are automatically cleared after running `rrk update`.

### Uninstall

```bash
# Remove shell integration only
rrk uninstall

# Remove data as well
rrk uninstall --remove-data

# Auto-confirm uninstall
rrk uninstall -y --remove-data
```

## Example Usage

```bash
# Check current directory history
$ rrk dir show
ID  TIME      SESSION        COMMAND
1   14:30:12  abc123...      git status
2   14:30:45  abc123...      git add .
3   14:31:02  abc123...      git commit -m "fix bug"

# Re-execute a specific command
$ rrk rerun 2
git add .

# Check session history
$ rrk session show
ID  TIME      DIRECTORY       COMMAND
1   14:30:12  ~/project       git status
2   14:30:45  ~/project       git add .
3   14:31:02  ~/project       git commit -m "fix bug"
4   14:32:15  ~/documents     vim README.md

# Check all history
$ rrk list
ID  TIME      DIRECTORY       SESSION        COMMAND
1   14:30:12  ~/project       abc123...      git status
2   14:30:45  ~/project       abc123...      git add .
3   14:31:02  ~/project       abc123...      git commit -m "fix bug"
4   14:32:15  ~/documents     def456...      vim README.md

# List directories with IDs
$ rrk dir list
ID  DIRECTORY        STATUS
0   ~/project        (current)
1   ~/documents
2   /tmp

# Show history by directory ID
$ rrk dir show 1
ID  TIME      SESSION        COMMAND
4   14:32:15  def456...      vim README.md
```

## Data Storage

- History data is stored in `~/rrk/history.jsonl`
- Session information is stored in `~/rrk/current_session`
- No external database required

## For Developers

See [`docs/DEVELOPERS.md`](./docs/DEVELOPERS.md) for detailed information.

### Contributing

- Pull requests merged to the `main` branch automatically trigger a patch version release
- The CI/CD pipeline handles version bumping and GitHub releases automatically

## License

MIT License
