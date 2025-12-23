# Orbit

```
  ██████╗  ██████╗  ██████╗  ██╗ ████████╗
 ██╔═══██╗ ██╔══██╗ ██╔══██╗ ██║ ╚══██╔══╝
 ██║   ██║ ██████╔╝ ██████╔╝ ██║    ██║   
 ██║   ██║ ██╔══██╗ ██╔══██╗ ██║    ██║   
 ╚██████╔╝ ██║  ██║ ██████╔╝ ██║    ██║   
  ╚═════╝  ╚═╝  ╚═╝ ╚═════╝  ╚═╝    ╚═╝   
```

> Keep your side projects in orbit

Orbit is a powerful CLI tool and TUI application for managing your side projects. Built with Go, using [Cobra](https://github.com/spf13/cobra) for CLI, [Bubble Tea](https://github.com/charmbracelet/bubbletea) for TUI, and [Lip Gloss](https://github.com/charmbracelet/lipgloss) for styling.

![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/license-MIT-green)

## Features

- **Beautiful TUI** - Interactive terminal UI with ASCII banner
- **Workspace Management** - Initialize and manage multiple workspaces
- **Project Organization** - Create projects with repo, docs, and secret folders
- **Status Tracking** - Track project status (active, archived, done)
- **Aliases** - Set short aliases for projects with long names
- **README Viewer** - Beautiful markdown rendering in terminal

## Installation

### Quick Setup (Recommended)

Run the automated setup script to install all dependencies:

```bash
git clone https://github.com/henrynguci/orbit.git
cd orbit
./setup.sh
```

This script will automatically install:
- Go 1.21+
- glow (markdown viewer)
- gp (Charm's Git tool)
- glamour (markdown renderer)
- bat (better cat)
- fzf (fuzzy finder)
- golangci-lint
- All Go project dependencies

### Manual Installation

If you prefer to install manually:

```bash
git clone https://github.com/henrynguci/orbit.git
cd orbit
go build -o orbit .
```

### Move to PATH

```bash
sudo mv orbit /usr/local/bin/
```

## Usage

### TUI Mode

Simply run `orbit` to open the interactive TUI with the ORBIT banner:

```bash
orbit
```

The TUI provides a menu with:
- **Init Workspace** - Create workspace with optional project
- **List Projects** - View all projects
- **Alias** - Set project aliases
- **Info** - View README.md
- **Set Status** - Change project status
- **Status** - View project status

### CLI Commands

#### Initialize a Workspace

```bash
orbit init ~/my-workspace

orbit init ~/my-workspace --project my-project
orbit init ~/my-workspace -p my-project
```

With project, creates:
```
~/my-workspace/
└── project/
    └── my-project/
        ├── repo/
        ├── docs/
        └── secret/
```

#### List Projects

```bash
orbit ls
```

#### View Project Info

```bash
orbit info <project-name>
```

#### Set Project Status

```bash
orbit set <project-name> <status>

orbit set myproject active
orbit set oldproject archived
orbit set completed-project done
```

#### Get Project Status

```bash
orbit status <project-name>
```

#### Set Alias

```bash
orbit alias <project-name> <alias>
```

## TUI Features

- **ASCII Banner** - Beautiful ORBIT logo on startup
- **Menu Navigation** - Use arrow keys or j/k
- **Multi-step Forms** - Guided workspace/project creation
- **Project List** - Filterable with status icons
- **README Viewer** - Glamour markdown rendering with scroll

## Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `↑/↓` or `j/k` | Navigate |
| `Enter` | Select/Confirm |
| `/` | Filter (in lists) |
| `Esc` | Go back |
| `q` | Quit |

## Configuration

Configuration is stored in `~/.config/orbit/config.json`:

```json
{
  "workspaces": [
    "/home/user/workspace1"
  ],
  "projects": {
    "myproject": {
      "name": "myproject",
      "alias": "mp",
      "path": "/home/user/workspace1/project/myproject",
      "status": "active"
    }
  }
}
```

## Rich Toolset

The `setup.sh` script installs a comprehensive set of CLI tools that enhance the Orbit experience:

| Category | Tool | Description |
|----------|------|-------------|
| **Core** | [Go](https://go.dev) | Programming language (v1.21+) |
| **Build** | [Make](https://www.gnu.org/software/make/) | Build automation tool |
| **Git** | [Git](https://git-scm.com) | Version control system |
| **Markdown** | [Glow](https://github.com/charmbracelet/glow) | terminal markdown reader |
| **Markdown** | [Glamour](https://github.com/charmbracelet/glamour) | Markdown rendering library |
| **Input** | [Gum](https://github.com/charmbracelet/gum) | Glamorous shell scripts |
| **File Viewer** | [bat](https://github.com/sharkdp/bat) | Cat clone with syntax highlighting |
| **Fuzzy Finder** | [fzf](https://github.com/junegunn/fzf) | Command-line fuzzy finder |
| **Linter** | [golangci-lint](https://golangci-lint.run) | Go linters aggregator |

All tools are automatically installed and configured when you run `./setup.sh`.

## Development

### Prerequisites

- Go 1.21+

### Build

```bash
make build

go build -o orbit .
```

### Install

```bash
make install
```

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Charm](https://charm.sh/) for the amazing TUI libraries
- [Cobra](https://github.com/spf13/cobra) for CLI framework
- [Glamour](https://github.com/charmbracelet/glamour) for markdown rendering

---

Made with by [henrynguci](https://github.com/henrynguci)
