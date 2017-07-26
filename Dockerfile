#
# ----- Go Builder Image ------
#
FROM golang:1.8-alpine AS builder

RUN apk add --no-cache git

# set working directory
RUN mkdir -p /go/src/drone-helm
WORKDIR /go/src/drone-helm

# copy sources
COPY . .

# run tests
RUN go test -v

# build binary
RUN go build -v -o "/drone-helm"

#
# ------ Drone-Helm plugin image ------
#

FROM alpine:3.6
MAINTAINER Ivan Pedrazas <ipedrazas@gmail.com>

# Helm version: can be passed at build time (default to v2.5.0)
ARG VERSION
ENV VERSION ${VERSION:-v2.5.0}
ENV FILENAME helm-${VERSION}-linux-amd64.tar.gz

ARG KUBECTL
ENV KUBECTL ${KUBECTL:-v1.6.6}

ADD http://storage.googleapis.com/kubernetes-helm/${FILENAME} /tmp

ADD https://storage.googleapis.com/kubernetes-release/release/${KUBECTL}/bin/linux/amd64/kubectl /tmp

RUN apk add --no-cache ca-certificates

RUN tar -zxvf /tmp/${FILENAME} -C /tmp \
  && mv /tmp/linux-amd64/helm /bin/helm \
  && chmod +x /tmp/kubectl \
  && mv /tmp/kubectl /bin/kubectl \
  && rm -rf /tmp

LABEL description="Kubeclt and Helm."
LABEL base="alpine"

COPY --from=builder /drone-helm /bin/drone-helm
COPY kubeconfig /root/.kube/kubeconfig

ENTRYPOINT [ "/bin/drone-helm" ]
