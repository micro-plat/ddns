# ddns

提供简单、快速的 DNS 缓存服务

√ 　支持 windows,linux,macos

√ 　支持从`hosts`(`/etc/hosts*`,`C:\Windows\System32\drivers\etc\hosts*`)读取配置信息

√ 　支持从`/etc/names.conf`,`C:\Windows\System32\drivers\etc\names.conf`读取上游 DNS 服务器 IP

√ 　支持从注册中心`/dns`读取配置

√ 　支持`hydra`应用注册的`dns`服务，并立即生效

√ 　缓存上游 DNS 解析结果，加快响应速度

√ 　所有配置热更新，无需重启服务

√ 　上游 DNS 服务器检测，优先使用速度最快的服务器

√ 　解决 windows 只能配置 2 个 DNS 服务器地址问题

√ 　基于[hydra](https://github.com/micro-plat/hydra)实现

## 1. 快速使用

- 下载

```sh
go get github.com/micro-plat/ddns
```

- 编译

```sh
go install github.com/micro-plat/ddns
```

- 安装

```sh
sudo ddns install -r fs://../
```

- 运行

```sh
sudo ddns start
```

- 测试

```sh
dig github.com @ip
```

- 本机使用

```sh
sudo vim /etc/resolv.conf

#修改内容如下:
nameserver 192.168.4.121
```

## 2. hosts 文件

所有将`ddns`作为`dns`服务器的用户，可直接使用`ddns`配置的`hosts`解析信息

- 修改`hosts`文件，添加需解析的域名,删除无需解析的域名
- `etc`目录下新建名称以`hosts`开头的文件，添加解析信息

```sh
sudo vim /etc/hosts_google
```

```sh
# Google Start
172.217.6.127	com.google
172.217.6.127	domains.google
172.217.6.127	environment.google
172.217.6.127	google.com
172.217.6.127	google.com.af
172.217.6.127	google.com.ag
172.217.6.127	google.com.ai
172.217.6.127	google.com.ar
172.217.6.127	google.com.au

"/etc/hosts_google" 9L
```

## 3.注册中心

- 进入注册中心(`fs`或`zookeeper`),在节点`/dns`目录下新建域名，和解析的 IP，如:

```sh
dns
-----google.com
--------172.217.6.127
```

## 4.hydra 服务

本地 IP 作为解析 IP:

```go
app.Conf.API.SetMain(conf.NewAPIServerConf(":8098").WithDNS("api.hydra.com"))
```

使用`LVS`或`nginx`IP 作为解析 IP:

```go
app.Conf.API.SetMain(conf.NewAPIServerConf(":8098").WithDNS("api.hydra.com","172.16.9.100"))
```

## 5. 接口提交

通过 ADSL 拨号的网络，每次拨号后的公网 IP 不同，可通过接口提交到 DDNS 服务器

```sh
curl  "http://127.0.0.1:9090/ddns/request?domain=api.bac.com&ip=192.168.4.121"
```

```sh
[2019/08/23 16:31:11.997083][i][6ad7395c4]api.request GET /ddns/request?domain=api.bac.com&ip=192.168.4.121 from 127.0.0.1
[2019/08/23 16:31:11.997521][i][6ad7395c4]--------------保存动态域名信息---------------
[2019/08/23 16:31:11.997528][i][6ad7395c4]1. 检查必须参数
[2019/08/23 16:31:11.997667][i][6ad7395c4]2. 获取分布式锁
[2019/08/23 16:31:12.16919][i][6ad7395c4]3. 检查并创建解析信息
[2019/08/23 16:31:12.56629][i][6ad7395c4]api.response GET /ddns/request?domain=api.bac.com&ip=192.168.4.121 200  59.578447ms
```

DDNS 服务器会实时收到解析信息

```sh
[2019/08/23 16:31:12.48855][i][9077e9414][缓存:api.bac.com,1条]
```

## 6.上游 DNS

未在`hosts*`,或`注册中心`配置的域名，直接使用`上游DNS服务器`查询域名解析结果

- 打开`etc/names.conf`文件，添加上游 DNS 服务器 IP

```sh
sudo vim /etc/names.conf
```

```sh
119.6.6.6
61.139.2.69
114.114.114.114
180.76.76.76
8.8.8.8

"/etc/names.conf" 2L
```

## 7. 优先级

注册中心 > 本地 HOSTS > 上游 DNS
