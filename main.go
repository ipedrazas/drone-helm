package main

import (
	"fmt"
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
			EnvVar: "PLUGIN_HELM_COMMAND",
		},
		cli.StringFlag{
			Name:   "api_server",
			Usage:  "Api Server url",
			EnvVar: "PLUGIN_API_SERVER",
		},
		cli.StringFlag{
			Name:   "token",
			Usage:  "Kubernetes Token",
			EnvVar: "PLUGIN_TOKEN",
		},
		cli.StringFlag{
			Name:   "namespace",
			Usage:  "Kubernetes namespace",
			EnvVar: "PLUGIN_NAMESPACE",
		},
		cli.BoolFlag{
			Name:   "skip_tls_verify",
			Usage:  "Skip TLS verification",
			EnvVar: "PLUGIN_SKIP_TLS_VERIFY",
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
		},
	}
	return plugin.Exec()
}
