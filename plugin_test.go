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
	os.Setenv("DRONE_BUILD_EVENT", "push")
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
			Wait:          true,
			ReuseValues:   true,
		},
	}
	setHelmCommand(plugin)
	res := strings.Join(plugin.Config.HelmCommand[:], " ")
	expected := "upgrade --install test-release ./chart/test --set image.tag=v.0.1.0,nameOverride=my-over-app --namespace default --dry-run --debug --wait --reuse-values"
	if res != expected {
		t.Errorf("Result is %s and we expected %s", res, expected)
	}
}

func TestResolveSecrets(t *testing.T) {
	tag := "v0.1.1"
	api := "http://apiserver"
	os.Setenv("MY_TAG", tag)
	os.Setenv("MY_API_SERVER", api)
	os.Setenv("MY_TOKEN", "12345")

	plugin := &Plugin{
		Config: Config{
			HelmCommand:   nil,
			Namespace:     "default",
			SkipTLSVerify: true,
			Debug:         true,
			DryRun:        true,
			Chart:         "./chart/test",
			Release:       "test-release",
			Prefix:        "MY",
			Values:        "image.tag=$TAG,api=${API_SERVER},nameOverride=my-over-app,second.tag=${TAG}",
		},
	}

	resolveSecrets(plugin)
	// test that the subsitution works
	if !strings.Contains(plugin.Config.Values, tag) {
		t.Errorf("env var ${TAG} not resolved %s", tag)
	}
	if strings.Contains(plugin.Config.Values, "${TAG}") {
		t.Errorf("env var ${TAG} not resolved %s", tag)
	}

	if plugin.Config.APIServer != api {
		t.Errorf("env var ${API_SERVER} not resolved %s", api)
	}
}

func TestGetEnvVars(t *testing.T) {

	testText := "this should be ${TAG} now"
	result := getEnvVars(testText)
	if len(result) == 0 {
		t.Error("No envvar was found")
	}
	envvar := result[0]
	if !strings.Contains(envvar[2], "TAG") {
		t.Errorf("envvar not found in %s", testText)
	}
}

func TestReplaceEnvvars(t *testing.T) {
	tag := "tagged"
	os.Setenv("MY_TAG", tag)
	prefix := "MY"
	testText := "this should be ${TAG} now ${TAG}"
	result := getEnvVars(testText)
	resolved := replaceEnvvars(result, prefix, testText)
	if !strings.Contains(resolved, tag) {
		t.Errorf("EnvVar MY_TAG no replaced by %s  -- %s \n", tag, resolved)
	}
}

func TestSetHelmHelp(t *testing.T) {
	plugin := &Plugin{
		Config: Config{
			HelmCommand:   nil,
			Namespace:     "default",
			SkipTLSVerify: true,
			Debug:         true,
			DryRun:        true,
			Chart:         "./chart/test",
			Release:       "test-release",
			Prefix:        "MY",
			Values:        "image.tag=$TAG,api=${API_SERVER},nameOverride=my-over-app,second.tag=${TAG}",
		},
	}
	setHelmHelp(plugin)
	if plugin.Config.HelmCommand == nil {
		t.Error("Helm help is not displayed")
	}
}

func TestDetHelmInit(t *testing.T) {
	plugin := &Plugin{
		Config: Config{
			HelmCommand:   nil,
			Namespace:     "default",
			SkipTLSVerify: true,
			Debug:         true,
			DryRun:        true,
			Chart:         "./chart/test",
			Release:       "test-release",
			Prefix:        "MY",
			Values:        "image.tag=$TAG,api=${API_SERVER},nameOverride=my-over-app,second.tag=${TAG}",
			TillerNs:      "system-test",
		},
	}
	init := doHelmInit(plugin)
	result := strings.Join(init, " ")
	expected := "init --tiller-namespace " + plugin.Config.TillerNs

	if expected != result {
		t.Error("Tiller not installed in proper namespace")
	}
}

func TestDetHelmInitClient(t *testing.T) {
	plugin := &Plugin{
		Config: Config{
			HelmCommand:   nil,
			Namespace:     "default",
			SkipTLSVerify: true,
			Debug:         true,
			DryRun:        true,
			Chart:         "./chart/test",
			Release:       "test-release",
			Prefix:        "MY",
			Values:        "image.tag=$TAG,api=${API_SERVER},nameOverride=my-over-app,second.tag=${TAG}",
			ClientOnly:    true,
		},
	}
	init := doHelmInit(plugin)
	result := strings.Join(init, " ")
	expected := "init "
	if plugin.Config.ClientOnly {
		expected = expected + "--client-only"
	}

	if expected != result {
		t.Error("Helm cannot init in client only")
	}
}

func TestDetHelmInitUpgrade(t *testing.T) {
	plugin := &Plugin{
		Config: Config{
			HelmCommand:   nil,
			Namespace:     "default",
			SkipTLSVerify: true,
			Debug:         true,
			DryRun:        true,
			Chart:         "./chart/test",
			Release:       "test-release",
			Prefix:        "MY",
			Values:        "image.tag=$TAG,api=${API_SERVER},nameOverride=my-over-app,second.tag=${TAG}",
			Upgrade:       true,
		},
	}
	init := doHelmInit(plugin)
	result := strings.Join(init, " ")
	expected := "init "
	if plugin.Config.Upgrade {
		expected = expected + "--upgrade"
	}

	if expected != result {
		t.Error("Helm cannot init in client only")
	}
}
func TestResolveSecretsFallback(t *testing.T) {
	tag := "v0.1.1"
	api := "http://apiserver"
	os.Setenv("MY_TAG", tag)
	os.Setenv("MY_API_SERVER", api)
	os.Setenv("MY_TOKEN", "12345")
	os.Setenv("NOTTOKEN", "99999")

	plugin := &Plugin{
		Config: Config{
			HelmCommand:   nil,
			Namespace:     "default",
			SkipTLSVerify: true,
			Debug:         true,
			DryRun:        true,
			Chart:         "./chart/test",
			Release:       "test-release",
			Prefix:        "MY",
			Values:        "image.tag=$TAG,api=${API_SERVER},nottoken=${NOTTOKEN},nameOverride=my-over-app,second.tag=${TAG}",
		},
	}

	resolveSecrets(plugin)
	// test that the subsitution works
	if !strings.Contains(plugin.Config.Values, tag) {
		t.Errorf("env var ${TAG} not resolved %s", tag)
	}
	if strings.Contains(plugin.Config.Values, "${TAG}") {
		t.Errorf("env var ${TAG} not resolved %s", tag)
	}

	if plugin.Config.APIServer != api {
		t.Errorf("env var ${API_SERVER} not resolved %s", api)
	}
	if !strings.Contains(plugin.Config.Values, "99999") {
		t.Errorf("envar ${NOTTOKEN} has not been resolved to 99999, not using prefix")
	}
}
