#
# ------ Drone-Helm build image ------
#
FROM golang:1.12-alpine3.9 as builder

RUN apk update
RUN apk add dep git

ENV GOOS=linux 
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
COPY  *.sh /bin/
COPY kubeconfig /root/.kube/kubeconfig

RUN set -ex \
  && apk add --no-cache curl ca-certificates bash \
  && curl -o /tmp/aws-iam-authenticator https://amazon-eks.s3-us-west-2.amazonaws.com/1.10.3/2018-07-26/bin/linux/amd64/aws-iam-authenticator \
  && chmod +x /tmp/aws-iam-authenticator \
  && mv /tmp/aws-iam-authenticator /bin/aws-iam-authenticator 
RUN /bin/setup-environments.sh
RUN rm -rf /tmp/*

LABEL description="Kubectl and Helm."
LABEL base="alpine"

ENTRYPOINT [ "/bin/entrypoint.sh" ]
