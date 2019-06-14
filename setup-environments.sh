#!/bin/bash

# will populate ~/.local/lib/helm-$version/ with helm and tiller binaries
helm_arch="linux-amd64"
helm_versions=(v2.14.1 v2.13.1 v2.12.3)

for version in "${helm_versions[@]}"
do
    curl http://storage.googleapis.com/kubernetes-helm/helm-${version}-${helm_arch}.tar.gz > /tmp/${version}.tar.gz
    mkdir -p ~/.local/lib/helm-${version}/
    tar -C ~/.local/lib/helm-${version}/ -xvf /tmp/${version}.tar.gz --strip 1
done

# will populate ~/.local/lib/kubectl-$version/ with kubectl binaries
kubectl_arch="linux/amd64"
kubectl_versions=(v1.14.3 v1.13.7 v1.12.9)

for version in "${kubectl_versions[@]}"
do
    mkdir -p ~/.local/lib/kubectl-${version}/
    curl https://storage.googleapis.com/kubernetes-release/release/${version}/bin/${kubectl_arch}/kubectl > ~/.local/lib/kubectl-${version}/kubectl
    chmod +x ~/.local/lib/kubectl-${version}/kubectl
done
