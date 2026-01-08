package ui

import (
	"fmt"
	"strings"

	"pulsefetch/internal/config"
	"pulsefetch/internal/fetcher"

	"github.com/charmbracelet/lipgloss"
)

var (
	keyStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("39")).Bold(true) // Removed fixed padding
	valueStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("255")) // White
	logoStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("39")).Bold(true).PaddingRight(4)
	titleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("39")).Bold(true)
	sepStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("255"))
)

type infoItem struct {
	Key   string
	Value string
	Multi []string // For multi-line values like GPUs
}

func Render(cfg *config.Config, info *fetcher.SystemInfo, logo string) {
	var items []infoItem

	// Helper to collect items
	add := func(key, value string) {
		if value != "" {
			items = append(items, infoItem{Key: key, Value: value})
		}
	}
	addMulti := func(key string, values []string) {
		if len(values) > 0 {
			items = append(items, infoItem{Key: key, Multi: values})
		}
	}

	if cfg.ShowOS { add("OS", info.OS) }
	if cfg.ShowHost { add("Host", info.Host) }
	if cfg.ShowKernel { add("Kernel", info.Kernel) }
	if cfg.ShowUptime { add("Uptime", info.Uptime) }
	if cfg.ShowPackages { add("Packages", info.Packages) }
	if cfg.ShowShell { add("Shell", info.Shell) }
	if cfg.ShowResolution { add("Resolution", info.Resolution) }
	if cfg.ShowDE { add("DE", info.DE) }
	if cfg.ShowWM { add("WM", info.WM) }
	if cfg.ShowWMTheme { add("WM Theme", info.WMTheme) }
	if cfg.ShowTheme { add("Theme", info.Theme) }
	if cfg.ShowIcons { add("Icons", info.Icons) }
	if cfg.ShowTerminal { add("Terminal", info.Terminal) }
	if cfg.ShowCPU { add("CPU", info.CPU) }
	if cfg.ShowCPUUsage { add("CPU Usage", info.CPUUsage) }
	
	if cfg.ShowGPU { addMulti("GPU", info.GPUs) }
	
	if cfg.ShowMemory { add("Memory", info.Memory) }
	if cfg.ShowMemoryUsage { add("Memory Usage", info.MemoryUsage) }
	if cfg.ShowDisk { add("Disk", info.Disk) }
	if cfg.ShowDiskUsage { add("Disk Usage", info.DiskUsage) }
	if cfg.ShowNetwork { add("Network", info.Network) }
	if cfg.ShowBattery { add("Battery", info.Battery) }
	if cfg.ShowSensors { add("Sensors", info.Sensors) }

	// Calculate Max Key Length
	maxKeyLen := 0
	for _, item := range items {
		if len(item.Key) > maxKeyLen {
			maxKeyLen = len(item.Key)
		}
	}

	var infoLines []string

	// Title
	if info.User != "" && info.Hostname != "" {
		title := fmt.Sprintf("%s@%s", info.User, info.Hostname)
		infoLines = append(infoLines, titleStyle.Render(title))
		sep := strings.Repeat("-", len(title))
		infoLines = append(infoLines, sepStyle.Render(sep))
	} else if info.Hostname != "" {
		infoLines = append(infoLines, titleStyle.Render(info.Hostname))
		infoLines = append(infoLines, sepStyle.Render(strings.Repeat("-", len(info.Hostname))))
	}

	// Render Items
	for _, item := range items {
		// Create padding string
		padLen := maxKeyLen - len(item.Key)
		padding := strings.Repeat(" ", padLen)
		
		// Render Key
		keyStr := fmt.Sprintf("%s%s  ", item.Key, padding) // Key + align padding + 2 spaces gap
		styledKey := keyStyle.Render(keyStr)

		if len(item.Multi) > 0 {
			// Multi-line (GPUs)
			// First line: Key ... Value[0]
			infoLines = append(infoLines, fmt.Sprintf("%s%s", styledKey, valueStyle.Render(item.Multi[0])))
			
			// Subsequent lines: Empty Key ... Value[i]
			emptyKeyPadding := strings.Repeat(" ", maxKeyLen + 2) // MaxLen + 2 spaces gap
			for i := 1; i < len(item.Multi); i++ {
				infoLines = append(infoLines, fmt.Sprintf("%s%s", emptyKeyPadding, valueStyle.Render(item.Multi[i])))
			}
		} else {
			// Single line
			infoLines = append(infoLines, fmt.Sprintf("%s%s", styledKey, valueStyle.Render(item.Value)))
		}
	}

	infoBlock := lipgloss.JoinVertical(lipgloss.Left, infoLines...)

	// If Image Mode (detected by config presence and non-empty logo)
	if cfg.ImagePath != "" && logo != "" {
		// Graphic Image Mode Layout
		// 1. Print Image
		fmt.Print(logo)

		// 2. Move Cursor UP to top of image
		// Calculate actual height of the logo
		logoHeight := strings.Count(logo, "\n")
		if logoHeight > 0 {
			fmt.Printf("\033[%dA", logoHeight)
		}

		// 3. Print Info Block Line by Line
		lines := strings.Split(infoBlock, "\n")
		for _, line := range lines {
			// Move Right 42 (40 width + 2 padding)
			fmt.Printf("\033[%dC%s\n", 42, line)
		}
		
		// If info has fewer lines than image, we need to move cursor down
		if len(lines) < logoHeight {
			remaining := logoHeight - len(lines)
			fmt.Printf("\033[%dB", remaining)
		}
	} else {
		// Standard Text/ASCII Layout
		styledLogo := logoStyle.Render(logo)
		finalOutput := lipgloss.JoinHorizontal(lipgloss.Top, styledLogo, infoBlock)
		fmt.Println(finalOutput)
	}
}
