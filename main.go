package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/joho/godotenv"
	"github.com/urfave/cli"
)

var build = "0" // build number set at compile-time

func main() {
	app := cli.NewApp()
	app.Name = "helm plugin"
	app.Usage = "helm plugin"
	app.Action = run
	app.Version = fmt.Sprintf("1.0.%s", build)
	app.Flags = []cli.Flag{
		cli.StringSliceFlag{
			Name:   "helm_command",
			Usage:  "add the command Helm has to execute",
			EnvVar: "PLUGIN_HELM_COMMAND,HELM_COMMAND",
		},
		cli.StringFlag{
			Name:   "api_server",
			Usage:  "Api Server url",
			EnvVar: "PLUGIN_API_SERVER,API_SERVER",
		},
		cli.StringFlag{
			Name:   "token",
			Usage:  "Kubernetes Token",
			EnvVar: "PLUGIN_TOKEN,KUBERNETES_TOKEN",
		},
		cli.StringFlag{
			Name:   "namespace",
			Usage:  "Kubernetes namespace",
			EnvVar: "PLUGIN_NAMESPACE,NAMESPACE",
		},
		cli.StringFlag{
			Name:   "release",
			Usage:  "Kubernetes helm release",
			EnvVar: "PLUGIN_RELEASE,RELEASE",
		},
		cli.StringFlag{
			Name:   "chart",
			Usage:  "Kubernetes helm release",
			EnvVar: "PLUGIN_CHART,CHART",
		},
		cli.StringFlag{
			Name:   "values",
			Usage:  "Kubernetes helm release",
			EnvVar: "PLUGIN_VALUES,VALUES",
		},
		cli.BoolFlag{
			Name:   "skip_tls_verify",
			Usage:  "Skip TLS verification",
			EnvVar: "PLUGIN_SKIP_TLS_VERIFY,SKIP_TLS_VERIFY",
		},
		cli.BoolFlag{
			Name:   "debug",
			Usage:  "Debug",
			EnvVar: "PLUGIN_DEBUG,DEBUG",
		},
		cli.BoolFlag{
			Name:   "dry-run",
			Usage:  "Helm dry-run",
			EnvVar: "PLUGIN_DRY_RUN,DRY_RUN",
		},
		cli.StringSliceFlag{
			Name:   "secrets",
			Usage:  "add the secrets used in the values field",
			EnvVar: "PLUGIN_SECRETS,SECRETS",
		},
		cli.StringSliceFlag{
			Name:   "prefix",
			Usage:  "Prefix for all the secrets",
			EnvVar: "PLUGIN_PREFIX,PREFIX",
		},
	}
	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func run(c *cli.Context) error {
	if c.String("env-file") != "" {
		_ = godotenv.Load(c.String("env-file"))
	}
	plugin := Plugin{
		Config: Config{
			APIServer:     c.String("api_server"),
			Token:         c.String("token"),
			HelmCommand:   c.StringSlice("helm_command"),
			Namespace:     c.String("namespace"),
			SkipTLSVerify: c.Bool("skip_tls_verify"),
			Values:        c.String("values"),
			Release:       c.String("release"),
			Chart:         c.String("chart"),
			Debug:         c.Bool("debug"),
			DryRun:        c.Bool("dry-run"),
			Secrets:       c.StringSlice("secrets"),
			Prefix:        c.String("prefix"),
		},
	}
	if plugin.Config.Debug {
		debug(&plugin)
	}

	resolveSecrets(&plugin)

	if plugin.Config.Debug {
		debug(&plugin)
	}

	return plugin.Exec()
}

func debug(plugin *Plugin) {
	// debug env vars
	for _, e := range os.Environ() {
		fmt.Println(e)
	}
	// debug plugin obj
	fmt.Printf("Api server: %s \n", plugin.Config.APIServer)
	fmt.Printf("Values: %s \n", plugin.Config.Values)
	fmt.Printf("Values: %s \n", plugin.Config.Secrets)

	kubeconfig, err := ioutil.ReadFile(KUBECONFIG)
	if err != nil {
		fmt.Println(string(kubeconfig))
	}
	config, err := ioutil.ReadFile(CONFIG)
	if err != nil {
		fmt.Println(string(config))
	}

}
