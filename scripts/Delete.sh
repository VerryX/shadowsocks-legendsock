#!/usr/bin/env bash
echo "Stop service"
systemctl stop shadowsocks-legendsock
systemctl disable shadowsocks-legendsock

echo "Delete service"
rm -f /etc/systemd/system/shadowsocks-legendsock.service

echo "Reload Systemd"
systemctl daemon-reload

echo "Delete binaries"
rm -rf /opt/shadowsocks-legendsock

echo "Done"

exit 0
