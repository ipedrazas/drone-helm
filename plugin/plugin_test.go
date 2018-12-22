package plugin

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestInitialiseKubeconfig(t *testing.T) {

	plugin := Plugin{
		Config: Config{
			APIServer:      "http://myapiserver",
			Token:          "secret-token",
			Certificate:    "my-cert-data",
			ServiceAccount: "default-account",
			HelmCommand:    "",
			Namespace:      "default",
			SkipTLSVerify:  false, // if set the true with Certificate, this test will fail
		},
	}

	configfile := "config3.test"
	initialiseKubeconfig(&plugin.Config, "../kubeconfig", configfile)
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
	if !strings.Contains(kubeConfigStr, "default-account") {
		t.Errorf("Kubeconfig doesn't render serviceaccount")
	}
	if !strings.Contains(kubeConfigStr, "my-cert-data") {
		t.Errorf("Kubeconfig doesn't render certificate")
	}
	if strings.Contains(kubeConfigStr, "aws-iam-authenticator") {
		t.Errorf("Kubeconfig renders EKS cluster configuration")
	}
}

func TestInitialiseKubeconfigEKS(t *testing.T) {

	plugin := Plugin{
		Config: Config{
			APIServer:     "http://myapiserver",
			Certificate:   "my-cert-data",
			EKSCluster:    "my-eks-cluster-name",
			EKSRoleARN:    "my-eks-role-arn",
			HelmCommand:   "",
			Namespace:     "default",
			SkipTLSVerify: false, // if set the true with Certificate, this test will fail
		},
	}

	configfile := "config3.test"
	initialiseKubeconfig(&plugin.Config, "../kubeconfig", configfile)
	data, err := ioutil.ReadFile(configfile)
	if err != nil {
		t.Errorf("Error reading file %v", err)
	}
	kubeConfigStr := string(data)

	if strings.Contains(kubeConfigStr, "token:") {
		t.Errorf("Kubeconfig renders token")
	}
	if !strings.Contains(kubeConfigStr, "http://myapiserver") {
		t.Errorf("Kubeconfig doesn't render APIServer")
	}
	if !strings.Contains(kubeConfigStr, "my-cert-data") {
		t.Errorf("Kubeconfig doesn't render certificate")
	}
	if !strings.Contains(kubeConfigStr, "my-eks-cluster-name") {
		t.Errorf("Kubeconfig doesn't render EKS cluster name")
	}
	if !strings.Contains(kubeConfigStr, "my-eks-role-arn") {
		t.Errorf("Kubeconfig doesn't render EKS role ARN")
	}
}

func TestGetHelmCommandEmptyPushEvent(t *testing.T) {
	os.Setenv("DRONE_BUILD_EVENT", "push")
	plugin := &Plugin{
		Config: Config{
			APIServer:     "http://myapiserver",
			Token:         "secret-token",
			HelmCommand:   "",
			Namespace:     "default",
			SkipTLSVerify: true,
			Debug:         true,
			DryRun:        true,
			Chart:         "./chart/test",
			Version:       "1.2.3",
			Release:       "test-release",
			Values:        `"image.tag=v.0.1.0,nameOverride=my-over-app"`,
			StringValues:  `"long_string_value=1234567890"`,
			Wait:          true,
			ReuseValues:   true,
			Timeout:       "500",
			Force:         true,
		},
	}
	setHelmCommand(plugin)
	res := strings.Join(plugin.command[:], " ")
	expected := "upgrade --install test-release ./chart/test --version 1.2.3 --set image.tag=v.0.1.0,nameOverride=my-over-app --set-string long_string_value=1234567890 --namespace default --dry-run --debug --wait --reuse-values --timeout 500 --force"
	if res != expected {
		t.Errorf("Result is %s and we expected %s", res, expected)
	}
}

