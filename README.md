# Helm (Kubernetes) plugin for drone.io

[![Build Status](https://drone.pelo.tech/api/badges/josmo/drone-helm/status.svg)](https://drone.pelo.tech/josmo/drone-helm)
[![Go Doc](https://godoc.org/github.com/josmo/drone-helm?status.svg)](http://godoc.org/github.com/josmo/drone-helm)
[![Go Report](https://goreportcard.com/badge/github.com/josmo/drone-helm)](https://goreportcard.com/report/github.com/josmo/drone-helm)
[![](https://images.microbadger.com/badges/image/peloton/drone-helm.svg)](https://microbadger.com/images/peloton/drone-helm "Get your own image badge on microbadger.com")

This plugin allows to deploy a [Helm](https://github.com/kubernetes/helm) chart into a [Kubernetes](https://github.com/kubernetes/kubernetes) cluster.

* Current `helm` version: 2.9.1
* Current `kubectl` version: 1.11.0

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
  -e PLUGIN_DEBUG=true \
  -e PLUGIN_DRY_RUN=true \
  -e DRONE_BUILD_EVENT=push \
  quay.io/ipedrazas/drone-helm
```

### Contribution

TBD

### Fork Notes

This is currently a fork of ipedrazas/drone-helm which is published to quai. The fork image will be published to dockerhub at peloton/drone-helm in the docs just replace until we know the long term path. 
