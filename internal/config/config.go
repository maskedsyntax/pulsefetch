package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	// Modules - keys match the TOML keys
	ShowOS           bool   `mapstructure:"show_os"`
	ShowHost         bool   `mapstructure:"show_host"`
	ShowKernel       bool   `mapstructure:"show_kernel"`
	ShowUptime       bool   `mapstructure:"show_uptime"`
	ShowPackages     bool   `mapstructure:"show_packages"`
	ShowShell        bool   `mapstructure:"show_shell"`
	ShowResolution   bool   `mapstructure:"show_resolution"`
	ShowDE           bool   `mapstructure:"show_de"`
	ShowWM           bool   `mapstructure:"show_wm"`
	ShowWMTheme      bool   `mapstructure:"show_wm_theme"`
	ShowTheme        bool   `mapstructure:"show_theme"`
	ShowIcons        bool   `mapstructure:"show_icons"`
	ShowTerminal     bool   `mapstructure:"show_terminal"`
	ShowCPU          bool   `mapstructure:"show_cpu"`
	ShowGPU          bool   `mapstructure:"show_gpu"`
	ShowMemory       bool   `mapstructure:"show_memory"`
	ShowDisk         bool   `mapstructure:"show_disk"`
	ShowNetwork      bool   `mapstructure:"show_network"`
	ShowBattery      bool   `mapstructure:"show_battery"`
	ShowSensors      bool   `mapstructure:"show_sensors"`
	
	// Usage
	ShowCPUUsage     bool   `mapstructure:"show_cpu_usage"`
	ShowMemoryUsage  bool   `mapstructure:"show_memory_usage"`
	ShowDiskUsage    bool   `mapstructure:"show_disk_usage"`
	ShowNetworkUsage bool   `mapstructure:"show_network_usage"`
	ShowBatteryUsage bool   `mapstructure:"show_battery_usage"`
	ShowSensorsUsage bool   `mapstructure:"show_sensors_usage"`

	// Image
	ImagePath string `mapstructure:"image_path"` // Path to custom image
	ImageMode string `mapstructure:"image_mode"` // "ascii", "none" (maybe "image" later)
}

func LoadConfig() (*Config, error) {
	viper.SetConfigType("toml") // Treat config as TOML regardless of extension if we force it

	// Defaults
	viper.SetDefault("show_os", true)
	viper.SetDefault("show_host", true)
	viper.SetDefault("show_kernel", true)
	viper.SetDefault("show_uptime", true)
	viper.SetDefault("show_packages", true)
	viper.SetDefault("show_shell", true)
	viper.SetDefault("show_resolution", true)
	viper.SetDefault("show_de", true)
	viper.SetDefault("show_wm", true)
	viper.SetDefault("show_terminal", true)
	viper.SetDefault("show_cpu", true)
	viper.SetDefault("show_gpu", true)
	viper.SetDefault("show_memory", true)
	viper.SetDefault("show_disk", true)
	viper.SetDefault("show_battery", true)
	viper.SetDefault("show_resolution", true)
	viper.SetDefault("image_mode", "ascii")

	// Check for specific config file: ~/.config/pulsefetch/pulsefetch.toml
	home, err := os.UserHomeDir()
	foundSpecific := false
	if err == nil {
		specificPath := filepath.Join(home, ".config", "pulsefetch", "pulsefetch.toml")
		if _, err := os.Stat(specificPath); err == nil {
			viper.SetConfigFile(specificPath)
			foundSpecific = true
		}
	}

	if !foundSpecific {
		viper.SetConfigName("pulsefetch")
		if err == nil {
			viper.AddConfigPath(filepath.Join(home, ".config", "pulsefetch"))
		}
		viper.AddConfigPath("/etc/pulsefetch")
		// Removed AddConfigPath(".") to avoid conflict with binary named 'pulsefetch'
	}

	err = viper.ReadInConfig()
	if err != nil {
		// If config file not found, we just return defaults
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// If we found the file but failed to parse, that's a real error
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
		// Config not found is fine, we use defaults
	}

	var cfg Config
	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	return &cfg, nil
}
