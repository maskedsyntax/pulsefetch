package fetcher

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"strings"
	"time"

	"pulsefetch/internal/config"

	//"github.com/shirou/gopsutil/v3/battery"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

type SystemInfo struct {
	User         string
	OS           string
	Hostname     string
	Host         string // Hardware Model
	Kernel       string
	Uptime       string
	Packages     string
	Shell        string
	Resolution   string
	DE           string
	WM           string
	WMTheme      string
	Theme        string
	Icons        string
	Terminal     string
	CPU          string
	GPUs         []string
	Memory       string
	Disk         string
	Network      string 
	Battery      string
	Sensors      string
	
	CPUUsage     string
	MemoryUsage  string
	DiskUsage    string
}

func Fetch(cfg *config.Config) (*SystemInfo, error) {
	info := &SystemInfo{}

	u, err := user.Current()
	if err == nil {
		info.User = u.Username
	}

	h, err := host.Info()
	if err == nil {
		info.Hostname = h.Hostname
		if cfg.ShowOS {
			info.OS = fmt.Sprintf("%s %s", h.Platform, h.PlatformVersion)
		}
		if cfg.ShowKernel {
			info.Kernel = h.KernelVersion
		}
		if cfg.ShowUptime {
			d := time.Duration(h.Uptime) * time.Second
			info.Uptime = formatDuration(d)
		}
	}

	if cfg.ShowHost {
		info.Host = getModel()
	}

	if cfg.ShowCPU || cfg.ShowCPUUsage {
		c, err := cpu.Info()
		if err == nil && len(c) > 0 {
			info.CPU = c[0].ModelName
		}
		if cfg.ShowCPUUsage {
			percent, err := cpu.Percent(0, false)
			if err == nil && len(percent) > 0 {
				info.CPUUsage = fmt.Sprintf("%.1f%%", percent[0])
			}
		}
	}

	if cfg.ShowGPU {
		info.GPUs = getGPU()
	}

	if cfg.ShowResolution {
		info.Resolution = getResolution()
	}

	if cfg.ShowMemory || cfg.ShowMemoryUsage {
		v, err := mem.VirtualMemory()
		if err == nil {
			if cfg.ShowMemory {
				info.Memory = fmt.Sprintf("%vMiB / %vMiB", v.Used/1024/1024, v.Total/1024/1024)
			}
			if cfg.ShowMemoryUsage {
				info.MemoryUsage = fmt.Sprintf("%.1f%%", v.UsedPercent)
			}
		}
	}

	if cfg.ShowDisk || cfg.ShowDiskUsage {
		parts, err := disk.Partitions(false)
		if err == nil && len(parts) > 0 {
			u, err := disk.Usage("/")
			if err == nil {
				if cfg.ShowDisk {
					info.Disk = fmt.Sprintf("%vGiB / %vGiB", u.Used/1024/1024/1024, u.Total/1024/1024/1024)
				}
				if cfg.ShowDiskUsage {
					info.DiskUsage = fmt.Sprintf("%.1f%%", u.UsedPercent)
				}
			}
		}
	}
	
	if cfg.ShowNetwork {
		ifaces, err := net.Interfaces()
		if err == nil {
			for _, i := range ifaces {
				isLoopback := false
				for _, flag := range i.Flags {
					if flag == "loopback" {
						isLoopback = true
						break
					}
				}
				if isLoopback { continue }
				for _, addr := range i.Addrs {
					if strings.Contains(addr.Addr, ".") {
						info.Network = addr.Addr
						break
					}
				}
				if info.Network != "" { break }
			}
		}
	}

	if cfg.ShowShell {
		info.Shell = os.Getenv("SHELL")
		if info.Shell != "" {
			parts := strings.Split(info.Shell, "/")
			info.Shell = parts[len(parts)-1]
		}
	}

	if cfg.ShowTerminal {
		info.Terminal = getTerminal()
	}

	if cfg.ShowDE {
		info.DE = getDE()
	}
	
	if cfg.ShowWM {
		wms := map[string]bool{
			"i3": true, "bspwm": true, "sway": true, "dwm": true, "awesome": true, "xmonad": true, "openbox": true,
		}
		if wms[info.DE] {
			info.WM = info.DE
			info.DE = "" 
		} else {
			info.WM = getWM() 
		}
	}

	if cfg.ShowWMTheme {
		info.WMTheme = "" 
	}

	if cfg.ShowTheme {
		info.Theme = getTheme()
	}

	if cfg.ShowIcons {
		info.Icons = getIcons()
	}

	if cfg.ShowPackages {
		info.Packages = getPackages()
	}

	return info, nil
}

func formatDuration(d time.Duration) string {
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	return fmt.Sprintf("%d hours, %d mins", h, m)
}

func getModel() string {
	data, err := os.ReadFile("/sys/class/dmi/id/product_name")
	if err == nil {
		return strings.TrimSpace(string(data))
	}
	data, err = os.ReadFile("/sys/class/dmi/id/board_name")
	if err == nil {
		return strings.TrimSpace(string(data))
	}
	return "Unknown"
}

func getGPU() []string {
	path, err := exec.LookPath("lspci")
	if err != nil {
		return nil
	}
	out, err := exec.Command(path, "-mm").Output()
	if err != nil {
		return nil
	}
	
	lines := strings.Split(string(out), "\n")
	var gpus []string
	for _, line := range lines {
		if line == "" { continue }
		lower := strings.ToLower(line)
		if strings.Contains(lower, "vga") || strings.Contains(lower, "3d controller") || strings.Contains(lower, "display controller") {
			parts := strings.Split(line, "\"")
			if len(parts) >= 6 {
				vendor := parts[3]
				device := parts[5]
				gpuName := fmt.Sprintf("%s %s", vendor, device)
				gpus = append(gpus, gpuName)
			}
		}
	}
	return gpus
}

