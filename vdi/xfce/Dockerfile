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
# TODO: set background in a different way, it's an ugly way
# set background
RUN rm /usr/share/backgrounds/xfce/*
COPY xfce-stripes.png /usr/share/backgrounds/xfce/
RUN cp /usr/share/backgrounds/xfce/xfce-stripes.png /usr/share/backgrounds/xfce/xfce-verticals.png
#
# copy configuation files and generate xorg.conf
COPY generate-xorg.sh /etc/vdi
COPY xorg.conf /etc/vdi
RUN chmod 755 /etc/vdi/generate-xorg.sh
# RUN /etc/vdi/generate-xorg.sh

COPY supervisord.conf /etc/vdi/supervisord/xfce.conf

#
#
# improve usability

# COPY ./SomeTheme.tar.xz /usr/share/themes/

RUN set -eux; \
    apt-get update; \
    apt-get install -y software-properties-common; \
    # install firefox as default browser
    add-apt-repository ppa:mozillateam/ppa; \
    echo 'Package: *\nPin: release o=LP-PPA-mozillateam\nPin-Priority: 1001' | sudo tee /etc/apt/preferences.d/mozilla-firefox; \
    apt-get install -y firefox;
    # install theme
    # rm -rf /usr/share/themes/Default; \
    # tar -xvf /usr/share/themes/SomeTheme.tar.xz; \
    # mv SomeTheme /usr/share/themes/Default

ENV BROWSER=/usr/bin/firefox

CMD ["/bin/bash", "-c", "/etc/vdi/generate-xorg.sh && supervisord -c /etc/vdi/supervisord.conf"]