package main

import (
	"fmt"
	"log"
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
		Debug         bool     `json:"debug"`
		DryRun        bool     `json:"dry_run"`
		Secrets       []string `json:"secrets"`
		Prefix        string   `json:"prefix"`
	}
	// Plugin default
	Plugin struct {
		Config Config
	}
)

func setHelmCommand(p *Plugin) {
	upgrade := make([]string, 2)
	upgrade[0] = "upgrade"
	upgrade[1] = "--install"
	if p.Config.Release != "" {
		upgrade = append(upgrade, p.Config.Release)
	}
	upgrade = append(upgrade, p.Config.Chart)
	if p.Config.Values != "" {
		upgrade = append(upgrade, "--set")
		upgrade = append(upgrade, p.Config.Values)
	}
	if p.Config.DryRun {
		upgrade = append(upgrade, "--dry-run")
	}
	if p.Config.Debug {
		upgrade = append(upgrade, "--debug")
	}
	p.Config.HelmCommand = upgrade
}

// Exec default method
func (p *Plugin) Exec() error {
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
	setHelmCommand(p)
	if p.Config.Debug {
		log.Println("helm comand: " + strings.Join(p.Config.HelmCommand[:], " "))
	}
	err = runCommand(p.Config.HelmCommand)
	if err != nil {
		return fmt.Errorf("Error running helm comand: " + strings.Join(p.Config.HelmCommand[:], " "))
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

func resolveSecrets(p *Plugin) {
	if len(p.Config.Secrets) > 0 {
		for _, secret := range p.Config.Secrets {
			envval := os.Getenv(secret)
			p.Config.Values = resolveEnvVar(p.Config.Values, secret, envval)
			p.Config.APIServer = resolveEnvVar(p.Config.APIServer, secret, envval)
			p.Config.Token = resolveEnvVar(p.Config.Token, secret, envval)
		}
	}
}

// this functions checks if $VAR or ${VAR} exists and
// returns the text with resolved vars
func resolveEnvVar(key string, envvar string, envval string) string {
	if strings.Contains(key, "$"+envvar) {
		key = strings.Replace(key, "$"+envvar, envval, -1)
	}
	if strings.Contains(key, "${"+envvar+"}") {
		key = strings.Replace(key, "${"+envvar+"}", envval, -1)
	}
	return key
}
