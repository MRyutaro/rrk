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

### Build from Source

```bash
git clone https://github.com/MRyutaro/rrk.git
cd rrk
make build
sudo mv rrk /usr/local/bin/
```

## Setup

After installation, enable shell integration:

```bash
# Automatic setup (recommended)
rrk setup

# Auto-confirm setup
rrk setup -y
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

# List directories with history
rrk dir list
```

### Command Re-execution

```bash
# Re-execute command by history ID
rrk rerun <HISTORY_ID>

# Example: re-execute command with ID=1
rrk rerun 1
```

### Version Information

```bash
# Show version (includes latest GitHub release info)
rrk -v
rrk --version
```

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
Re-running: git add .
Original directory: /Users/user/project
Current directory: /Users/user/project

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
```

## Data Storage

- History data is stored in `~/rrk/history.jsonl`
- Session information is stored in `~/rrk/current_session`
- No external database required

## For Developers

See [`docs/DEVELOPERS.md`](./docs/DEVELOPERS.md) for detailed information.

## License

MIT License
