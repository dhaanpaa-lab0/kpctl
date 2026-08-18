package cmd

import (
	"bytes"
	"testing"

	"github.com/pterm/pterm"
	"github.com/spf13/viper"
	"nexus-sds.com/kpctl/pkg/sysConfig"
)

func TestListGlobalsCmd(t *testing.T) {
	viper.Reset()
	sysConfig.AddGlobalsConfigMapValue("ENV_NAME", "production")
	sysConfig.AddGlobalsConfigMapValue("CLUSTER", "main-cluster")

	var buf bytes.Buffer
	pterm.SetDefaultOutput(&buf)

	rootCmd.SetArgs([]string{"listGlobals"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("rootCmd.Execute(listGlobals) failed: %v", err)
	}

	output := buf.String()
	if !bytes.Contains(buf.Bytes(), []byte("Key")) || !bytes.Contains(buf.Bytes(), []byte("Value")) {
		t.Errorf("expected table header with Key and Value in output, got: %s", output)
	}
	if !bytes.Contains(buf.Bytes(), []byte("ENV_NAME")) || !bytes.Contains(buf.Bytes(), []byte("production")) {
		t.Errorf("expected ENV_NAME and production in output, got: %s", output)
	}
	if !bytes.Contains(buf.Bytes(), []byte("CLUSTER")) || !bytes.Contains(buf.Bytes(), []byte("main-cluster")) {
		t.Errorf("expected CLUSTER and main-cluster in output, got: %s", output)
	}
}
