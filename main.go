package main

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"os/exec"

	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-plugin-go/plugin"
)

var build = "0" // build number set at compile-time

// PluginParams to execute Helm
type PluginParams struct {
	Command         []string `json:"command"`
	APIServer       string   `json:"api_server"`
	KubernetesToken string   `json:"kubernetes_token"`
	Namespace       string   `json:"namespace"`
	Debug           string   `json:"debug"`
	Kubeconfig      string   `json:"kubeconfig"`
	SkipTLSVerify   string   `json:"skip_tls_verify"`
}

func isValidConfig(params *PluginParams) bool {
	if params.APIServer == "" {
		return false
	}
	if params.KubernetesToken == "" {
		return false
	}
	return true
}
func initialiseKubeconfig(params *PluginParams) {
	if params.Kubeconfig == "" {
		params.Kubeconfig = "/root/.kube/config"
	}
	if isValidConfig(params) {
		t, _ := template.ParseFiles("/root/.kube/kubeconfig")
		f, err := os.Create(params.Kubeconfig)
		if err != nil {
			log.Println("create file: ", err)
			return
		}
		err = t.Execute(f, params)
		if err != nil {
			log.Print("execute: ", err)
			return
		}
		f.Close()
	}
}

func main() {
	var (
		repo         = new(drone.Repo)
		build        = new(drone.Build)
		sys          = new(drone.System)
		pluginParams = new(PluginParams)
	)

	plugin.Param("build", build)
	plugin.Param("repo", repo)
	plugin.Param("system", sys)
	plugin.Param("vargs", pluginParams)
	plugin.MustParse()

	initialiseKubeconfig(pluginParams)
	init := make([]string, 1)
	init[0] = "init"
	runCommand("/bin/helm", init)
	runCommand("/bin/helm", pluginParams.Command)

}

func runCommand(command string, params []string) {
	cmd := new(exec.Cmd)
	cmd = exec.Command(command, params...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Printf("%s\n", err)
	}
}
