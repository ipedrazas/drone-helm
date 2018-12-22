package plugin

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"text/template"
)

var HELM_BIN = "/bin/helm"
var KUBECONFIG = "/root/.kube/kubeconfig"

type (
	// Config maps the params we need to run Helm
	Config struct {
		APIServer          string   `json:"api_server"`
		Token              string   `json:"token"`
		Certificate        string   `json:"certificate"`
		ServiceAccount     string   `json:"service_account"`
		KubeConfig         string   `json:"kube_config"`
		HelmCommand        string   `json:"helm_command"`
		SkipTLSVerify      bool     `json:"tls_skip_verify"`
		Namespace          string   `json:"namespace"`
		Release            string   `json:"release"`
		Chart              string   `json:"chart"`
		Version            string   `json:"version"`
		EKSCluster         string   `json:"eks_cluster"`
		EKSRoleARN         string   `json:"eks_role_arn"`
		Values             string   `json:"values"`
		StringValues       string   `json:"string_values"`
		ValuesFiles        string   `json:"values_files"`
		Debug              bool     `json:"debug"`
		DryRun             bool     `json:"dry_run"`
		Secrets            []string `json:"secrets"`
		Prefix             string   `json:"prefix"`
		TillerNs           string   `json:"tiller_ns"`
		Wait               bool     `json:"wait"`
		RecreatePods       bool     `json:"recreate_pods"`
		Upgrade            bool     `json:"upgrade"`
		CanaryImage        bool     `json:"canary_image"`
		ClientOnly         bool     `json:"client_only"`
		ReuseValues        bool     `json:"reuse_values"`
		Timeout            string   `json:"timeout"`
		Force              bool     `json:"force"`
		HelmRepos          []string `json:"helm_repos"`
		Purge              bool     `json:"purge"`
		UpdateDependencies bool     `json:"update_dependencies"`
		StableRepoURL      string   `json:"stable_repo_url"`
	}
	// Plugin default
	Plugin struct {
		Config  Config
		command []string
	}
)

func setHelpCommand(p *Plugin) {
	p.command = []string{""}
}
func setDeleteCommand(p *Plugin) {
	delete := make([]string, 2)
	delete[0] = "delete"
	delete[1] = p.Config.Release

	if p.Config.TillerNs != "" {
		delete = append(delete, "--tiller-namespace")
		delete = append(delete, p.Config.TillerNs)
	}
	if p.Config.DryRun {
		delete = append(delete, "--dry-run")
	}
	if p.Config.Purge {
		delete = append(delete, "--purge")
	}

	p.command = delete
}

