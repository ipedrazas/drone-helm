#!/bin/sh

function error {
    echo "$1"
    exit 1
}  
                
# Will symlink `helm` and `kubectl` binaries based on the `$HELM_VERSION` and
# `$KUBECTL_VERSION` env variables into `$PREFIX/bin`, add `$PREFIX/bin` to 
# the `$PATH` and delegate further execution to `/bin/drone-helm`
# See `set-environments.sh` for baked in versions.

HELM_VERSION="${HELM_VERSION:-v2.14.1}"
KUBECTL_VERSION="${KUBECTL_VERSION:-v1.14.3}"

mkdir -p ~/.local/bin
ln -s -f ~/.local/lib/helm-${HELM_VERSION}/helm ~/.local/bin
ln -s -f ~/.local/lib/kubectl-${KUBECTL_VERSION}/kubectl ~/.local/bin

export PATH=~/.local/bin:$PATH

echo "Using helm    Version: ${HELM_VERSION} installed into ~/.local/bin"
echo "Using kubectl Version: ${KUBECTL_VERSION} installed into ~/.local/bin"

helm version --client || error "Helm installation is not functional"
kubectl version --client || error "Kubectl installation is not functional"

/bin/drone-helm "$@"
