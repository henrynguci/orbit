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
    
    local MIN_VERSION="1.21"
    local TARGET_VERSION="1.22.0"
    
    if command_exists go; then
        local current_version=$(go version | awk '{print $3}' | sed 's/go//')
        local current_major_minor=$(echo "$current_version" | cut -d. -f1,2)
        
        if awk "BEGIN {exit !($current_major_minor >= $MIN_VERSION)}"; then
            print_success "Go $current_version is already installed (>= $MIN_VERSION required)"
            return 0
        else
            print_warning "Go $current_version found but < $MIN_VERSION, upgrading..."
        fi
    fi
    
    print_info "Installing Go $TARGET_VERSION..."
    
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
        local glow_version=$(glow --version 2>/dev/null | head -n1 | awk '{print $NF}' || echo "unknown")
        print_success "Glow $glow_version is already installed"
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

install_make() {
    print_step "Checking Make"
    
    if command_exists make; then
        local make_version=$(make --version 2>/dev/null | head -n1 | awk '{print $3}' || echo "unknown")
        print_success "Make $make_version is already installed"
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

main() {
    print_header
    
    print_info "Detecting operating system..."
    detect_os
    print_success "Detected OS: $OS"
    
    echo ""
    print_info "Starting installation of dependencies..."
    echo ""
    
    install_make
    install_go
    
    install_glow
    
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
    echo ""
    if [[ -n "$SHELL" ]]; then
        exec "$SHELL" -l
    else
        exec bash -l
    fi
}

main
