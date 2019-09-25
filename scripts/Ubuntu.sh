#!/usr/bin/env bash
echo "Get Go 1.13"
cd /usr/local || exit 1
wget -O go.tar.gz https://dl.google.com/go/go1.13.linux-amd64.tar.gz

echo "Unzip Go 1.13"
tar -zxf go.tar.gz

echo "Delete the Go 1.13 zip file"
rm go.tar.gz

echo "Setting environment variables"
OLDPATH="$PATH"
PATH="$PATH:/usr/local/go/bin"

echo "Get source code"
go get -d -u -v github.com/shadowsocks-server/shadowsocks-legendsock
cd ~/go/src/github.com/shadowsocks-server/shadowsocks-legendsock

echo "Compiling"
go build -ldflags "-w -s" -o shadowsocks-legendsock

echo "Create directories"
mkdir -p /opt/shadowsocks-legendsock

echo "Copy the compiled file"
cp shadowsocks-legendsock /opt/shadowsocks-legendsock

echo "Change permissions"
chmod +x /opt/shadowsocks-legendsock/shadowsocks-legendsock

echo "Configuration service"
cp scripts/systemd/shadowsocks-legendsock.service /etc/systemd/system

echo "Reload Systemd"
systemctl daemon-reload

echo "Delete source code"
rm -rf ~/go

echo "Restore environment variables"
PATH="$OLDPATH"

echo "Delete Go 1.13"
cd /usr/local
rm -rf go

echo "Done"
echo
echo "Next, you need to do it manually"
echo "1. Change '/etc/systemd/system/shadowsocks-legendsock.service'"
echo "2. Execute 'systemctl daemon-reload'"
echo "3. Execute 'systemctl enable shadowsocks-legendsock'"
echo "4. Execute 'systemctl restart shadowsocks-legendsock'"
echo "5. If there is no error, then you have configured it."

exit 0
