package main

import (
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/ipedrazas/drone-helm/plugin"
	"github.com/joho/godotenv"
	"github.com/urfave/cli"
)

var (
	version = "0.0.0"
	build   = "0"
)

func main() {
	app := cli.NewApp()
	app.Name = "helm plugin"
	app.Usage = "helm plugin"
	app.Action = run
	app.Version = fmt.Sprintf("%s+%s", version, build)
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "helm_command",
			Usage:  "add the command Helm has to execute",
			EnvVar: "PLUGIN_HELM_COMMAND,HELM_COMMAND",
		},
		cli.StringFlag{
			Name:   "kube-config",
			Usage:  "Kubernetes configuration file path",
			EnvVar: "PLUGIN_KUBE_CONFIG,KUBE_CONFIG",
			Value:  "/root/.kube/config",
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
		cli.StringSliceFlag{
			Name:   "helm_repos",
			Usage:  "Repos helm should add",
			EnvVar: "PLUGIN_HELM_REPOS,HELM_REPOS",
		},
		cli.StringFlag{
			Name:   "chart",
			Usage:  "Kubernetes helm chart name",
			EnvVar: "PLUGIN_CHART,CHART",
		},
		cli.StringFlag{
			Name:   "chart-version",
			Usage:  "specify the exact chart version to use. If this is not specified, the latest version is used",
			EnvVar: "PLUGIN_CHART_VERSION,CHART_VERSION",
		},
		cli.StringFlag{
			Name:   "eks_cluster",
			Usage:  "Name of EKS cluster. Requires AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, and AWS_DEFAULT_REGION secrets AND/OR EKS_ROLE and proper role configuration",
			EnvVar: "PLUGIN_EKS_CLUSTER,EKS_CLUSTER",
		},
		cli.StringFlag{
			Name:   "eks_role_arn",
			Usage:  "ARN of EKS role to assume for EKS authentication.",
			EnvVar: "PLUGIN_EKS_ROLE_ARN,EKS_ROLE_ARN",
		},
		cli.StringFlag{
			Name:   "values",
			Usage:  "Kubernetes helm release",
			EnvVar: "PLUGIN_VALUES,VALUES",
		},
		cli.StringFlag{
			Name:   "string_values",
			Usage:  "Kubernetes helm release",
			EnvVar: "PLUGIN_STRING_VALUES,STRING_VALUES",
		},
		cli.StringFlag{
			Name:   "values_files",
			Usage:  "Helm values override files",
			EnvVar: "PLUGIN_VALUES_FILES,VALUES_FILES",
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
		cli.StringFlag{
			Name:   "prefix",
			Usage:  "Prefix for all the secrets",
			EnvVar: "PLUGIN_PREFIX,PREFIX",
		},
		cli.StringFlag{
			Name:   "tiller-ns",
			Usage:  "Namespace to install Tiller",
			EnvVar: "PLUGIN_TILLER_NS,TILLER_NS",
		},
		cli.BoolFlag{
			Name:   "wait",
			Usage:  "if set, will wait until all Pods, PVCs, and Services are in a ready state before marking the release as successful.",
			EnvVar: "PLUGIN_WAIT,WAIT",
		},
		cli.BoolFlag{
			Name:   "recreate-pods",
			Usage:  "performs pods restart for the resource if applicable",
			EnvVar: "PLUGIN_RECREATE_PODS,RECREATE_PODS",
		},
		cli.BoolFlag{
			Name:   "upgrade",
			Usage:  "if set, will upgrade tiller to the latest version",
			EnvVar: "PLUGIN_UPGRADE,UPGRADE",
		},
		cli.BoolFlag{
			Name:   "canary-image",
			Usage:  "if set, Helm will use the canary tiller image",
			EnvVar: "PLUGIN_CANARY_IMAGE,CANARY_IMAGE",
		},
		cli.BoolFlag{
			Name:   "client-only",
			Usage:  "if set, it will initilises helm in the client side only",
			EnvVar: "PLUGIN_CLIENT_ONLY,CLIENT_ONLY",
		},
		cli.BoolFlag{
			Name:   "reuse-values",
			Usage:  "when upgrading, reuse the last release's values, and merge in any new values",
			EnvVar: "PLUGIN_REUSE_VALUES,REUSE_VALUES",
		},
		cli.StringFlag{
			Name:   "timeout",
			Usage:  "time in seconds to wait for any individual kubernetes operation (like Jobs for hooks) (default 300)",
			EnvVar: "PLUGIN_TIMEOUT,TIMEOUT",
		},
		cli.BoolFlag{
			Name:   "force",
			Usage:  "force resource update through delete/recreate if needed",
			EnvVar: "PLUGIN_FORCE,FORCE",
		},
		cli.BoolFlag{
			Name:   "purge",
			Usage:  "purge on delete",
			EnvVar: "PLUGIN_PURGE,PURGE",
		},
		cli.BoolFlag{
			Name:   "update-dependencies",
			Usage:  "update dependency charts based on the contents of requirements.yaml file of the local chart",
			EnvVar: "PLUGIN_UPDATE_DEPENDENCIES,UPDATE_DEPENDENCIES",
		},
		cli.StringFlag{
			Name:   "stable_repo_url",
			Usage:  "URL for stable repository (default 'https://charts.helm.sh/stable')",
			EnvVar: "PLUGIN_STABLE_REPO_URL,STABLE_REPO_URL",
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
	p := plugin.Plugin{
		Config: plugin.Config{
			APIServer:          c.String("api_server"),
			Token:              c.String("token"),
			Certificate:        c.String("certificate"),
			ServiceAccount:     c.String("service-account"),
			KubeConfig:         c.String("kube-config"),
			HelmCommand:        c.String("helm_command"),
			Namespace:          c.String("namespace"),
			SkipTLSVerify:      c.Bool("skip_tls_verify"),
			Values:             c.String("values"),
			StringValues:       c.String("string_values"),
			ValuesFiles:        c.String("values_files"),
			Release:            c.String("release"),
			HelmRepos:          c.StringSlice("helm_repos"),
			Chart:              c.String("chart"),
			Version:            c.String("chart-version"),
			EKSCluster:         c.String("eks_cluster"),
			EKSRoleARN:         c.String("eks_role_arn"),
			Debug:              c.Bool("debug"),
			DryRun:             c.Bool("dry-run"),
			Secrets:            c.StringSlice("secrets"),
			Prefix:             c.String("prefix"),
			TillerNs:           c.String("tiller-ns"),
			Wait:               c.Bool("wait"),
			RecreatePods:       c.Bool("recreate-pods"),
			ClientOnly:         c.Bool("client-only"),
			CanaryImage:        c.Bool("canary-image"),
			Upgrade:            c.Bool("upgrade"),
			ReuseValues:        c.Bool("reuse-values"),
			Timeout:            c.String("timeout"),
			Force:              c.Bool("force"),
			UpdateDependencies: c.Bool("update-dependencies"),
			StableRepoURL:      c.String("stable_repo_url"),
		},
	}
	return p.Exec()
}