func TestGetHelmCommandUpgrade(t *testing.T) {
	os.Setenv("DRONE_BUILD_EVENT", "push")
	plugin := &Plugin{
		Config: Config{
			APIServer:     "http://myapiserver",
			Token:         "secret-token",
			HelmCommand:   "upgrade",
			Namespace:     "default",
			SkipTLSVerify: true,
			Debug:         true,
			DryRun:        true,
			Chart:         "./chart/test",
			Version:       "1.2.3",
			Release:       "test-release",
			Values:        `"image.tag=v.0.1.0,nameOverride=my-over-app"`,
			StringValues:  `"long_string_value=1234567890"`,
			Wait:          true,
			ReuseValues:   true,
			Timeout:       "500",
			Force:         true,
		},
	}
	setHelmCommand(plugin)
	res := strings.Join(plugin.command[:], " ")
	expected := "upgrade --install test-release ./chart/test --version 1.2.3 --set image.tag=v.0.1.0,nameOverride=my-over-app --set-string long_string_value=1234567890 --namespace default --dry-run --debug --wait --reuse-values --timeout 500 --force"
	if res != expected {
		t.Errorf("Result is %s and we expected %s", res, expected)
	}
}

func TestGetHelmDeleteCommand(t *testing.T) {
	os.Setenv("DRONE_BUILD_EVENT", "delete")
	plugin := &Plugin{
		Config: Config{
			APIServer:     "http://myapiserver",
			Token:         "secret-token",
			HelmCommand:   "delete",
			Namespace:     "default",
			SkipTLSVerify: true,
			Debug:         true,
			Chart:         "./chart/test",
			Release:       "test-release",
			Values:        "image.tag=v.0.1.0,nameOverride=my-over-app",
			StringValues:  "long_string_value=1234567890",
			Wait:          true,
		},
	}
	setHelmCommand(plugin)
	res := strings.Join(plugin.command[:], " ")
	expected := "delete test-release"
	if res != expected {
		t.Errorf("Result is %s and we expected %s", res, expected)
	}
}

func TestGetHelmCommandLint(t *testing.T) {
	os.Setenv("DRONE_BUILD_EVENT", "push")
	plugin := &Plugin{
		Config: Config{
			APIServer:     "http://myapiserver",
			Token:         "secret-token",
			HelmCommand:   "lint",
			Namespace:     "default",
			SkipTLSVerify: true,
			Chart:         "./chart/test",
			Values:        `"image.tag=v.0.1.0,nameOverride=my-over-app"`,
			StringValues:  `"long_string_value=1234567890"`,
		},
	}
	setHelmCommand(plugin)
	res := strings.Join(plugin.command[:], " ")
	expected := "lint ./chart/test --set image.tag=v.0.1.0,nameOverride=my-over-app --set-string long_string_value=1234567890 --namespace default"
	if res != expected {
		t.Errorf("Result is %s and we expected %s", res, expected)
	}
}



func TestGetHelmDeleteCommandOverried(t *testing.T) {
	os.Setenv("DRONE_BUILD_EVENT", "deployment")
	plugin := &Plugin{
		Config: Config{
			APIServer:     "http://myapiserver",
			Token:         "secret-token",
			HelmCommand:   "delete",
			Namespace:     "default",
			TillerNs:      "default-tiller-ns",
			SkipTLSVerify: true,
			Debug:         true,
			DryRun:        true,
			Chart:         "./chart/test",
			Release:       "test-release",
			Values:        "image.tag=v.0.1.0,nameOverride=my-over-app",
			StringValues:  "long_string_value=1234567890",
			Wait:          true,
			Purge:         true,
		},
	}
	setHelmCommand(plugin)
	res := strings.Join(plugin.command[:], " ")
	expected := "delete test-release --tiller-namespace default-tiller-ns --dry-run --purge"
	if res != expected {
		t.Errorf("Result is %s and we expected %s", res, expected)
	}
}

