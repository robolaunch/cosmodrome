ARG BASE_IMAGE=robolaunchio/vdi:base
FROM $BASE_IMAGE

ENV DEBIAN_FRONTEND=noninteractive
ENV TZ UTC

#
# install xfce
RUN set -eux; apt-get update; \
    apt-get install -y --no-install-recommends xserver-xorg xorg-dev xfce4 xfce4-terminal sudo; \
    # clean up
    apt-get clean -y; \
    rm -rf /var/lib/apt/lists/* /var/cache/apt/*

#
#
# set background
RUN rm /usr/share/backgrounds/xfce/xfce-stripes.png
COPY xfce-stripes.png /usr/share/backgrounds/xfce/

#
# copy configuation files and generate xorg.conf
COPY generate-xorg.sh /etc/vdi
COPY xorg.conf /etc/vdi
RUN chmod 755 /etc/vdi/generate-xorg.sh
# RUN /etc/vdi/generate-xorg.sh

COPY supervisord.conf /etc/vdi/supervisord/xfce.conf

CMD ["/bin/bash", "-c", "/etc/vdi/generate-xorg.sh && supervisord -c /etc/vdi/supervisord.conf"]