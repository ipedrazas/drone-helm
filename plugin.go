package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
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
		ValuesFiles   string   `json:"values_files"`
		Debug         bool     `json:"debug"`
		DryRun        bool     `json:"dry_run"`
		Secrets       []string `json:"secrets"`
		Prefix        string   `json:"prefix"`
		TillerNs      string   `json:"tiller_ns"`
		Wait          bool     `json:"wait"`
		Upgrade       bool     `json:"upgrade"`
		ClientOnly    bool     `json:"client_only"`
	}
	// Plugin default
	Plugin struct {
		Config Config
	}
)

func setHelmHelp(p *Plugin) {
	p.Config.HelmCommand = []string{""}
}
func setDeleteEventCommand(p *Plugin) {
	upgrade := make([]string, 2)
	upgrade[0] = "delete"
	upgrade[1] = p.Config.Release

	p.Config.HelmCommand = upgrade
}

func setPushEventCommand(p *Plugin) {
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
	if p.Config.ValuesFiles != "" {
		for _, valuesFile := range strings.Split(p.Config.ValuesFiles, ",") {
			upgrade = append(upgrade, "--values")
			upgrade = append(upgrade, valuesFile)
		}
	}
	if p.Config.Namespace != "" {
		upgrade = append(upgrade, "--namespace")
		upgrade = append(upgrade, p.Config.Namespace)
	}
	if p.Config.TillerNs != "" {
		upgrade = append(upgrade, "--tiller-namespace")
		upgrade = append(upgrade, p.Config.TillerNs)
	}
	if p.Config.DryRun {
		upgrade = append(upgrade, "--dry-run")
	}
	if p.Config.Debug {
		upgrade = append(upgrade, "--debug")
	}
	if p.Config.Wait {
		upgrade = append(upgrade, "--wait")
	}
	p.Config.HelmCommand = upgrade

}

func setHelmCommand(p *Plugin) {
	buildEvent := os.Getenv("DRONE_BUILD_EVENT")
	switch buildEvent {
	case "push":
		setPushEventCommand(p)
	case "tag":
		setPushEventCommand(p)
	case "deployment":
		setPushEventCommand(p)
	case "delete":
		setDeleteEventCommand(p)
	default:
		setHelmHelp(p)
	}

}

func doHelmInit(p *Plugin) []string {
	init := make([]string, 1)
	init[0] = "init"
	if p.Config.TillerNs != "" {
		init = append(init, "--tiller-namespace")
		init = append(init, p.Config.TillerNs)
	}
	if p.Config.ClientOnly {
		init = append(init, "--client-only")
	}
	if p.Config.Upgrade {
		init = append(init, "--upgrade")
	}

	return init

}

// Exec default method
func (p *Plugin) Exec() error {
	resolveSecrets(p)
	if p.Config.APIServer == "" {
		return fmt.Errorf("Error: API Server is needed to deploy.")
	}
	if p.Config.Token == "" {
		return fmt.Errorf("Error: Token is needed to deploy.")
	}
	initialiseKubeconfig(&p.Config, KUBECONFIG, CONFIG)

	if p.Config.Debug {
		p.debug()
	}

	init := doHelmInit(p)
	err := runCommand(init)
	if err != nil {
		return fmt.Errorf("Error running helm command: " + strings.Join(init[:], " "))
	}
	setHelmCommand(p)

	if p.Config.Debug {
		log.Println("helm command: " + strings.Join(p.Config.HelmCommand[:], " "))
	}
	err = runCommand(p.Config.HelmCommand)
	if err != nil {
		return fmt.Errorf("Error running helm command: " + strings.Join(p.Config.HelmCommand[:], " "))
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
	p.Config.Values = resolveEnvVar(p.Config.Values, p.Config.Prefix)
	p.Config.APIServer = resolveEnvVar("${API_SERVER}", p.Config.Prefix)
	p.Config.Token = resolveEnvVar("${KUBERNETES_TOKEN}", p.Config.Prefix)
}

// getEnvVars will return [${TAG} {TAG} TAG]
func getEnvVars(envvars string) [][]string {
	re := regexp.MustCompile(`\$(\{?(\w+)\}?)\.?`)
	extracted := re.FindAllStringSubmatch(envvars, -1)
	return extracted
}

func resolveEnvVar(key string, prefix string) string {
	envvars := getEnvVars(key)
	return replaceEnvvars(envvars, prefix, key)
}

func replaceEnvvars(envvars [][]string, prefix string, s string) string {
	for _, envvar := range envvars {
		envvarName := envvar[0]
		envvarKey := envvar[2]
		envval := os.Getenv(prefix + "_" + envvarKey)
		if envval == "" {
			envval = os.Getenv(envvarKey)
		}
		if strings.Contains(s, envvarName) {
			s = strings.Replace(s, envvarName, envval, -1)
		}
	}
	return s
}

func (p *Plugin) debug() {
	fmt.Println(p)
	// debug env vars
	for _, e := range os.Environ() {
		fmt.Println("-Var:--", e)
	}
	// debug plugin obj
	fmt.Printf("Api server: %s \n", p.Config.APIServer)
	fmt.Printf("Values: %s \n", p.Config.Values)
	fmt.Printf("Secrets: %s \n", p.Config.Secrets)
	fmt.Printf("ValuesFiles: %s \n", p.Config.ValuesFiles)

	kubeconfig, err := ioutil.ReadFile(KUBECONFIG)
	if err == nil {
		fmt.Println(string(kubeconfig))
	}
	config, err := ioutil.ReadFile(CONFIG)
	if err == nil {
		fmt.Println(string(config))
	}

}
