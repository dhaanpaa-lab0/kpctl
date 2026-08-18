package sysConfig

import (
	"errors"
	"fmt"
	"os"
	"sort"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"nexus-sds.com/kpctl/pkg/helpers"
)

const (
	// DefaultProfileName is the profile used when none has been selected yet.
	DefaultProfileName = "default"

	profilesKey       = "profiles"
	currentProfileKey = "current_profile"
)

// legacyProfileFields lists the top-level config keys that existed before
// profiles were introduced.
var legacyProfileFields = []string{"k8s_context", "k8s_cluster_ns", "gitops_url", "gitops_folder"}

// Profile groups the per-environment settings configured via initConfig.
type Profile struct {
	K8sContext   string
	K8sClusterNs string
	GitopsURL    string
	GitopsFolder string
}

// MigrateLegacyProfileConfig moves top-level k8s_context, k8s_cluster_ns,
// gitops_url and gitops_folder settings (from a config file written before
// profiles existed) into a "default" profile. It only updates viper's
// in-memory state; the migrated layout is written to disk the next time a
// command persists the config (e.g. initConfig).
func MigrateLegacyProfileConfig() {
	if len(viper.GetStringMap(profilesKey)) > 0 {
		return
	}

	profile := map[string]string{}
	for _, field := range legacyProfileFields {
		if v := viper.GetString(field); v != "" {
			profile[field] = v
		}
	}
	if len(profile) == 0 {
		return
	}

	viper.Set(profilesKey+"."+DefaultProfileName, profile)
	if viper.GetString(currentProfileKey) == "" {
		viper.Set(currentProfileKey, DefaultProfileName)
	}
}

// stripLegacyProfileFields removes the pre-profile top-level keys from the
// config file on disk, if present. It is called after every config write so
// a migrated config file no longer carries the now-unused duplicate values.
// A missing or unparsable file is left alone.
func stripLegacyProfileFields(path string) error {
	if path == "" {
		return nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	root := map[string]any{}
	if err := yaml.Unmarshal(data, &root); err != nil {
		return err
	}

	removed := false
	for _, field := range legacyProfileFields {
		if _, ok := root[field]; ok {
			delete(root, field)
			removed = true
		}
	}
	if !removed {
		return nil
	}

	out, err := yaml.Marshal(root)
	if err != nil {
		return err
	}
	return os.WriteFile(path, out, 0644)
}

func profileFieldKey(profileName string, field string) string {
	return profilesKey + "." + profileName + "." + field
}

// GetActiveProfileName returns the name of the currently selected profile,
// falling back to DefaultProfileName if none has been set.
func GetActiveProfileName() string {
	if name := viper.GetString(currentProfileKey); name != "" {
		return name
	}
	return DefaultProfileName
}

// SetActiveProfileName selects the profile used by GetActiveProfile.
func SetActiveProfileName(name string) {
	viper.Set(currentProfileKey, name)
}

// ListProfileNames returns the names of all configured profiles, sorted.
func ListProfileNames() []string {
	m := viper.GetStringMap(profilesKey)
	names := make([]string, 0, len(m))
	for name := range m {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// ProfileExists reports whether a profile with the given name is configured.
func ProfileExists(name string) bool {
	_, ok := viper.GetStringMap(profilesKey)[name]
	return ok
}

// GetProfile returns the settings stored under the named profile.
func GetProfile(name string) Profile {
	return Profile{
		K8sContext:   viper.GetString(profileFieldKey(name, "k8s_context")),
		K8sClusterNs: viper.GetString(profileFieldKey(name, "k8s_cluster_ns")),
		GitopsURL:    viper.GetString(profileFieldKey(name, "gitops_url")),
		GitopsFolder: viper.GetString(profileFieldKey(name, "gitops_folder")),
	}
}

// GetActiveProfile returns the settings for the currently selected profile.
func GetActiveProfile() Profile {
	return GetProfile(GetActiveProfileName())
}

// SetProfileValue sets a single field on the named profile.
func SetProfileValue(profileName string, field string, value string) {
	viper.Set(profileFieldKey(profileName, field), value)
}

// SetActiveProfileValue sets a single field on the currently selected profile.
func SetActiveProfileValue(field string, value string) {
	SetProfileValue(GetActiveProfileName(), field, value)
}

// DeleteProfile removes the named profile. If it was the active profile, the
// active profile reverts to DefaultProfileName.
func DeleteProfile(name string) {
	profiles := viper.GetStringMap(profilesKey)
	delete(profiles, name)
	viper.Set(profilesKey, profiles)
	if GetActiveProfileName() == name {
		SetActiveProfileName(DefaultProfileName)
	}
}

// SetProfileValuePrompted prompts for a field's value, seeding the prompt
// with the field's current value when defaultValue is "^".
func SetProfileValuePrompted(profileName string, field string, prompt string, defaultValue string) {
	if defaultValue == "^" {
		defaultValue = viper.GetString(profileFieldKey(profileName, field))
	}
	v := helpers.GetInputAsString(prompt, defaultValue)
	SetProfileValue(profileName, field, v)
}

// SetProfileValuePromptedSelect prompts for a field's value from a list of options.
func SetProfileValuePromptedSelect(profileName string, field string, prompt string, options []string) {
	selected := helpers.SelectOptionAsString(prompt, options)
	SetProfileValue(profileName, field, selected)
}

// ProfileExport is the on-disk YAML representation of a profile, used by
// ExportProfile and ImportProfile.
type ProfileExport struct {
	Name         string `yaml:"name"`
	K8sContext   string `yaml:"k8s_context"`
	K8sClusterNs string `yaml:"k8s_cluster_ns"`
	GitopsURL    string `yaml:"gitops_url"`
	GitopsFolder string `yaml:"gitops_folder"`
}

// ExportProfile writes the named profile out to path as YAML.
func ExportProfile(name string, path string) error {
	if !ProfileExists(name) {
		return fmt.Errorf("profile %q does not exist", name)
	}

	p := GetProfile(name)
	export := ProfileExport{
		Name:         name,
		K8sContext:   p.K8sContext,
		K8sClusterNs: p.K8sClusterNs,
		GitopsURL:    p.GitopsURL,
		GitopsFolder: p.GitopsFolder,
	}

	data, err := yaml.Marshal(export)
	if err != nil {
		return fmt.Errorf("failed to marshal profile %q: %w", name, err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write profile export to %s: %w", path, err)
	}
	return nil
}

// ImportProfile reads a YAML file written by ExportProfile and creates or
// updates a profile from its contents. If name is non-empty it overrides the
// name embedded in the file. It returns the name of the imported profile.
func ImportProfile(path string, name string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read profile import file %s: %w", path, err)
	}

	var imported ProfileExport
	if err := yaml.Unmarshal(data, &imported); err != nil {
		return "", fmt.Errorf("failed to parse profile import file %s: %w", path, err)
	}

	profileName := name
	if profileName == "" {
		profileName = imported.Name
	}
	if profileName == "" {
		return "", errors.New("no profile name specified and none found in import file")
	}

	SetProfileValue(profileName, "k8s_context", imported.K8sContext)
	SetProfileValue(profileName, "k8s_cluster_ns", imported.K8sClusterNs)
	SetProfileValue(profileName, "gitops_url", imported.GitopsURL)
	SetProfileValue(profileName, "gitops_folder", imported.GitopsFolder)

	return profileName, nil
}
