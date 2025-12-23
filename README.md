# 🌍 Orbit

```
  ██████╗  ██████╗  ██████╗  ██╗ ████████╗
 ██╔═══██╗ ██╔══██╗ ██╔══██╗ ██║ ╚══██╔══╝
 ██║   ██║ ██████╔╝ ██████╔╝ ██║    ██║   
 ██║   ██║ ██╔══██╗ ██╔══██╗ ██║    ██║   
 ╚██████╔╝ ██║  ██║ ██████╔╝ ██║    ██║   
  ╚═════╝  ╚═╝  ╚═╝ ╚═════╝  ╚═╝    ╚═╝   
```

> Keep your side projects in orbit 🚀

Orbit is a powerful CLI tool and TUI application for managing your side projects. Built with Go, using [Cobra](https://github.com/spf13/cobra) for CLI, [Bubble Tea](https://github.com/charmbracelet/bubbletea) for TUI, and [Lip Gloss](https://github.com/charmbracelet/lipgloss) for styling.

![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/license-MIT-green)

## ✨ Features

- **Beautiful TUI** - Interactive terminal UI with ASCII banner
- **Workspace Management** - Initialize and manage multiple workspaces
- **Project Organization** - Create projects with repo, docs, and secret folders
- **Status Tracking** - Track project status (active, archived, done)
- **Aliases** - Set short aliases for projects with long names
- **README Viewer** - Beautiful markdown rendering in terminal

##  Installation

### From Source

```bash
git clone https://github.com/henrynguci/orbit.git
cd orbit
go build -o orbit .
```

### Move to PATH

```bash
sudo mv orbit /usr/local/bin/
```

## 🚀 Usage

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
# Create workspace only
orbit init ~/my-workspace

# Create workspace with a project
orbit init ~/my-workspace --project my-project
orbit init ~/my-workspace -p my-project
```

With project, creates:
```
~/my-workspace/
└── project/
    └── my-project/
        ├── repo/      # Your project repositories
        ├── docs/      # Documentation files
        └── secret/    # Secret/sensitive files
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

# Examples:
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

## 🎨 TUI Features

- **ASCII Banner** - Beautiful ORBIT logo on startup
- **Menu Navigation** - Use arrow keys or j/k
- **Multi-step Forms** - Guided workspace/project creation
- **Project List** - Filterable with status icons
  - 🟢 Active
  -  Archived
  -  Done
- **README Viewer** - Glamour markdown rendering with scroll

## ⌨️ Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `↑/↓` or `j/k` | Navigate |
| `Enter` | Select/Confirm |
| `/` | Filter (in lists) |
| `Esc` | Go back |
| `q` | Quit |

## 📁 Configuration

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

## 🛠️ Development

### Prerequisites

- Go 1.21+

### Build

```bash
make build
# or
go build -o orbit .
```

### Install

```bash
make install
```

## 📝 License

MIT License - see [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Charm](https://charm.sh/) for the amazing TUI libraries
- [Cobra](https://github.com/spf13/cobra) for CLI framework
- [Glamour](https://github.com/charmbracelet/glamour) for markdown rendering

---

Made with ❤️ by [henrynguci](https://github.com/henrynguci)
