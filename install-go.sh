#!/bin/bash

# Function to check if a command exists
command_exists() {
    command -v "$@" > /dev/null 2>&1
}

# Install curl if not already installed
if ! command_exists curl; then
    if command_exists apt; then
        sudo apt install -y curl
    elif command_exists yum; then
        sudo yum install -y curl
    elif command_exists dnf; then
        sudo dnf install -y curl
    else
        echo "Package manager not supported. Install curl manually."
        exit 1
    fi
fi

# Install curl
sudo apt install curl

# Get CPU architecture
cpuArch=$(arch)
if [ "$cpuArch" = "aarch64" ]; then
	    cpuArch="arm64"
fi

# Get Operating system
osys=$(uname)
# Latest Go programming language version
release=$(curl --silent https://go.dev/doc/devel/release | grep -Eo 'go[0-9]+(\.[0-9]+)+' | sort -V | uniq | tail -1 | sed 's/go//')

# Prompt the user for the Go version
read -p "Enter the Go version you want to install (e.g., 1.20.1) or press enter to install the latest version ($release): " GO_VERSION

# Use latest version if the user did not enter a version
if [[ -z "$GO_VERSION" ]]; then
    GO_VERSION=$release
    echo "Installing latest version of Go: $GO_VERSION"
fi

# Download Go for the specified architecture
curl -OL https://dl.google.com/go/go${GO_VERSION}.${osys}-${cpuArch}.tar.gz

b=go${GO_VERSION}.${osys}-${cpuArch}.tar.gz
du -hcs $b

# Check if download was successful
if [ $? -ne 0 ]; then
    echo "Download failed. Check the entered version and try again."
    exit 1
fi

# Extract the tarball
sudo tar -C /usr/local -xzf go${GO_VERSION}.${osys}-${cpuArch}.tar.gz

# Remove the tarball
rm go${GO_VERSION}.${osys}-${cpuArch}.tar.gz

# Set Go environment variables
echo "export PATH=\$PATH:/usr/local/go/bin" >> $HOME/.zshrc
echo "export GOROOT=/usr/local/go" >> $HOME/.zshrc
echo "export GOPATH=\$HOME/go" >> $HOME/.zshrc

echo "export PATH=\$PATH:/usr/local/go/bin" >> $HOME/.bashrc
echo "export GOROOT=/usr/local/go" >> $HOME/.bashrc
echo "export GOPATH=\$HOME/go" >> $HOME/.bashrc

# Apply environment variables
source $HOME/.zshrc
source $HOME/.bashrc
# Print Go version to verify installation
go version