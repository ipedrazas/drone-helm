FROM alpine:latest
MAINTAINER Ivan Pedrazas <ipedrazas@gmail.com>


RUN apk -Uuv add curl bash && rm /var/cache/apk/*

ENV HELM_VERSION v2.1.0
ENV FILENAME helm-${HELM_VERSION}-linux-amd64.tar.gz
ENV KUBECTL v1.4.6

ADD http://storage.googleapis.com/kubernetes-helm/${FILENAME} /tmp

ADD https://storage.googleapis.com/kubernetes-release/release/${KUBECTL}/bin/linux/amd64/kubectl /tmp


RUN tar -zxvf /tmp/${FILENAME} -C /tmp \
  && mv /tmp/linux-amd64/helm /bin/helm \
  && chmod +x /bin/helm \
  && chmod +x /tmp/kubectl \
  && mv /tmp/kubectl /bin/kubectl \
  && rm -rf /tmp

LABEL description="Drone.io plugin for Helm (kubernetes)."
LABEL base="alpine"
LABEL helm=${FILENAME}
LABEL kubectl=${KUBECTL}

COPY drone-helm /bin/drone-helm
COPY kubeconfig /root/.kube/kubeconfig

ENTRYPOINT [ "/bin/drone-helm" ]