func TestResolveSecrets(t *testing.T) {

	// Test resolving secrets from env
	testEnvs := []struct {
		prefix  string
		tag     string
		api     string
		token   string
		account string
	}{
		{prefix: "PROD", tag: "v0.1.1", api: "http://apiserver", token: "12345", account: "helm"},
		{prefix: "STAGING", tag: "12345678", api: "http://apiserver", token: "12345", account: "helm"},
	}
	for _, env := range testEnvs {
		envMap := map[string]string{
			"TAG":              env.tag,
			"API_SERVER":       env.api,
			"KUBERNETES_TOKEN": env.token,
			"SERVICE_ACCOUNT":  env.account,
		}

		for envKey, envValue := range envMap {
			os.Setenv(fmt.Sprintf("%s_%s", env.prefix, envKey), envValue)
		}

		plugin := &Plugin{
			Config: Config{
				HelmCommand:   "",
				Namespace:     "default",
				SkipTLSVerify: true,
				Debug:         true,
				DryRun:        true,
				Chart:         "./chart/test",
				Release:       "test-release",
				Prefix:        env.prefix,
				Values:        "image.tag=$TAG,api=${API_SERVER},nameOverride=my-over-app,second.tag=${TAG}",
				StringValues:  "long_string_value=1234567890",
			},
		}

		resolveSecrets(plugin)
		// test that the subsitution works
		fmt.Println(plugin.Config.Values)
		if !strings.Contains(plugin.Config.Values, env.tag) {
			t.Errorf("env var ${TAG} not resolved %s", env.tag)
		}
		if strings.Contains(plugin.Config.Values, "${TAG}") {
			t.Errorf("env var ${TAG} not resolved %s", env.tag)
		}

		if plugin.Config.APIServer != env.api {
			t.Errorf("env var ${API_SERVER} not resolved %s", env.api)
		}
		if plugin.Config.Token != env.token {
			t.Errorf("env var ${KUBERNETES_TOKEN} not resolved %s", env.token)
		}
		if plugin.Config.ServiceAccount != env.account {
			t.Errorf("env var ${SERVICE_ACCOUNT} not resolved %s", env.account)
		}

		// clean up
		for envKey := range envMap {
			os.Unsetenv(fmt.Sprintf("%s_%s", env.prefix, envKey))
		}
	}

	// Test resolving provided values
	testInput := []struct {
		server       string
		values       string
		stringValues string
		token        string
		account      string
	}{
		{server: "http://apiserver2", token: "123456", account: "helm2", values: "aval=test", stringValues: "long_string_value=1234567890"},
	}
	for _, input := range testInput {
		plugin := &Plugin{
			Config: Config{
				APIServer:      input.server,
				ServiceAccount: input.account,
				Token:          input.token,
				Values:         input.values,
				StringValues:   input.stringValues,
			},
		}

		resolveSecrets(plugin)
		if plugin.Config.APIServer != input.server {
			t.Errorf("failed to keep APIServer '%s' got '%s'", input.server, plugin.Config.APIServer)
		}
		if plugin.Config.ServiceAccount != input.account {
			t.Errorf("failed to keep ServiceAccount '%s' got '%s'", input.account, plugin.Config.ServiceAccount)
		}
		if plugin.Config.Token != input.token {
			t.Errorf("failed to keep Token '%s' got '%s'", input.token, plugin.Config.Token)
		}
		if plugin.Config.Values != input.values {
			t.Errorf("failed to keep Values '%s' got '%s'", input.values, plugin.Config.Values)
		}
		if plugin.Config.StringValues != input.stringValues {
			t.Errorf("failed to keep StringValues '%s' got '%s'", input.stringValues, plugin.Config.StringValues)
		}
	}
}

func TestDetHelmRepoAdd(t *testing.T) {
	plugin := &Plugin{
		Config: Config{
			HelmCommand:   "",
			Namespace:     "default",
			SkipTLSVerify: true,
			Debug:         true,
			DryRun:        true,
			Chart:         "./chart/test",
			Release:       "test-release",
			Prefix:        "MY",
			Values:        "image.tag=$TAG,api=${API_SERVER},nameOverride=my-over-app,second.tag=${TAG}",
			StringValues:  "long_string_value=1234567890",
			ClientOnly:    true,
			HelmRepos: []string{
				`"r1=http://r1.example.com"`, //handle quoted strings
				`r2=http://r2.example.com`,   //and unquoted strings
			},
		},
	}
	expected := []string{
		"repo add r1 http://r1.example.com",
		"repo add r2 http://r2.example.com",
	}

	for i, r := range plugin.Config.HelmRepos {
		repos, err := doHelmRepoAdd(r)
		if err != nil {
			t.Error(err)
		}
		result := strings.Join(repos, " ")
		if expected[i] != result {
			t.Errorf("Helm cannot add remote repositories - expected %q - got %q",
				expected[i],
				result,
			)
		}
	}
}

func TestDependencyUpdate(t *testing.T) {
	chart := "./chart"
	expected := "dependency update ./chart"
	result := strings.Join(doDependencyUpdate(chart), " ")

	if expected != result {
		t.Errorf("Helm cannot update repositories - expected %q - got %q", expected, result)
	}
}

