package sysConfig

import (
	"testing"

	"github.com/spf13/viper"
)

func TestGlobalsConfigMap(t *testing.T) {
	viper.Reset()

	AddGlobalsConfigMapValue("FOO", "bar")
	AddGlobalsConfigMapValue("BAZ", "qux")

	if val := GetConfigValueFromMap("FOO"); val != "bar" {
		t.Fatalf("expected FOO to be 'bar', got %q", val)
	}

	globals := GetGlobalsConfigMap()
	if len(globals) != 2 {
		t.Fatalf("expected 2 globals, got %d", len(globals))
	}
	if globals["FOO"] != "bar" || globals["BAZ"] != "qux" {
		t.Fatalf("unexpected globals map: %v", globals)
	}

	DeleteGlobalConfigMapValue("FOO")
	globalsAfterDelete := GetGlobalsConfigMap()
	if len(globalsAfterDelete) != 1 {
		t.Fatalf("expected 1 global after delete, got %d", len(globalsAfterDelete))
	}
	if _, exists := globalsAfterDelete["FOO"]; exists {
		t.Fatalf("FOO should have been deleted")
	}
}
