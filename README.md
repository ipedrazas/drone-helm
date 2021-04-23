# Drone helm MoneySmart
This is a forked repository from https://github.com/ipedrazas/drone-helm
The image built is used in drone tasks to deploy helm charts to EKS.

Two changes have been introduced to the base repo
1. The Helm version has been hardcoded to use version 2.11
2. The stable repo url default value has been set to "https://charts.helm.sh/stable" since helm version needs updating for the [actual fix](https://github.com/hashicorp/terraform-provider-helm/issues/649) for this.

Image has been built and pushed to [dockerhub](https://hub.docker.com/repository/docker/moneysmartco/drone-helm) manually.

To build and upload image run
```
docker build -t moneysmartco/drone-helm:tag-name
docker push moneysmartco/drone-helm:tag-name
```

# Helm (Kubernetes) plugin for drone.io

[![Build Status](https://cloud.drone.io/api/badges/ipedrazas/drone-helm/status.svg)](https://cloud.drone.io/ipedrazas/drone-helm)
[![Docker Repository on Quay](https://quay.io/repository/ipedrazas/drone-helm/status "Docker Repository on Quay")](https://quay.io/repository/ipedrazas/drone-helm)
[![Go Doc](https://godoc.org/github.com/ipedrazas/drone-helm?status.svg)](http://godoc.org/github.com/ipedrazas/drone-helm)
[![Go Report](https://goreportcard.com/badge/github.com/ipedrazas/drone-helm)](https://goreportcard.com/report/github.com/ipedrazas/drone-helm)
[![](https://images.microbadger.com/badges/image/ipedrazas/drone-helm.svg)](https://microbadger.com/images/ipedrazas/drone-helm "Get your own image badge on microbadger.com")

This plugin allows to deploy a [Helm](https://github.com/kubernetes/helm) chart into a [Kubernetes](https://github.com/kubernetes/kubernetes) cluster.

* Current `helm` version: 3.4.2
* Current `kubectl` version: 1.21.0

## Drone Pipeline Usage

For the usage information and a listing of the available options please take a look at [the docs](DOCS.md).


Run the local image (or replace `drone-helm` with `quay.io/ipedrazas/drone-helm`:

```bash
docker run --rm \
  -e API_SERVER="https://$(minikube ip):8443" \
  -e KUBERNETES_TOKEN="${KUBERNETES_TOKEN}" \
  -e PLUGIN_NAMESPACE=default \
  -e PLUGIN_SKIP_TLS_VERIFY=true \
  -e PLUGIN_RELEASE=my-release \
  -e PLUGIN_CHART=stable/redis \
  -e PLUGIN_VALUES="tag=TAG,api=API" \
  -e PLUGIN_STRING_VALUES="long_string_value=1234567890" \
  -e PLUGIN_DEBUG=true \
  -e PLUGIN_DRY_RUN=true \
  -e DRONE_BUILD_EVENT=push \
  quay.io/ipedrazas/drone-helm
```

### Troubleshooting

If you see this problem: `Error: API Server is needed to deploy.` It's usually because you haven't a secret that specifies the `API_SERVER` or the `KUBERNETES_TOKEN`.

As [one000mph](https://github.com/one000mph) commented in an issue, setting the right `PREFIX` and secrets usually solves the problem.

```
export ACTION=add
    export REPO=org/myrepo
    export PREFIX=prod_
    # export CLUSTER_URI, UNENCODED_TOKEN, BASE64_CERT
    drone secret $ACTION --repository $REPO --name "${PREFIX}api_server" --value $CLUSTER_URI
    drone secret $ACTION --repository $REPO --name "${PREFIX}kubernetes_token" --value $UNENCODED_TOKEN
    drone secret $ACTION --repository $REPO --name "${PREFIX}kubernetes_certificate" --value $BASE64_CERT```
```

### Contribution

This repo is setup in a way that if you enable a personal drone server to build your fork it will
 build and publish your image (makes it easier to test PRs and use the image till the contributions get merged)
 
* Build local ```DRONE_REPO_OWNER=ipedrazas DRONE_REPO_NAME=drone-helm drone exec```
* on your server just make sure you have DOCKER_USERNAME, DOCKER_PASSWORD, and DOCKERHUB_REPO set as secrets
