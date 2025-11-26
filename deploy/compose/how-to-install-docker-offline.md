# How to Install Docker and Docker Compose Offline

## Docker

### Download Docker Package

Download from [Docker official website](https://download.docker.com/linux/static/stable/x86_64/), recommended version [20.10.24](https://download.docker.com/linux/static/stable/x86_64/docker-20.10.24.tgz)

### Extract and Copy

```bash
# docker-20.10.24.tgz is the actual location of the downloaded file
tar -zxvf docker-20.10.24.tgz -C /tmp
cp /tmp/docker/* /usr/bin/
```

### Register Docker Service

vim /etc/systemd/system/docker.service

```
[Unit]
Description=Docker Application Container Engine
Documentation=https://docs.docker.com
After=network-online.target firewalld.service
Wants=network-online.target
[Service]
Type=notify
ExecStart=/usr/bin/dockerd --selinux-enabled=false --insecure-registry=127.0.0.1
ExecReload=/bin/kill -s HUP $MAINPID
LimitNOFILE=infinity
LimitNPROC=infinity
LimitCORE=infinity
TimeoutStartSec=0
Delegate=yes
KillMode=process
Restart=on-failure
StartLimitBurst=3
StartLimitInterval=60s
[Install]
WantedBy=multi-user.target
```

### Start

```bash
chmod 777 /etc/systemd/system/docker.service
systemctl daemon-reload
systemctl enable docker
systemctl start docker
```

### Verify

```bash
systemctl status docker
docker -v
docker info
```

## Docker Compose

### Download

Download from [GitHub release](https://github.com/docker/compose/releases), recommended version: [v2.39.3](https://github.com/docker/compose/releases/download/v2.39.3/docker-compose-linux-x86_64)

### Install

```bash
# docker-compose-linux-x86_64 is the actual location of the downloaded file
sudo mv docker-compose-linux-x86_64 /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
sudo ln -s /usr/local/bin/docker-compose /usr/bin/docker-compose
cp -a /usr/local/bin/docker-compose  ~/.docker/cli-plugins/docker-compose # Add to docker subcommands
```

### Verify

```bash
docker-compose 
docker compose
```