FROM ubuntu:focal

SHELL [ "/bin/bash", "-c" ]

COPY . .

RUN set eux; \
    apt-get update; \
    apt-get install -y ifstat; \
    chmod +x gpu-util.sh network-load.sh

CMD [ "./network-load.sh" ]