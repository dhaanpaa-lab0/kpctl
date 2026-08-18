package sysConfig

import (
	"errors"
	"fmt"
	"log"

	"github.com/spf13/viper"
	"nexus-sds.com/kpctl/pkg/platform"
)

func AddGlobalsConfigMapValue(key string, value string) {
	globals := viper.GetStringMapString("globals")
	globals[key] = value
	viper.Set("globals", globals)
}

func GetConfigValueFromMap(key string) string {
	globals := viper.GetStringMapString("globals")
	return globals[key]
}

func GetGlobalsConfigMap() map[string]string {
	return viper.GetStringMapString("globals")
}

func DeleteGlobalConfigMapValue(key string) {
	globals := viper.GetStringMapString("globals")
	delete(globals, key)
	viper.Set("globals", globals)
}

func UpdateGlobalConfig() {
	errWritingConfig := viper.WriteConfig()
	var configFileNotFoundError viper.ConfigFileNotFoundError

	if errWritingConfig != nil {
		if errors.As(errWritingConfig, &configFileNotFoundError) {
			configPath := platform.GetConfigFolderFileName("config.yaml")
			fmt.Printf("Config not found. Creating default at: %s\n", configPath)

			if writeErr := viper.WriteConfigAs(configPath); writeErr != nil {
				log.Fatalf("Failed to write initial config: %v", writeErr)
			}
		} else {
			// Real syntax or permission error found during loading
			log.Fatalf("Error writing config file: %v", errWritingConfig)
		}

		return
	}

	if errStrippingLegacyFields := stripLegacyProfileFields(viper.ConfigFileUsed()); errStrippingLegacyFields != nil {
		fmt.Println("Error cleaning up legacy config fields:", errStrippingLegacyFields)
	}
}
