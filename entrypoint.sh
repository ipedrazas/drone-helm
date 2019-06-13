#!/bin/sh

# Will symlink `helm` and `kubectl` binaries based on the `$HELM_VERSION` and
# `$KUBECTL_VERSION` env variables into `$PREFIX/bin`, add `$PREFIX/bin` to 
# the `$PATH` and delegate further execution to `/bin/drone-helm`
# See `set-environments.sh` for baked in versions.

PREFIX=~/.local
HELM_VERSION="${HELM_VERSION:-v2.14.1}"
KUBECTL_VERSION="${KUBECTL_VERSION:-v1.14.3}"

mkdir -p ${PREFIX}/bin
ln -s -f ${PREFIX}/lib/helm-${HELM_VERSION}/helm ${prefix}/bin
ln -s -f ${PREFIX}/lib/kubectl-${KUBECTL_VERSION}/kubectl ${prefix}/bin

export PATH=${PREFIX}/bin:$PATH

echo "Using helm    Version: ${HELM_VERSION}"
echo "Using kubectl Version: ${KUBECTL_VERSION}"

/bin/drone-helm "$@"