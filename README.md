# shadowsocks-legendsock
**注意，此后端仅支持原版 SS 而且加密方式较少！**

## 特性
- 使用 Go 编写，性能高
- 已经对接好 LegendSock （目前 YoYu 的 Air 在使用此后端）

## 支持的加密方式
```
AEAD_AES_128_GCM (AES-128-GCM)
AEAD_AES_192_GCM (AES-192-GCM)
AEAD_AES_256_GCM (AES-256-GCM)
AEAD_CHACHA20_POLY1305 (CHACHA20-POLY1305)
AEAD_XCHACHA20_POLY1305 (XCHACHA20-POLY1305)
RC4-MD5
AES-128-CFB
AES-192-CFB
AES-256-CFB
AES-128-CTR
AES-192-CTR
AES-256-CTR
CHACHA20
CHACHA20-IETF
XCHACHA20
```

## Ubuntu 自动安装脚本
1. 运行脚本
```bash
curl -fsSL https://raw.githubusercontent.com/shadowsocks-server/shadowsocks-legendsock/master/scripts/Ubuntu.sh | bash
```

2. 编辑配置文件调整参数
```bash
vim /etc/systemd/system/shadowsocks-legendsock.service
```

3. 重启 Systemd 并开启服务设置自启
```bash
systemctl daemon-reload
systemctl start shadowsocks-legendsock
systemctl enable shadowsocks-legendsock
systemctl status shadowsocks-legendsock
```