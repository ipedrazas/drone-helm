package main

import (
	"io/ioutil"
	"strings"
	"testing"
)

func TestInitialiseKubeconfig(t *testing.T) {

	cmd := make([]string, 2)
	cmd[0] = "install"
	cmd[1] = "--debug"

	plugin := Plugin{
		Config: Config{
			APIServer:     "http://myapiserver",
			Token:         "secret-token",
			HelmCommand:   cmd,
			Namespace:     "default",
			SkipTLSVerify: true,
		},
	}

	configfile := "config3.test"
	initialiseKubeconfig(&plugin.Config, "kubeconfig", configfile)
	data, err := ioutil.ReadFile(configfile)
	if err != nil {
		t.Errorf("Error reading file %v", err)
	}
	kubeConfigStr := string(data)

	if !strings.Contains(kubeConfigStr, "secret-token") {
		t.Errorf("Kubeconfig doesn't render token")
	}
	if !strings.Contains(kubeConfigStr, "http://myapiserver") {
		t.Errorf("Kubeconfig doesn't render APIServer")
	}

}

func TestGetHelmCommand(t *testing.T) {
	plugin := &Plugin{
		Config: Config{
			APIServer:     "http://myapiserver",
			Token:         "secret-token",
			HelmCommand:   nil,
			Namespace:     "default",
			SkipTLSVerify: true,
			Debug:         true,
			DryRun:        true,
			Chart:         "./chart/test",
			Release:       "test-release",
			Values:        "image.tag=v.0.1.0,nameOverride=my-over-app",
		},
	}
	setHelmCommand(plugin)
	res := strings.Join(plugin.Config.HelmCommand[:], " ")
	expected := "upgrade --install test-release ./chart/test --debug --set image.tag=v.0.1.0,nameOverride=my-over-app --dry-run"
	if res != expected {
		t.Errorf("Result is %s and we expected %s", res, expected)
	}
}
