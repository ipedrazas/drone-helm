#!/bin/bash

PREFIX=~/.local

# will populate $PREFIX/lib/helm-$version/ with helm and tiller binaries
helm_arch="linux-amd64"
helm_versions=(v2.14.1 v2.13.1 v2.12.3)

for version in "${helm_versions[@]}"
do
    curl http://storage.googleapis.com/kubernetes-helm/helm-${version}-${helm_arch}.tar.gz > /tmp/${version}.tar.gz
    mkdir -p $PREFIX/lib/helm-${version}/
    tar -C $PREFIX/lib/helm-${version}/ -xvf /tmp/${version}.tar.gz --strip 1
done

# will populate $PREFIX/lib/kubectl-$version/ with kubectl binaries
kubectl_arch="linux/amd64"
kubectl_versions=(v1.14.3 v1.13.7 v1.12.9)

for version in "${kubectl_versions[@]}"
do
    mkdir -p $PREFIX/lib/kubectl-${version}/
    curl https://storage.googleapis.com/kubernetes-release/release/${version}/bin/${kubectl_arch}/kubectl > $PREFIX/lib/kubectl-${version}/kubectl
    chmod +x $PREFIX/lib/kubectl-${version}/kubectl
done