func TestHelmAddRepositoryError(t *testing.T) {
	_, err := doHelmRepoAdd("drone-helm=bad://drone-helm.example.com:443/stable")
	if err == nil {
		t.Errorf("Expect to see error when repo URL is invalid")
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
	resolved := replaceEnvvars(result, prefix, testText, false)
	if !strings.Contains(resolved, tag) {
		t.Errorf("EnvVar MY_TAG no replaced by %s  -- %s \n", tag, resolved)
	}
}

func TestSetHelpCommand(t *testing.T) {
	plugin := &Plugin{
		Config: Config{
			HelmCommand:   "",
			Namespace:     "default",
			SkipTLSVerify: true,
			Debug:         true,
			DryRun:        true,
			Chart:         "./chart/test",
			Release:       "test-release",
			Prefix:        "MY",
			Values:        "image.tag=$TAG,api=${API_SERVER},nameOverride=my-over-app,second.tag=${TAG}",
			StringValues:  "long_string_value=1234567890",
		},
	}
	setHelpCommand(plugin)
	if plugin.command == nil {
		t.Error("Helm help is not displayed")
	}
}

func TestDetHelmInit(t *testing.T) {
	plugin := &Plugin{
		Config: Config{
			HelmCommand:   "",
			Namespace:     "default",
			SkipTLSVerify: true,
			Debug:         true,
			DryRun:        true,
			Chart:         "./chart/test",
			Release:       "test-release",
			Prefix:        "MY",
			Values:        "image.tag=$TAG,api=${API_SERVER},nameOverride=my-over-app,second.tag=${TAG}",
			StringValues:  "long_string_value=1234567890",
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
			HelmCommand:   "",
			Namespace:     "default",
			SkipTLSVerify: true,
			Debug:         true,
			DryRun:        true,
			Chart:         "./chart/test",
			Release:       "test-release",
			Prefix:        "MY",
			Values:        "image.tag=$TAG,api=${API_SERVER},nameOverride=my-over-app,second.tag=${TAG}",
			StringValues:  "long_string_value=1234567890",
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
			HelmCommand:   "",
			Namespace:     "default",
			SkipTLSVerify: true,
			Debug:         true,
			DryRun:        true,
			Chart:         "./chart/test",
			Release:       "test-release",
			Prefix:        "MY",
			Values:        "image.tag=$TAG,api=${API_SERVER},nameOverride=my-over-app,second.tag=${TAG}",
			StringValues:  "long_string_value=1234567890",
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

func TestDetHelmInitCanary(t *testing.T) {
	plugin := &Plugin{
		Config: Config{
			HelmCommand:   "",
			Namespace:     "default",
			SkipTLSVerify: true,
			Debug:         true,
			DryRun:        true,
			Chart:         "./chart/test",
			Release:       "test-release",
			Prefix:        "MY",
			Values:        "image.tag=$TAG,api=${API_SERVER},nameOverride=my-over-app,second.tag=${TAG}",
			StringValues:  "long_string_value=1234567890",
			CanaryImage:   true,
		},
	}
	init := doHelmInit(plugin)
	result := strings.Join(init, " ")
	expected := "init "
	if plugin.Config.CanaryImage {
		expected = expected + "--canary-image"
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
			HelmCommand:   "",
			Namespace:     "default",
			SkipTLSVerify: true,
			Debug:         true,
			DryRun:        true,
			Chart:         "./chart/test",
			Release:       "test-release",
			Prefix:        "MY",
			Values:        "image.tag=$TAG,api=${API_SERVER},nottoken=${NOTTOKEN},nameOverride=my-over-app,second.tag=${TAG}",
			StringValues:  "long_string_value=1234567890",
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
func TestHelmInitByAliyunStableRepo(t *testing.T) {
	plugin := &Plugin{
		Config: Config{
			HelmCommand:   "",
			Namespace:     "default",
			SkipTLSVerify: true,
			Debug:         true,
			DryRun:        true,
			Chart:         "./chart/test",
			Release:       "test-release",
			Prefix:        "MY",
			Values:        "image.tag=$TAG,api=${API_SERVER},nameOverride=my-over-app,second.tag=${TAG}",
			StringValues:  "long_string_value=1234567890",
			StableRepoURL: "https://kubernetes.oss-cn-hangzhou.aliyuncs.com/charts",
		},
	}
	init := doHelmInit(plugin)
	result := strings.Join(init, " ")
	expected := "init --stable-repo-url " + plugin.Config.StableRepoURL
	if expected != result {
		t.Error("Helm cannot init for stable repository")
	}
}
