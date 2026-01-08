# PulseFetch

PulseFetch is a minimalist, Neofetch-inspired program written in Go. It displays system information along with a custom image or logo, designed for terminal enthusiasts.

## Features

- Display essential system information:
  - OS, Host, Kernel, Uptime, Packages, Shell, Resolution, WM, Terminal, CPU, GPU, Memory, and Disk.
- Image rendering:
  - Uses Chafa for terminal graphics.
  - Supports Kitty Graphics Protocol, Sixel, and Unicode block symbols fallback.
  - Automatically detects terminal capabilities for the best possible image quality.
- Customizable configuration:
  - TOML-based configuration file.
  - Support for custom images and logos.
  - Toggle individual information modules.
- Lightweight and fast.

## Requirements

- Go 1.25.3 or later.
- Chafa (optional, for image support).

## Installation

### Debian / Ubuntu (APT)

You can install PulseFetch directly from the automated APT repository:

```bash
# 1. Add the repository
echo "deb [trusted=yes] https://maskedsyntax.github.io/pulsefetch/apt/ /" | sudo tee /etc/apt/sources.list.d/pulsefetch.list

# 2. Update and install
sudo apt update && sudo apt install pulsefetch
```

### From Source

1. Clone the repository.
2. Build the binary:
   ```bash
   go build -o pulsefetch ./cmd/pulsefetch/main.go
   ```
3. Run the program:
   ```bash
   ./pulsefetch
   ```

## Configuration

PulseFetch looks for a configuration file at `~/.config/pulsefetch/pulsefetch.toml` or `/etc/pulsefetch/pulsefetch.toml`. A default template is provided in the repository.
