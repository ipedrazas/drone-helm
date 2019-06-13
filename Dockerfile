#
# ------ Drone-Helm build image ------
#
FROM golang:1.12-alpine3.9 as builder

RUN apk update
RUN apk add dep git

ENV GOOS linux 
ENV GOARCH=386 

WORKDIR /go/src/github.com/ipedrazas/drone-helm
COPY . .

RUN dep ensure
RUN go build

#
# ------ Drone-Helm plugin image ------
#
FROM alpine:3.9 as final
MAINTAINER Ivan Pedrazas <ipedrazas@gmail.com>

COPY --from=builder /go/src/github.com/ipedrazas/drone-helm/drone-helm /bin/

# Helm version: can be passed at build time
ARG VERSION
ENV VERSION ${VERSION:-v2.14.1}
ENV FILENAME helm-${VERSION}-linux-amd64.tar.gz

ARG KUBECTL
ENV KUBECTL ${KUBECTL:-v1.14.3}

RUN set -ex \
  && apk add --no-cache curl ca-certificates \
  && curl -o /tmp/${FILENAME} http://storage.googleapis.com/kubernetes-helm/${FILENAME} \
  && curl -o /tmp/kubectl https://storage.googleapis.com/kubernetes-release/release/${KUBECTL}/bin/linux/amd64/kubectl \
  && curl -o /tmp/aws-iam-authenticator https://amazon-eks.s3-us-west-2.amazonaws.com/1.10.3/2018-07-26/bin/linux/amd64/aws-iam-authenticator \
  && tar -zxvf /tmp/${FILENAME} -C /tmp \
  && mv /tmp/linux-amd64/helm /bin/helm \
  && chmod +x /tmp/kubectl \
  && mv /tmp/kubectl /bin/kubectl \
  && chmod +x /tmp/aws-iam-authenticator \
  && mv /tmp/aws-iam-authenticator /bin/aws-iam-authenticator \
  && rm -rf /tmp/*

LABEL description="Kubectl and Helm."
LABEL base="alpine"
COPY kubeconfig /root/.kube/kubeconfig

ENTRYPOINT [ "/bin/drone-helm" ]
