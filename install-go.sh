#!/bin/bash

# Install curl
sudo apt install curl

# Prompt the user for the Go version
read -p "Enter the Go version you want to install (e.g., 1.18.1): " GO_VERSION

# Get CPU architecture
cpuArch = arch

# Latest Go programming language
release=$(curl --silent https://go.dev/doc/devel/release | grep -Eo 'go[0-9]+(\.[0-9]+)+' | sort -V | uniq | tail -1 | sed 's/go//') && echo $release

# Check if the user entered a version
if [[ -z "$GO_VERSION" ]]; then
    echo "No version entered. Exiting."
    exit 1
fi

# Download Go for ARM64
curl -O  https://dl.google.com/go/go${GO_VERSION}.linux-${cpuArch}.tar.gz


# Check if download was successful
if [ $? -ne 0 ]; then
    echo "Download failed. Check the entered version and try again."
    exit 1
fi

# Extract the tarball
sudo tar -C /usr/local -xzf go${GO_VERSION}.linux-arm64.tar.gz

# Remove the tarball
rm go${GO_VERSION}.linux-arm64.tar.gz

# Set Go environment variables
echo "export PATH=\$PATH:/usr/local/go/bin" >> $HOME/.profile
echo "export GOROOT=/usr/local/go" >> $HOME/.profile
echo "export GOPATH=\$HOME/go" >> $HOME/.profile

# Apply environment variables
source $HOME/.profile

# Print Go version to verify installation
go version
