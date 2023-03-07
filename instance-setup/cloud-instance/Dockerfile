# docker run -it \
#     -v $KUBECONFIG:$KUBECONFIG \
#     -e KUBECONFIG \
#     -e ORGANIZATON \
#     -e TEAM \
#     -e REGION \
#     -e CLOUD_INSTANCE \
#     -e CLOUD_INSTANCE_ALIAS \
#     -e PUBLIC_IP \
#     robolaunchio/autorun:v1

# docker run -it \
#     -v $KUBECONFIG:$KUBECONFIG \
#     -e KUBECONFIG \
#     --env-file=.env \
#     --gpus all \
#     robolaunchio/autorun:v1

FROM ubuntu:jammy
ARG TARGETARCH=amd64
ARG DEBIAN_FRONTEND=noninteractive
SHELL [ "/bin/bash", "-c" ]
COPY run.sh /root/

RUN apt-get update && apt-get install -y curl wget figlet
RUN set -eux; \
    # install yq
    wget -qO /usr/local/bin/yq https://github.com/mikefarah/yq/releases/latest/download/yq_linux_${TARGETARCH}; \
    chmod a+x /usr/local/bin/yq; \
    # install kubectl
    curl -LO https://dl.k8s.io/release/v1.24.10/bin/linux/${TARGETARCH}/kubectl; \
    install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl; \
    # install helm
    curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash;

# RUN set -eux; \
#     apt-get update && apt-get install -y gnupg linux-headers-$(uname -r) && apt-get install -y --no-install-recommends nvidia-driver-470; \
#     distribution=$(. /etc/os-release;echo $ID$VERSION_ID) && curl -s -L https://nvidia.github.io/nvidia-docker/gpgkey | apt-key add - && curl -s -L https://nvidia.github.io/nvidia-docker/$distribution/nvidia-docker.list | tee /etc/apt/sources.list.d/nvidia-docker.list; \
#     apt-get update && apt-get install -y nvidia-container-runtime;

RUN chmod +x /root/run.sh
CMD [ "/root/run.sh" ]