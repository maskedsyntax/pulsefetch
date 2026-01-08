package ui

import (
	"os"

	"pulsefetch/internal/config"
	"pulsefetch/internal/fetcher"
)

// Default "Pulse" Logo (Electric Shock / Pulse Wave)
const defaultLogo = `
      ____
     /    \
    |      |
   _|      |_
  / |      | \
 |  |      |  |
 |__|      |__|
    |      |
    |______|
`
// A better pulse wave ASCII
const pulseLogo = `
       _     
      | |    
 _____| |__  
|  _  | '_ \ 
| | | | | | |
| ||_||_|_|_|
|___|        
`

// Even better "Electric" Pulse
const electricLogo = `
       __
      /  /
     /  /
    /  /____
   /_______/
      /  /
     /__/
`

// Let's use a simple stylized bolt
const boltLogo = `
    _/|
   /_ |
    | |
    | |
   _| |_
  |_____|
`

func GetLogo(cfg *config.Config, info *fetcher.SystemInfo) string {
	if cfg.ImageMode == "none" {
		return ""
	}

	if cfg.ImagePath != "" {
		// Try to load image
		if _, err := os.Stat(cfg.ImagePath); err == nil {
			logo, _, _, err := RenderImage(cfg.ImagePath, info.Terminal)
			if err == nil {
				return logo
			}
		}
	}

	// Fallback to default ASCII
	return electricLogo
}