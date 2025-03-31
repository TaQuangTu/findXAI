#!/bin/bash

set -e  # Exit immediately if any command fails

# Function to check if a command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

echo "Checking system requirements..."

# Install Python if not installed
if ! command_exists python3; then
    echo "Python is not installed. Installing Python..."
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        sudo apt update && sudo apt install -y python3 python3-pip
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        brew install python3
    else
        echo "Unsupported OS. Please install Python manually."
        exit 1
    fi
fi

# Install pip if not installed
if ! command_exists pip3; then
    echo "pip3 is not installed. Installing pip..."
    python3 -m ensurepip --default-pip
fi

# Install pre-commit if not installed
if ! command_exists pre-commit; then
    echo "Installing pre-commit..."
    pip3 install pre-commit
else
    echo "pre-commit is already installed."
fi

# # Navigate to the Go project directory (assumed to be the script's location)
# cd "$(dirname "$0")"

# Install pre-commit hooks
if [[ -f ".pre-commit-config.yaml" ]]; then
    echo "Installing pre-commit hooks..."
    pre-commit install
else
    echo "No .pre-commit-config.yaml found. Skipping hook installation."
    exit 1
fi

echo "✅ Pre-commit setup completed successfully!"
