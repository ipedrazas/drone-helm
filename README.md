# Helm (Kubernetes) plugin for drone.io

[![Build Status](https://build.kube.camp/api/badges/ipedrazas/drone-helm/status.svg)](https://build.kube.camp/ipedrazas/drone-helm)

This plugin allows to deploy a [Helm](https://github.com/kubernetes/helm) chart into a [Kubernetes](https://github.com/kubernetes/kubernetes) cluster.

* Current `helm` version: 2.6.0
* Current `kubectl` version: 1.6.6

## Drone Pipeline Usage

### Simple Usage

For example, this configuration will deploy my-app using a chart located in the repo called `my-chart`

```YAML
pipeline:
  helm_deploy:
    image: quay.io/ipedrazas/drone-helm
    skip_tls_verify: true
    chart: ./charts/my-chart
    release: ${DRONE_BRANCH}
    values: secret.password=${SECRET_PASSWORD},image.tag=${DRONE_BRANCH}-${DRONE_COMMIT_SHA:0:7}
    prefix: STAGING
    debug: true
    wait: true
    when:
      branch: [master]
```

Last update of Drone expect you to declare the secrets you want to use:

```YAML
  helm_deploy:
    image: quay.io/ipedrazas/drone-helm
    chart: ./chart/blog
    release: ${DRONE_BRANCH}-blog
    values: image.tag=${DRONE_BRANCH}-${DRONE_COMMIT_SHA:0:7}
    prefix: PROD
    secrets: [ prod_api_server, prod_kubernetes_token ]
    when:
      branch: [master]
```

### Using Values and Value files

Values can be passed using the `values_files` key. Use this option to define your values in a set of files
and pass them to `helm`. This option trigger the `-f` or ``--values`` flag in `helm`:

```
--values valueFiles   specify values in a YAML file (can specify multiple) (default [])
```

For example:

```YAML
pipeline:
  helm_deploy:
    image: quay.io/ipedrazas/drone-helm
    skip_tls_verify: true
    chart: ./charts/my-chart
    release: ${DRONE_BRANCH}
    values_files: ["global-values.yaml", "myenv-values.yaml"]
    when:
      branch: [master]
```

### Using private Repositories

Charts can also be fetched from your own private Chart Repository. `helm_repos` accepts a comma separated list of key value pairs where the key is the repository name and the value is the repository url.

For Example:
```
helm_deploy_staging:
    image: quay.io/ipedrazas/drone-helm
    skip_tls_verify: true
    helm_repos: hb-charts=http://helm-charts.honestbee.com
    chart: hb-charts/hello-world
    values: image.repository=quay.io/honestbee/hello-drone-helm,image.tag=${DRONE_BRANCH}-${DRONE_COMMIT_SHA:0:7}
    release: ${DRONE_REPO_NAME}-${DRONE_BRANCH}
    prefix: STAGING
    when:
      branch:
        exclude: [ master ]
```

## Drone Secrets

There are two secrets you have to create (Note that if you specify the prefix, your secrets have to be created using that prefix):

```Bash
drone secret add --image=quay.io/ipedrazas/drone-helm \
  your-user/your-repo STAGING_API_SERVER https://mykubernetesapiserver

drone secret add --image=quay.io/ipedrazas/drone-helm \
  your-user/your-repo STAGING_KUBERNETES_TOKEN eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJrdWJ...

drone secret add --image=quay.io/ipedrazas/drone-helm \
  your-user/your-repo STAGING_SECRET_PASSWORD Sup3rS3cr3t
```

`Prefix` helps you to use the same block in different environments:
```YAML
pipeline:
  helm_deploy_staging:
    image: quay.io/ipedrazas/drone-helm
    skip_tls_verify: true
    chart: ./charts/my-chart
    release: ${DRONE_BRANCH}
    values: secret.password=${SECRET_PASSWORD},image.tag=${DRONE_BRANCH}-${DRONE_COMMIT_SHA:0:7}
    prefix: STAGING
    debug: true
    wait: true
    when:
      branch:
        exclude: [ master ]

pipeline_production:
  helm_deploy:
    image: quay.io/ipedrazas/drone-helm
    skip_tls_verify: true
    chart: ./charts/my-chart
    release: ${DRONE_BRANCH}
    values: secret.password=${SECRET_PASSWORD},image.tag=${DRONE_BRANCH}-${DRONE_COMMIT_SHA:0:7}
    prefix: PROD
    debug: true
    wait: true
    when:
      branch: [master]
```

This last block defines how the plugin will deploy


## Testing with Minikube

To test the plugin, you can run `minikube` and just run the docker image as follows:

By using the docker daemon of minikube we can test local builds without having to push to a registry:

```bash
eval $(minikube docker-env)
```

Build the image locally

```bash
./build.sh
```

Get the token for the default service account in the default namespace:

```bash
KUBERNETES_TOKEN=$(kubectl get secret $(kubectl get sa default -o jsonpath='{.secrets[].name}{"\n"}') -o jsonpath="{.data.token}" | base64 -D)
```

Run the local image (or replace `drone-helm` with `quay.io/honestbee/drone-helm`:
```Bash
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

## Advanced customisations and debugging

This plugin installs [Tiller](https://github.com/kubernetes/helm/blob/master/docs/architecture.md) in the cluster, if you want to specify the namespace where `tiller` ins installed, use the `tiller_ns` attribute.

The following example will install `tiller` in the `operations` namespace:
```YAML
pipeline_production:
  helm_deploy:
    image: quay.io/ipedrazas/drone-helm
    skip_tls_verify: true
    chart: ./charts/my-chart
    release: ${DRONE_BRANCH}
    values: image.tag=${DRONE_BRANCH}-${DRONE_COMMIT_SHA:0:7}
    prefix: PROD
    tiller_ns: operations
    when:
      branch: [master]
```

There's an option to do a `dry-run` in case you want to verify that the secrets and envvars are replaced correctly. Just add the attribute `dry-run` to true:

```YAML
pipeline_production:
  helm_deploy:
    image: quay.io/ipedrazas/drone-helm
    skip_tls_verify: true
    chart: ./charts/my-chart
    release: ${DRONE_BRANCH}
    values: image.tag=${DRONE_BRANCH}-${DRONE_COMMIT_SHA:0:7}
    prefix: STAGING
    dry-run:true
    when:
      branch: [master]
```
Happy Helming!
