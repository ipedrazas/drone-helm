package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/alecthomas/template"
)

var HELM_BIN = "/bin/helm"
var KUBECONFIG = "/root/.kube/kubeconfig"
var CONFIG = "/root/.kube/config"

type (
	// Config maps the params we need to run Helm
	Config struct {
		APIServer     string   `json:"api_server"`
		Token         string   `json:"token"`
		HelmCommand   []string `json:"helm_command"`
		SkipTLSVerify bool     `json:"tls_skip_verify"`
		Namespace     string   `json:"namespace"`
		Release       string   `json:"release"`
		Chart         string   `json:"chart"`
		Values        string   `json:"values"`
	}
	// Plugin default
	Plugin struct {
		Config Config
	}
)

// Exec default method
func (p Plugin) Exec() error {
	if p.Config.APIServer == "" {
		return fmt.Errorf("Error: API Server is needed to deploy.")
	}
	if p.Config.Token == "" {
		return fmt.Errorf("Error: Token is needed to deploy.")
	}
	initialiseKubeconfig(&p.Config, KUBECONFIG, CONFIG)
	fmt.Println(p)
	init := make([]string, 1)
	init[0] = "init"
	err := runCommand(init)
	if err != nil {
		return fmt.Errorf("Error running helm comand: " + strings.Join(init[:], " "))
	}
	upgrade := make([]string, 7)
	upgrade[0] = "upgrade"
	upgrade[1] = "--install"
	upgrade[2] = p.Config.Release
	upgrade[3] = p.Config.Chart
	upgrade[4] = "--debug"
	if p.Config.Values != "" {
		upgrade[5] = "--set"
		upgrade[6] = p.Config.Values
	}

	err = runCommand(upgrade)
	if err != nil {
		return fmt.Errorf("Error running helm comand: " + strings.Join(upgrade[:], " "))
	}
	return nil
}

func initialiseKubeconfig(params *Config, source string, target string) error {
	t, _ := template.ParseFiles(source)
	f, err := os.Create(target)
	err = t.Execute(f, params)
	f.Close()
	return err
}

func runCommand(params []string) error {
	cmd := new(exec.Cmd)
	cmd = exec.Command(HELM_BIN, params...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	return err
}
