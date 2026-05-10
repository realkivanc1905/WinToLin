# WinToLin

WinToLin is a high-performance, P2P file transfer utility designed to seamlessly move files between Windows and Linux machines. It works over Local Networks (Wi-Fi) and the Internet.

## Features

- **Blazing Fast**: Direct TCP connection for maximum speed.
- **Cross-Platform**: Works on Windows and Linux with zero dependencies.
- **Progress Tracking**: Real-time speed and percentage display.
- **Internet Support**: Automatically detects your Public IP for remote transfers.

## Installation

### Windows
Open PowerShell and run:
```powershell
irm https://raw.githubusercontent.com/realkivanc1905/WinToLin/main/install.ps1 | iex
```
*Restart your terminal after installation.*

### Linux (Ubuntu/Debian)
Open Terminal and run:
```bash
curl -sL https://raw.githubusercontent.com/realkivanc1905/WinToLin/main/install.sh | bash
```

## Usage

### 1. Send a File (Windows or Linux)
```bash
wintolin send my_file.zip
```
The program will display a command for the receiver to run.

### 2. Receive a File (Windows or Linux)
Run the command provided by the sender, for example:
```bash
wintolin receive 192.168.1.10:54321
```

## Troubleshooting
- **Connection Refused**: Ensure the sender is still running and waiting for the receiver.
- **Over Internet**: If transferring over the internet, you may need to open the displayed port on your router's Port Forwarding settings.
