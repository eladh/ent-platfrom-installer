FROM sickp/alpine-sshd:latest

ARG K8S_VERSION=v1.12.9
ARG HELM_VERSION=v2.13.0

ENV HELM_FILENAME=helm-${HELM_VERSION}-linux-amd64.tar.gz

# System
RUN apk update && apk add --no-cache --virtual .build-deps \
    ca-certificates \
    apache2-utils \
    curl \
    tar \
    bash \
    openssl \
    python \
    git

#Install gcloud sdk
RUN curl -sSL https://sdk.cloud.google.com | bash
ENV PATH $PATH:/root/google-cloud-sdk/bin
RUN gcloud components update && gcloud components install beta --quiet
RUN ln -sf /root/google-cloud-sdk/bin/gcloud /usr/bin/gcloud

#Install Helm + kubectl
RUN apk add --update ca-certificates \
 && apk add --update -t deps curl  \
 && apk add --update gettext tar gzip \
 && curl -L https://storage.googleapis.com/kubernetes-release/release/${K8S_VERSION}/bin/linux/amd64/kubectl -o /usr/local/bin/kubectl \
 && curl -L https://storage.googleapis.com/kubernetes-helm/${HELM_FILENAME} | tar xz && mv linux-amd64/helm /bin/helm && rm -rf linux-amd64 \
 && chmod +x /usr/local/bin/kubectl \
 && apk del --purge deps \
 && rm /var/cache/apk/*

RUN apk update \
    && apk add --no-cache \
        sudo nano screen bash rsync rdiff-backup \
    && rm -rf /tmp/* /var/tmp/*

# add openssh and clean
RUN apk add --no-cache openssh \
  && sed -i s/#PermitRootLogin.*/PermitRootLogin\ yes/ /etc/ssh/sshd_config \
  && echo "root:root" | chpasswd

WORKDIR /

# add entrypoint script
ADD ./docker-entrypoint.sh /usr/local/bin

#expose SSH Service
EXPOSE 22

ENTRYPOINT [ "bash", "docker-entrypoint.sh" ]