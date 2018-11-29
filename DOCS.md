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
pipeline:
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

Use Kubernetes Certificate Authority Data. Just add the `<prefix>_kubernetes_certificate` secret

```diff
  helm_deploy:
    image: quay.io/ipedrazas/drone-helm
    chart: ./chart/blog
    release: ${DRONE_BRANCH}-blog
    values: image.tag=${DRONE_BRANCH}-${DRONE_COMMIT_SHA:0:7}
    prefix: PROD
    - secrets: [ prod_api_server, prod_kubernetes_token ]
    + secrets: [ prod_api_server, prod_kubernetes_token, prod_kubernetes_certificate ]
    when:
      branch: [master]
```

### Using Values and Value files

Values can be passed using the `values_files` key. Use this option to define your values in a set of files
and pass them to `helm`. This option trigger the `-f` or ``--values`` flag in `helm`:

```plain
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

```YAML
pipeline:
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

## Updating Chart dependencies

In some cases, the local Chart might contain external dependencies defined in `./charts/my-chart/requirements.yaml`, e.g.:

```YAML
dependencies:
  - name: redis
    version: 3.3.6
    repository: '@stable'
```

To restore these dependecies before the deployment `update_dependencies` parameter should be used, e.g.:

```YAML
pipeline:
  helm_deploy:
    image: quay.io/ipedrazas/drone-helm
    skip_tls_verify: true
    chart: ./charts/my-chart
    update_dependencies: true
    release: ${DRONE_BRANCH}
    values_files: ["global-values.yaml", "myenv-values.yaml"]
    when:
      branch: [master]
```

## Drone Secrets

There are two secrets you have to create (Note that if you specify the prefix, your secrets have to be created using that prefix):

```bash
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

## Deploying to EKS

To deploy to EKS, you should have `api_server` and `kubernetes_certificate` secrets set in drone. If drone is deploying from outside of AWS, you should also have an `aws_access_key_id`, `aws_secret_access_key` secret. An `AWS_DEFAULT_REGION` environmental variable should also be set for the deployment.

### Use the AWS IAM keys for the user who created the EKS cluster (not recommended):

```YAML
pipeline
  helm_deploy:
    image: quay.io/ipedrazas/drone-helm
    chart: ./charts/my-chart
    release: ${DRONE_BRANCH}
    values: image.tag=${DRONE_BRANCH}-${DRONE_COMMIT_SHA:0:7}
    prefix: STAGING
    namespace: staging
    eks_cluster: my-eks-cluster
    environment:
      - AWS_DEFAULT_REGION=us-east-1

    secrets: [ aws_access_key_id, aws_secret_access_key, api_server, kubernetes_certificate ]
```

### Use role based EKS access
You must first [configure an IAM Role for cluster access](https://docs.aws.amazon.com/eks/latest/userguide/add-user-role.html) and configure policy to allow an iam user or iam role to assume the new role.

Running drone agent on an ec2 instance that has Role based access to assume the cluster access Role created above:
```YAML
pipeline
  helm_deploy:
    image: quay.io/ipedrazas/drone-helm
    chart: ./charts/my-chart
    release: ${DRONE_BRANCH}
    values: image.tag=${DRONE_BRANCH}-${DRONE_COMMIT_SHA:0:7}
    prefix: STAGING
    namespace: staging
    eks_cluster: my-eks-cluster
    eks_role_arn: arn:aws:iam::[ACCOUNT ID HERE]:role/eks-master
    environment:
      - AWS_DEFAULT_REGION=us-east-1

    secrets: [ api_server, kubernetes_certificate ]
```

Using IAM keys with access to assume the cluster access Role created above:
```YAML
pipeline
  helm_deploy:
    image: quay.io/ipedrazas/drone-helm
    chart: ./charts/my-chart
    release: ${DRONE_BRANCH}
    values: image.tag=${DRONE_BRANCH}-${DRONE_COMMIT_SHA:0:7}
    prefix: STAGING
    namespace: staging
    eks_cluster: my-eks-cluster
    eks_role_arn: arn:aws:iam::[ACCOUNT ID HERE]:role/eks-master
    environment:
      - AWS_DEFAULT_REGION=us-east-1

    secrets: [ aws_access_key_id, aws_secret_access_key, api_server, kubernetes_certificate ]
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
    dry-run: true
    when:
      branch: [master]
```

This plugin init stable repository in the cluster, if you want to specify the stable repository, use the `stable_repo_url` attribute.

The following example will init `stable_repo_url` in the `https://kubernetes.oss-cn-hangzhou.aliyuncs.com/charts` repo:

```YAML
pipeline_production:
  helm_deploy:
    image: quay.io/ipedrazas/drone-helm
    skip_tls_verify: true
    chart: ./charts/my-chart
    release: ${DRONE_BRANCH}
    values: image.tag=${DRONE_BRANCH}-${DRONE_COMMIT_SHA:0:7}
    prefix: PROD
    stable_repo_url: https://kubernetes.oss-cn-hangzhou.aliyuncs.com/charts
    when:
      branch: [master]
```


Happy Helming!

## Known issues

* Drone secrets that are part of `values` can be leaked in debug mode and in case of error as the whole helm command will be printed in the logs. See #52

