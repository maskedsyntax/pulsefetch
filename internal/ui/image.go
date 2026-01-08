package ui

import (
	"fmt"
	"os/exec"
	"strings"
)

func RenderImage(path string, terminalName string) (string, int, int, error) {
	width := 40
	height := 20
	size := "40x20"
	
	args := []string{"-s", size, "--animate=off", "--polite=on"}

	term := strings.ToLower(terminalName)

	if strings.Contains(term, "kitty") {
		args = append(args, "-f", "kitty")
	} else if strings.Contains(term, "wezterm") || strings.Contains(term, "foot") || strings.Contains(term, "mlterm") {
		args = append(args, "-f", "sixel")
	} else if strings.Contains(term, "ghostty") {
		args = append(args, "-f", "kitty")
	} else {
		// Fallback for Alacritty, Gnome-Terminal, etc.
		args = append(args, "-f", "symbols")
	}

	args = append(args, path)

	if _, err := exec.LookPath("chafa"); err != nil {
		return "", 0, 0, fmt.Errorf("chafa executable not found in PATH")
	}

	cmd := exec.Command("chafa", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", 0, 0, fmt.Errorf("chafa failed: %v, output: %s", err, string(out))
	}

	return string(out), width, height, nil
}
