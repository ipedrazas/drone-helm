package main

import "testing"

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

	initialiseKubeconfig(&plugin.Config, "kubeconfig", "config3.test")

}
