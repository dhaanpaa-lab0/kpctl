package platform

import (
	"log"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
)

var appName = "kpctl"

func GetConfigFolder() string {
	configPath := filepath.Join(xdg.ConfigHome, appName)
	if errMakeConfigFolder := os.MkdirAll(configPath, 0755); errMakeConfigFolder != nil {
		log.Fatalf("Failed to create config folder: %v", errMakeConfigFolder)
	}
	return configPath
}

func GetConfigFolderFileName(fileName string) string {
	return filepath.Join(GetConfigFolder(), fileName)
}
