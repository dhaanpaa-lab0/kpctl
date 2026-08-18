package sysConfig

import (
	"fmt"

	"github.com/spf13/viper"
	"nexus-sds.com/kpctl/pkg/helpers"
)

func AddGlobalsConfigMapValue(key string, value string) {
	globals := viper.GetStringMapString("globals")
	globals[key] = value
	viper.Set("globals", globals)
}

func SetConfigValueFromString(key string, value string) {
	viper.Set(key, value)
}

func GetConfigValueFromMap(key string) string {
	globals := viper.GetStringMapString("globals")
	return globals[key]
}

func GetGlobalsConfigMap() map[string]string {
	return viper.GetStringMapString("globals")
}

func GetConfigValueFromString(key string) string {
	return viper.GetString(key)
}

func SetConfigValuePrompted(key string, prompt string, defaultValue string) {
	if defaultValue == "^" {
		defaultValue = GetConfigValueFromString(key)
	}
	v := helpers.GetInputAsString(prompt, defaultValue)
	viper.Set(key, v)
}

func SetConfigValuePromptedSelect(key string, prompt string, options []string) {
	selectedOption := helpers.SelectOptionAsString(prompt, options)
	viper.Set(key, selectedOption)
}

func DeleteGlobalConfigMapValue(key string) {
	globals := viper.GetStringMapString("globals")
	delete(globals, key)
	viper.Set("globals", globals)
}

func UpdateGlobalConfig() {
	errWritingConfig := viper.WriteConfig()
	if errWritingConfig != nil {
		fmt.Println("Error writing config:", errWritingConfig)
		return
	}
}
