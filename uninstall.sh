#!/bin/bash

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

print_step() {
    echo ""
    echo -e "${CYAN}$1${NC}"
    echo "----------------------------------------"
}

print_success() {
    echo -e "${GREEN}$1${NC}"
}

print_info() {
    echo -e "${BLUE}$1${NC}"
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

uninstall_orbit() {
    print_step "Uninstalling Orbit"
    
    if [[ -f "/usr/local/bin/orbit" ]]; then
        print_info "Removing orbit binary from /usr/local/bin..."
        sudo rm "/usr/local/bin/orbit"
        print_success "Orbit binary removed"
    fi
    
    if [[ -d "$HOME/.config/orbit" ]]; then
        print_info "Removing configuration directory..."
        rm -rf "$HOME/.config/orbit"
        print_success "Orbit configuration removed"
    fi

    if [[ -d "bin" ]]; then
        print_info "Removing local build directory..."
        rm -rf bin
    fi
}

uninstall_glow() {
    print_step "Uninstalling Glow"
    
    if ! command_exists glow; then
        print_info "Glow not found"
        return 0
    fi
    
    case $OS in
        debian)
            sudo apt-get remove -y glow
            sudo rm -f /etc/apt/sources.list.d/charm.list
            sudo rm -f /etc/apt/keyrings/charm.gpg
            ;;
        fedora|rhel)
            sudo yum remove -y glow
            sudo rm -f /etc/yum.repos.d/charm.repo
            ;;
        arch)
            sudo pacman -Rs --noconfirm glow
            ;;
        macos)
            brew uninstall glow
            ;;
        *)
            if [[ -f "$HOME/go/bin/glow" ]]; then
                rm "$HOME/go/bin/glow"
            fi
            ;;
    esac
    print_success "Glow uninstalled"
}

uninstall_go() {
    print_step "Uninstalling Go"
    
    if [[ -d "/usr/local/go" ]]; then
        print_info "Removing /usr/local/go..."
        sudo rm -rf /usr/local/go
        print_success "Go binaries removed"
    fi
    
    print_info "Cleaning up PATH in shell configs..."
    sed -i '/\/usr\/local\/go\/bin/d' ~/.bashrc
    sed -i '/\$HOME\/go\/bin/d' ~/.bashrc
    
    if [[ -f ~/.zshrc ]]; then
        sed -i '/\/usr\/local\/go\/bin/d' ~/.zshrc
        sed -i '/\$HOME\/go\/bin/d' ~/.zshrc
    fi
    
    print_success "Shell configs cleaned"
}

main() {
    detect_os
    
    uninstall_orbit
    uninstall_glow
    uninstall_go
    
    echo ""
    echo -e "${GREEN}╔════════════════════════════════════════════════╗${NC}"
    echo -e "${GREEN}║                                                ║${NC}"
    echo -e "${GREEN}║  All components uninstalled successfully!     ║${NC}"
    echo -e "${GREEN}║                                                ║${NC}"
    echo -e "${GREEN}╚════════════════════════════════════════════════╝${NC}"
    echo ""
}

main
