#!/bin/bash

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

print_header() {
    echo -e "${CYAN}"
    echo "  ██████╗  ██████╗  ██████╗  ██╗ ████████╗"
    echo " ██╔═══██╗ ██╔══██╗ ██╔══██╗ ██║ ╚══██╔══╝"
    echo " ██║   ██║ ██████╔╝ ██████╔╝ ██║    ██║   "
    echo " ██║   ██║ ██╔══██╗ ██╔══██╗ ██║    ██║   "
    echo " ╚██████╔╝ ██║  ██║ ██████╔╝ ██║    ██║   "
    echo "  ╚═════╝  ╚═╝  ╚═╝ ╚═════╝  ╚═╝    ╚═╝   "
    echo ""
    echo "         Setup Dependencies Script"
    echo -e "${NC}"
}

print_success() {
    echo -e "${GREEN}$1${NC}"
}

print_error() {
    echo -e "${RED}$1${NC}"
}

print_warning() {
    echo -e "${YELLOW}$1${NC}"
}

print_info() {
    echo -e "${BLUE}$1${NC}"
}

print_step() {
    echo ""
    echo -e "${CYAN}$1${NC}"
    echo "----------------------------------------"
}

command_exists() {
    command -v "$1" >/dev/null 2>&1
}

detect_os() {
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        if command_exists apt-get; then
            OS="debian"
        elif command_exists dnf; then
            OS="fedora"
        elif command_exists yum; then
            OS="rhel"
        elif command_exists pacman; then
            OS="arch"
        else
            OS="linux"
        fi
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        OS="macos"
    else
        OS="unknown"
    fi
}

install_go() {
    print_step "Installing Go"
    
    if command_exists go; then
        local current_version=$(go version | awk '{print $3}' | sed 's/go//')
        print_success "Go is already installed (version $current_version)"
        return 0
    fi
    
    print_info "Installing Go 1.21+..."
    
    case $OS in
        debian)
            sudo apt-get update
            sudo apt-get install -y wget tar
            ;;
        macos)
            if ! command_exists brew; then
                print_error "Homebrew not found. Please install Homebrew first."
                return 1
            fi
            brew install go
            print_success "Go installed via Homebrew"
            return 0
            ;;
    esac
    
    local GO_VERSION="1.22.0"
    local ARCH=$(uname -m)
    
    if [[ "$ARCH" == "x86_64" ]]; then
        ARCH="amd64"
    elif [[ "$ARCH" == "aarch64" ]]; then
        ARCH="arm64"
    fi
    
    local GO_TARBALL="go${GO_VERSION}.linux-${ARCH}.tar.gz"
    local GO_URL="https://go.dev/dl/${GO_TARBALL}"
    
    print_info "Downloading Go ${GO_VERSION}..."
    wget -q --show-progress "$GO_URL" -O "/tmp/${GO_TARBALL}"
    
    print_info "Extracting Go..."
    sudo rm -rf /usr/local/go
    sudo tar -C /usr/local -xzf "/tmp/${GO_TARBALL}"
    rm "/tmp/${GO_TARBALL}"
    
    if [[ ! "$PATH" =~ "/usr/local/go/bin" ]]; then
        echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
        echo 'export PATH=$PATH:$HOME/go/bin' >> ~/.bashrc
        export PATH=$PATH:/usr/local/go/bin
        export PATH=$PATH:$HOME/go/bin
        
        if [[ -f ~/.zshrc ]]; then
            echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.zshrc
            echo 'export PATH=$PATH:$HOME/go/bin' >> ~/.zshrc
        fi
    fi
    
    print_success "Go installed successfully"
}

install_glow() {
    print_step "Installing Glow (Markdown Viewer)"
    
    if command_exists glow; then
        print_success "Glow is already installed"
        return 0
    fi
    
    print_info "Installing Glow..."
    
    case $OS in
        debian)
            sudo mkdir -p /etc/apt/keyrings
            curl -fsSL https://repo.charm.sh/apt/gpg.key | sudo gpg --dearmor -o /etc/apt/keyrings/charm.gpg
            echo "deb [signed-by=/etc/apt/keyrings/charm.gpg] https://repo.charm.sh/apt/ * *" | sudo tee /etc/apt/sources.list.d/charm.list
            sudo apt-get update
            sudo apt-get install -y glow
            ;;
        fedora|rhel)
            echo '[charm]
name=Charm
baseurl=https://repo.charm.sh/yum/
enabled=1
gpgcheck=1
gpgkey=https://repo.charm.sh/yum/gpg.key' | sudo tee /etc/yum.repos.d/charm.repo
            sudo yum install -y glow
            ;;
        arch)
            sudo pacman -S --noconfirm glow
            ;;
        macos)
            brew install glow
            ;;
        *)
            go install github.com/charmbracelet/glow@latest
            ;;
    esac
    
    print_success "Glow installed successfully"
}

