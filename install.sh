#!/bin/bash
set -e

REPO="realkivanc1905/WinToLin"
echo "Downloading WinToLin ($REPO)..."

# Fetch latest release URL
URL=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep -Eo 'https://[^"]+wintolin-linux-amd64' | head -1)

if [ -z "$URL" ]; then
    echo "ERROR: Latest release not found. Please ensure you have created a Release on GitHub."
    exit 1
fi

echo "Downloading: $URL..."
curl -sL "$URL" -o wintolin
chmod +x wintolin

if [ -w /usr/local/bin ]; then
    mv wintolin /usr/local/bin/wintolin
    echo "wintolin installed to /usr/local/bin"
else
    mkdir -p ~/.local/bin
    mv wintolin ~/.local/bin/wintolin
    echo "wintolin installed to ~/.local/bin"
    
    # Auto-add to PATH
    if [[ ":$PATH:" != *":$HOME/.local/bin:"* ]]; then
        echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc
        if [ -f ~/.zshrc ]; then
            echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.zshrc
        fi
        echo "PATH configuration added to ~/.bashrc"
        echo "Please restart your terminal or run: source ~/.bashrc"
    fi
fi

echo "Installation complete! You can now use the 'wintolin' command."
