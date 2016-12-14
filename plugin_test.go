package main

import (
	"io/ioutil"
	"os"
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
	expected := "upgrade --install test-release ./chart/test --set image.tag=v.0.1.0,nameOverride=my-over-app --dry-run --debug"
	if res != expected {
		t.Errorf("Result is %s and we expected %s", res, expected)
	}
}

func TestResolveSecrets(t *testing.T) {
	secrets := make([]string, 2)
	secrets[0] = "TAG"
	secrets[1] = "API_SERVER"
	tag := "v0.1.1"
	api := "http://apiserver"
	os.Setenv("TAG", tag)
	os.Setenv("API_SERVER", api)

	plugin := &Plugin{
		Config: Config{
			APIServer:     "${API_SERVER}",
			Token:         "secret-token",
			HelmCommand:   nil,
			Namespace:     "default",
			SkipTLSVerify: true,
			Debug:         true,
			DryRun:        true,
			Chart:         "./chart/test",
			Release:       "test-release",
			Values:        "image.tag=$TAG,api=${API_SERVER},nameOverride=my-over-app,second.tag=${TAG}",
			Secrets:       secrets,
		},
	}

	resolveSecrets(plugin)
	// test that the subsitution works
	if !strings.Contains(plugin.Config.Values, tag) {
		t.Errorf("env var %s not resolved %s", secrets[0], tag)
	}
	// test that subistutes more than 1 envvar
	if strings.Contains(plugin.Config.Values, secrets[0]) {
		t.Errorf("env var %s not resolved %s", secrets[0], tag)
	}
	// // test that the subsitution works with more than one envvar
	if strings.Contains(plugin.Config.Values, secrets[1]) {
		t.Errorf("env var %s not resolved %s", secrets[1], api)
	}
	// // test that the subsitution works with more than one envvar
	if strings.Contains(plugin.Config.APIServer, secrets[1]) {
		t.Errorf("env var %s not resolved %s", secrets[1], api)
	}
}
