FROM ipedrazas/k8s:latest
MAINTAINER Ivan Pedrazas <ipedrazas@gmail.com>

COPY drone-helm /bin/drone-helm
COPY kubeconfig /root/.kube/kubeconfig

ENTRYPOINT [ "/bin/drone-helm" ]
