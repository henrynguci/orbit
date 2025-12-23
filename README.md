<p align="center">
  <img src="images/logo.png" alt="Orbit Logo">
</p>

<p align="center"><i>Keep your side projects in orbit</i></p>

Orbit is a powerful CLI tool and TUI application for managing your side projects. Built with Go, using [Cobra](https://github.com/spf13/cobra) for CLI, [Bubble Tea](https://github.com/charmbracelet/bubbletea) for TUI, and [Lip Gloss](https://github.com/charmbracelet/lipgloss) for styling.

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go" alt="Go Version">
  <img src="https://img.shields.io/badge/license-MIT-green" alt="License">
</p>

## Features

- **Workspace Management** - Initialize and manage multiple workspaces
- **Project Organization** - Create projects with repo, docs, and secret folders
- **Status Tracking** - Track project status (active, archived, done)
- **Aliases** - Set short aliases for projects with long names
- **README Viewer** - Beautiful markdown rendering in terminal

## Installation

![Orbit Setup](images/orbit.gif)

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
- [Glow](https://github.com/charmbracelet/glow) for markdown rendering

---

Made with by [henrynguci](https://github.com/henrynguci)
