#!/bin/bash

set -e

ln -snf "/usr/share/zoneinfo/$TZ" /etc/localtime && echo "$TZ" | tee /etc/timezone > /dev/null
ln -snf /dev/ptmx /dev/tty7

if grep -Fxq "allowed_users=console" /etc/X11/Xwrapper.config; then
  sed -i "s/allowed_users=console/allowed_users=anybody/;$ a needs_root_rights=yes" /etc/X11/Xwrapper.config
fi

if [ "$NVIDIA_VISIBLE_DEVICES" == "all" ]; then
  export GPU_SELECT=$(nvidia-smi --query-gpu=uuid --format=csv | sed -n 2p)
elif [ -z "$NVIDIA_VISIBLE_DEVICES" ]; then
  export GPU_SELECT=$(nvidia-smi --query-gpu=uuid --format=csv | sed -n 2p)
else
  export GPU_SELECT=$(nvidia-smi --id=$(echo "$NVIDIA_VISIBLE_DEVICES" | cut -d ',' -f1) --query-gpu=uuid --format=csv | sed -n 2p)
  if [ -z "$GPU_SELECT" ]; then
    export GPU_SELECT=$(nvidia-smi --query-gpu=uuid --format=csv | sed -n 2p)
  fi
fi

if [ -z "$GPU_SELECT" ]; then
  echo "No NVIDIA GPUs detected."
  apt-get install -y xserver-xorg-video-dummy
  cp /etc/vdi/xorg.conf /etc/X11/xorg.conf
  sed -i 's|::X-ENTRYPOINT::|/usr/bin/X -config /etc/X11/xorg.conf %(ENV_DISPLAY)s|' /etc/vdi/supervisord.conf
  sed -i 's|::NEKO_ENTRYPOINT::|/usr/bin/robolaunch-vdi serve -d --static "/var/www" --bind "0.0.0.0:%(ENV_NEKO_BIND)s"   --epr "%(ENV_NEKO_UDP_PORT)s" --display "%(ENV_DISPLAY)s" --h264|' /etc/vdi/supervisord/xfce.conf
  exit 0
else
  sed -i 's|::X-ENTRYPOINT::|/usr/bin/X vt7 -novtswitch -sharevts +extension "MIT-SHM" %(ENV_DISPLAY)s -config /etc/X11/xorg.conf %(ENV_DISPLAY)s|' /etc/vdi/supervisord.conf
  sed -i 's|::NEKO_ENTRYPOINT::|/usr/bin/robolaunch-vdi serve -d --static "/var/www" --bind "0.0.0.0:%(ENV_NEKO_BIND)s"   --epr "%(ENV_NEKO_UDP_PORT)s" --display "%(ENV_DISPLAY)s" --h264 --hwenc "NVENC" --rc_mode "1"|' /etc/vdi/supervisord/xfce.conf
fi

if [ -f "/etc/X11/xorg.conf" ]; then
    rm /etc/X11/xorg.conf
fi

HEX_ID=$(nvidia-smi --query-gpu=pci.bus_id --id="$GPU_SELECT" --format=csv | sed -n 2p)
IFS=":." ARR_ID=($HEX_ID)
unset IFS
BUS_ID=PCI:$((16#${ARR_ID[1]})):$((16#${ARR_ID[2]})):$((16#${ARR_ID[3]}))
export MODELINE=$(cvt -r 2048 1152 | sed -n 2p)
nvidia-xconfig --depth="$CDEPTH" --mode=$(echo $MODELINE | awk '{print $2}' | tr -d '"') --allow-empty-initial-configuration --no-probe-all-gpus --busid="$BUS_ID" --only-one-x-screen --connected-monitor="$VIDEO_PORT"
sed -i '/Driver\s\+"nvidia"/a\    Option         "ModeValidation" "NoMaxPClkCheck, NoEdidMaxPClkCheck, NoMaxSizeCheck, NoHorizSyncCheck, NoVertRefreshCheck, NoVirtualSizeCheck, NoExtendedGpuCapabilitiesCheck, NoTotalSizeCheck, NoDualLinkDVICheck, NoDisplayPortBandwidthCheck, AllowNon3DVisionModes, AllowNonHDMI3DModes, AllowNonEdidModes, NoEdidHDMI2Check, AllowDpInterlaced"\n    Option         "DPI" "96 x 96"' /etc/X11/xorg.conf
sed -i '/Section\s\+"Monitor"/a\    '"$MODELINE" /etc/X11/xorg.conf