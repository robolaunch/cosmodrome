#!/bin/bash

sudo /usr/bin/dbus &
export DISPLAY=:0
/usr/bin/pulseaudio --disallow-module-loading -vvvv --disallow-exit --exit-idle-time=-1 &
cd /etc/X11
/usr/bin/X vt7 -novtswitch -sharevts +extension "MIT-SHM" :0 -config xorg.conf :0 &
#mate-session &
#/usr/bin/neko serve -d --static "/var/www" --bind "0.0.0.0:$NEKO_BIND"   --epr "$NEKO_UDP_PORT" --display ":0" --h264 --max_fps "30" &
echo "Session Running. Press [Return] to exit. TEST CM"
read