func setUpgradeCommand(p *Plugin) {
	upgrade := make([]string, 2)
	upgrade[0] = "upgrade"
	upgrade[1] = "--install"

	if p.Config.Release != "" {
		upgrade = append(upgrade, p.Config.Release)
	}
	upgrade = append(upgrade, p.Config.Chart)
	if p.Config.Version != "" {
		upgrade = append(upgrade, "--version")
		upgrade = append(upgrade, p.Config.Version)
	}
	if p.Config.Values != "" {
		upgrade = append(upgrade, "--set")
		upgrade = append(upgrade, unQuote(p.Config.Values))
	}
	if p.Config.StringValues != "" {
		upgrade = append(upgrade, "--set-string")
		upgrade = append(upgrade, unQuote(p.Config.StringValues))
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
	if p.Config.RecreatePods {
		upgrade = append(upgrade, "--recreate-pods")
	}
	if p.Config.ReuseValues {
		upgrade = append(upgrade, "--reuse-values")
	}
	if p.Config.Timeout != "" {
		upgrade = append(upgrade, "--timeout")
		upgrade = append(upgrade, p.Config.Timeout)
	}
	if p.Config.Force {
		upgrade = append(upgrade, "--force")
	}
	p.command = upgrade
}

func setLintCommand(p *Plugin) {
	lint := make([]string, 2)
	lint[0] = "lint"
	lint[1] = p.Config.Chart

	if p.Config.Values != "" {
		lint = append(lint, "--set")
		lint = append(lint, unQuote(p.Config.Values))
	}

	if p.Config.StringValues != "" {
		lint = append(lint, "--set-string")
		lint = append(lint, unQuote(p.Config.StringValues))
	}

	if p.Config.ValuesFiles != "" {
		for _, valuesFile := range strings.Split(p.Config.ValuesFiles, ",") {
			lint = append(lint, "--values")
			lint = append(lint, valuesFile)
		}
	}

	if p.Config.Namespace != "" {
		lint = append(lint, "--namespace")
		lint = append(lint, p.Config.Namespace)
	}

	if p.Config.TillerNs != "" {
		lint = append(lint, "--tiller-namespace")
		lint = append(lint, p.Config.TillerNs)
	}

	if p.Config.Debug {
		lint = append(lint, "--debug")
	}

	p.command = lint
}

func setHelmCommand(p *Plugin) {

	switch p.Config.HelmCommand {
	case "upgrade":
		setUpgradeCommand(p)
	case "delete":
		setDeleteCommand(p)
	case "lint":
		setLintCommand(p)
	default:
		switch os.Getenv("DRONE_BUILD_EVENT") {
		case "push", "tag", "deployment", "pull_request", "promote", "rollback":
			setUpgradeCommand(p)
		case "delete":
			setDeleteCommand(p)
		default:
			setHelpCommand(p)
		}
	}

}

var repoExp = regexp.MustCompile(`^(?P<name>[\w-]+)=(?P<url>(http|https)://[\w-./:@-]+)`)

// parseRepo returns map of regex capture groups (name, url)
func parseRepo(repo string) (map[string]string, error) {
	matches := repoExp.FindStringSubmatch(repo)
	if len(matches) < 1 {
		return nil, fmt.Errorf("Invalid repo definition: %s", repo)
	}
	result := make(map[string]string)
	for i, name := range repoExp.SubexpNames() {
		if i != 0 {
			result[name] = matches[i]
		}
	}
	return result, nil
}

func doHelmRepoAdd(repo string) ([]string, error) {
	repoMap, err := parseRepo(unQuote(repo))
	if err != nil {
		return nil, err
	}
	repoAdd := []string{
		"repo",
		"add",
		repoMap["name"],
		repoMap["url"],
	}
	return repoAdd, nil
}

func doHelmInit(p *Plugin) []string {
	init := make([]string, 1)
	init[0] = "init"
	if p.Config.StableRepoURL != "" {
		init = append(init, "--stable-repo-url")
		init = append(init, p.Config.StableRepoURL)
	}
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
	if p.Config.CanaryImage {
		init = append(init, "--canary-image")
	}

	return init

}

func doDependencyUpdate(chart string) []string {
	dependencyUpdate := []string{
		"dependency",
		"update",
		chart,
	}

	return dependencyUpdate
}

// Exec default method
func (p *Plugin) Exec() error {
	if p.Config.Debug {
		p.debugEnv()
	}

	// create /root/.kube/config file if not exists
	if _, err := os.Stat(p.Config.KubeConfig); os.IsNotExist(err) {
		resolveSecrets(p)
		if p.Config.APIServer == "" {
			return fmt.Errorf("Error: API Server is needed to deploy.")
		}
		if p.Config.EKSCluster == "" {
			if p.Config.Token == "" {
				return fmt.Errorf("Error: Token is needed to deploy.")
			}
		}
		initialiseKubeconfig(&p.Config, KUBECONFIG, p.Config.KubeConfig)
	}

	if p.Config.Debug {
		p.debug()
	}

	init := doHelmInit(p)
	err := runCommand(init)
	if err != nil {
		return fmt.Errorf("Error running helm command: " + strings.Join(init[:], " "))
	}

	if len(p.Config.HelmRepos) > 0 {
		for _, repo := range p.Config.HelmRepos {
			repoAdd, err := doHelmRepoAdd(repo)
			if err == nil {
				if p.Config.Debug {
					log.Println("adding helm repo: " + strings.Join(repoAdd[:], " "))
				}

				if err = runCommand(repoAdd); err != nil {
					return fmt.Errorf("Error adding helm repo: " + err.Error())
				}
			} else {
				return err
			}
		}
	}

	if p.Config.UpdateDependencies {
		if err = runCommand(doDependencyUpdate(p.Config.Chart)); err != nil {
			return fmt.Errorf("Error updating dependencies: " + err.Error())
		}
	}

	setHelmCommand(p)

	if p.Config.Debug {
		log.Println("helm command: " + strings.Join(p.command, " "))
	}

	err = runCommand(p.command)
	if err != nil {
		return fmt.Errorf("Error running helm command: " + strings.Join(p.command[:], " "))
	}

	return nil
}

func initialiseKubeconfig(params *Config, source string, target string) error {
	f, err := os.Create(target)
	if err != nil {
		return err
	}
	defer f.Close()
	// parse template
	t, _ := template.ParseFiles(source)
	// execute template
	return t.Execute(f, params)
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
	p.Config.Values = resolveEnvVar(p.Config.Values, p.Config.Prefix, p.Config.Debug)
	p.Config.StringValues = resolveEnvVar(p.Config.StringValues, p.Config.Prefix, p.Config.Debug)

	if p.Config.APIServer == "" {
		p.Config.APIServer = resolveEnvVar("${API_SERVER}", p.Config.Prefix, p.Config.Debug)
	}
	if p.Config.Token == "" {
		p.Config.Token = resolveEnvVar("${KUBERNETES_TOKEN}", p.Config.Prefix, p.Config.Debug)
	}
	if p.Config.Certificate == "" {
		p.Config.Certificate = resolveEnvVar("${KUBERNETES_CERTIFICATE}", p.Config.Prefix, p.Config.Debug)
	}
	if p.Config.ServiceAccount == "" {
		p.Config.ServiceAccount = resolveEnvVar("${SERVICE_ACCOUNT}", p.Config.Prefix, p.Config.Debug)
		if p.Config.ServiceAccount == "" {
			p.Config.ServiceAccount = "helm"
		}
	}
}

// getEnvVars will return [${TAG} {TAG} TAG]
func getEnvVars(envvars string) [][]string {
	re := regexp.MustCompile(`\$(\{?(\w+)\}?)\.?`)
	extracted := re.FindAllStringSubmatch(envvars, -1)
	return extracted
}

func resolveEnvVar(key string, prefix string, debug bool) string {
	envvars := getEnvVars(key)
	return replaceEnvvars(envvars, prefix, key, debug)
}

func replaceEnvvars(envvars [][]string, prefix string, s string, debug bool) string {
	for _, envvar := range envvars {
		envvarName := envvar[0]
		envvarKey := envvar[2]
		prefixedKey := strings.ToUpper(prefix + "_" + envvarKey)
		envval := os.Getenv(prefixedKey)
		if debug {
			fmt.Printf("-ReplVar: %s => %s-- %s\n", prefixedKey, envvarKey, envval)
		}

		if envval == "" {
			envval = os.Getenv(envvarKey)
		}

		if strings.Contains(s, envvarName) {
			s = strings.Replace(s, envvarName, envval, -1)
		}
	}

	return s
}

// unQuote removes quotes if present
func unQuote(s string) string {
	unquoted, err := strconv.Unquote(s)
	if err != nil {
		// ignore error and return original string
		return s
	}
	return unquoted
}

func (p *Plugin) debugEnv() {
	// debug env vars
	for _, e := range os.Environ() {
		fmt.Println("-Var:--", e)
	}
}

func (p *Plugin) debug() {
	fmt.Println(p)
	// debug plugin obj
	fmt.Printf("Api server: %s \n", p.Config.APIServer)
	fmt.Printf("Values: %s \n", p.Config.Values)
	fmt.Printf("StringValues: %s \n", p.Config.StringValues)
	fmt.Printf("Secrets: %s \n", p.Config.Secrets)
	fmt.Printf("Helm Repos: %s \n", p.Config.HelmRepos)
	fmt.Printf("ValuesFiles: %s \n", p.Config.ValuesFiles)
	fmt.Printf("StableRepoURL: %s \n", p.Config.StableRepoURL)
	kubeconfig, err := ioutil.ReadFile(KUBECONFIG)
	if err == nil {
		fmt.Println(string(kubeconfig))
	}

	config, err := ioutil.ReadFile(p.Config.KubeConfig)
	if err == nil {
		fmt.Println(string(config))
	}
}