install_gp() {
    print_step "Installing gp (Charm's Git Tool)"
    
    if command_exists gp; then
        print_success "gp is already installed"
        return 0
    fi
    
    print_info "Installing gp..."
    
    case $OS in
        debian)
            sudo apt-get install -y gum
            ;;
        fedora|rhel)
            sudo yum install -y gum
            ;;
        arch)
            sudo pacman -S --noconfirm gum
            ;;
        macos)
            brew install gum
            ;;
        *)
            go install github.com/charmbracelet/gum@latest
            ;;
    esac
    
    print_success "Charm tools installed successfully"
}

install_bat() {
    print_step "Installing bat (Better cat)"
    
    if command_exists bat || command_exists batcat; then
        print_success "bat is already installed"
        return 0
    fi
    
    print_info "Installing bat..."
    
    case $OS in
        debian)
            sudo apt-get install -y bat
            if [[ ! -f ~/.local/bin/bat ]] && command_exists batcat; then
                mkdir -p ~/.local/bin
                ln -s /usr/bin/batcat ~/.local/bin/bat
            fi
            ;;
        fedora|rhel)
            sudo yum install -y bat
            ;;
        arch)
            sudo pacman -S --noconfirm bat
            ;;
        macos)
            brew install bat
            ;;
    esac
    
    print_success "bat installed successfully"
}

install_fzf() {
    print_step "Installing fzf (Fuzzy Finder)"
    
    if command_exists fzf; then
        print_success "fzf is already installed"
        return 0
    fi
    
    print_info "Installing fzf..."
    
    case $OS in
        debian)
            sudo apt-get install -y fzf
            ;;
        fedora|rhel)
            sudo yum install -y fzf
            ;;
        arch)
            sudo pacman -S --noconfirm fzf
            ;;
        macos)
            brew install fzf
            ;;
        *)
            git clone --depth 1 https://github.com/junegunn/fzf.git ~/.fzf
            ~/.fzf/install --all
            ;;
    esac
    
    print_success "fzf installed successfully"
}

install_golangci_lint() {
    print_step "Installing golangci-lint"
    
    if command_exists golangci-lint; then
        print_success "golangci-lint is already installed"
        return 0
    fi
    
    print_info "Installing golangci-lint..."
    
    curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
    
    print_success "golangci-lint installed successfully"
}

install_git() {
    print_step "Checking Git"
    
    if command_exists git; then
        print_success "Git is already installed"
        return 0
    fi
    
    print_info "Installing Git..."
    
    case $OS in
        debian)
            sudo apt-get install -y git
            ;;
        fedora|rhel)
            sudo yum install -y git
            ;;
        arch)
            sudo pacman -S --noconfirm git
            ;;
        macos)
            brew install git
            ;;
    esac
    
    print_success "Git installed successfully"
}

install_make() {
    print_step "Checking Make"
    
    if command_exists make; then
        print_success "Make is already installed"
        return 0
    fi
    
    print_info "Installing Make..."
    
    case $OS in
        debian)
            sudo apt-get install -y build-essential
            ;;
        fedora|rhel)
            sudo yum groupinstall -y "Development Tools"
            ;;
        arch)
            sudo pacman -S --noconfirm base-devel
            ;;
        macos)
            xcode-select --install
            ;;
    esac
    
    print_success "Make installed successfully"
}

install_glamour() {
    print_step "Installing Glamour (CLI Markdown Renderer)"
    
    if command_exists glamour; then
        print_success "Glamour is already installed"
        return 0
    fi
    
    print_info "Installing Glamour..."
    go install github.com/charmbracelet/glamour@latest
    
    print_success "Glamour installed successfully"
}

main() {
    print_header
    
    print_info "Detecting operating system..."
    detect_os
    print_success "Detected OS: $OS"
    
    echo ""
    print_info "Starting installation of dependencies..."
    echo ""
    
    install_git
    install_make
    install_go
    
    install_golangci_lint
    
    install_glow
    install_gp
    install_glamour
    
    install_bat
    install_fzf
    
    print_step "Installing Go project dependencies"
    if [[ -f "go.mod" ]]; then
        print_info "Running go mod download..."
        go mod download
        print_info "Running go mod tidy..."
        go mod tidy
        print_success "Go modules installed successfully"
    else
        print_warning "go.mod not found, skipping Go module installation"
    fi
    
    echo ""
    echo -e "${GREEN}╔════════════════════════════════════════════════╗${NC}"
    echo -e "${GREEN}║                                                ║${NC}"
    echo -e "${GREEN}║  All dependencies installed successfully!     ║${NC}"
    echo -e "${GREEN}║                                                ║${NC}"
    echo -e "${GREEN}╚════════════════════════════════════════════════╝${NC}"
    echo ""
    
    print_info "You can now build the project with:"
    echo "  make build"
    echo ""
    print_info "Or install it system-wide with:"
    echo "  make install"
    echo ""
    
    print_warning "Please restart your shell or run: source ~/.bashrc (or ~/.zshrc)"
    echo ""
}

main