func getResolution() string {
	path, err := exec.LookPath("xrandr")
	if err != nil {
		return ""
	}
	out, err := exec.Command(path).Output()
	if err != nil {
		return ""
	}
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.Contains(line, "*") {
			fields := strings.Fields(line)
			if len(fields) > 0 {
				return fields[0] 
			}
		}
	}
	return ""
}

func getDE() string {
	de := os.Getenv("XDG_CURRENT_DESKTOP")
	if de == "" {
		de = os.Getenv("DESKTOP_SESSION")
	}
	return de
}

func getWM() string {
	return "Unknown" 
}

func getPackages() string {
	var counts []string

	// Pacman
	if _, err := exec.LookPath("pacman"); err == nil {
		out, _ := exec.Command("pacman", "-Qq").Output()
		if len(out) > 0 {
			counts = append(counts, fmt.Sprintf("%d (pacman)", strings.Count(string(out), "\n")))
		}
	}
	// Dpkg
	if _, err := exec.LookPath("dpkg"); err == nil {
		out, _ := exec.Command("dpkg-query", "-f", "${binary:Package}\n", "-W").Output()
		if len(out) > 0 {
			counts = append(counts, fmt.Sprintf("%d (dpkg)", strings.Count(string(out), "\n")))
		}
	}
	// Rpm
	if _, err := exec.LookPath("rpm"); err == nil {
		out, _ := exec.Command("rpm", "-qa").Output()
		if len(out) > 0 {
			counts = append(counts, fmt.Sprintf("%d (rpm)", strings.Count(string(out), "\n")))
		}
	}
	// Snap
	if _, err := exec.LookPath("snap"); err == nil {
		out, _ := exec.Command("snap", "list").Output()
		if len(out) > 0 {
			lines := strings.Count(string(out), "\n")
			if lines > 1 { // header row
				counts = append(counts, fmt.Sprintf("%d (snap)", lines-1))
			}
		}
	}
	// Flatpak
	if _, err := exec.LookPath("flatpak"); err == nil {
		out, _ := exec.Command("flatpak", "list", "--app").Output()
		if len(out) > 0 {
			counts = append(counts, fmt.Sprintf("%d (flatpak)", strings.Count(string(out), "\n")))
		}
	}

	return strings.Join(counts, ", ")
}

func getTerminal() string {
    // 1. Check Standard Env Vars
    if tp := os.Getenv("TERM_PROGRAM"); tp != "" { return tp }
    
    // 2. Check Specific Env Vars
    if os.Getenv("KITTY_PID") != "" { return "kitty" }
    if os.Getenv("GNOME_TERMINAL_SCREEN") != "" || os.Getenv("GNOME_TERMINAL_SERVICE") != "" { return "gnome-terminal" }
    if os.Getenv("ALACRITTY_SOCKET") != "" || os.Getenv("ALACRITTY_LOG") != "" { return "alacritty" }

    // 3. Walk Process Tree: Shell -> Terminal
    // Parent of pulsefetch is Shell. Parent of Shell is Terminal.
    shellPID := os.Getppid()
    termPID, err := getPPID(shellPID)
    if err == nil {
        name, err := getProcessName(termPID)
        if err == nil {
             // Handle some names
             name = strings.TrimSuffix(name, "-") // sometimes gnome-terminal-
             if strings.Contains(name, "gnome-terminal") { return "gnome-terminal" }
             if strings.Contains(name, "alacritty") { return "alacritty" }
             if strings.Contains(name, "kitty") { return "kitty" }
             if strings.Contains(name, "termite") { return "termite" }
             if strings.Contains(name, "urxvt") { return "urxvt" }
             return name
        }
    }

    // 4. Fallback
    return os.Getenv("TERM")
}

func getPPID(pid int) (int, error) {
    data, err := os.ReadFile(fmt.Sprintf("/proc/%d/stat", pid))
    if err != nil { return 0, err }
    // Format: pid (comm) state ppid ...
    // comm can contain spaces and parentheses, so strictly finding the last ) is safer
    s := string(data)
    lastParen := strings.LastIndex(s, ")")
    if lastParen == -1 { return 0, fmt.Errorf("bad stat format") }
    
    rest := s[lastParen+1:]
    fields := strings.Fields(rest)
    if len(fields) < 2 { return 0, fmt.Errorf("bad stat format") }
    // fields[0] is state, fields[1] is ppid
    return strconv.Atoi(fields[1])
}

func getProcessName(pid int) (string, error) {
    data, err := os.ReadFile(fmt.Sprintf("/proc/%d/comm", pid))
    if err != nil { return "", err }
    return strings.TrimSpace(string(data)), nil
}

func getTheme() string {
	home, err := os.UserHomeDir()
	if err != nil { return "" }
	path := fmt.Sprintf("%s/.config/gtk-3.0/settings.ini", home)
	return parseGtkSetting(path, "gtk-theme-name")
}

func getIcons() string {
	home, err := os.UserHomeDir()
	if err != nil { return "" }
	path := fmt.Sprintf("%s/.config/gtk-3.0/settings.ini", home)
	return parseGtkSetting(path, "gtk-icon-theme-name")
}

func parseGtkSetting(path, key string) string {
	data, err := os.ReadFile(path)
	if err != nil { return "" }
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, key) {
			parts := strings.Split(line, "=")
			if len(parts) == 2 {
				return strings.TrimSpace(parts[1])
			}
		}
	}
	return ""
}
