# ddns

提供简单、快速的 DNS 缓存服务

√ 　支持 linux,macos

√ 　支持从`hosts`(`/etc/hosts*`,`C:\Windows\System32\drivers\etc\hosts*`)读取配置信息

√ 　支持从`//etc/resolv.conf`,`windows注册表`读取上游 DNS 服务器 IP

√ 　支持从注册中心`/dns`读取配置

√ 　支持`hydra`应用注册的`dns`服务

√ 　缓存上游 DNS 解析结果，加快响应速度

√ 　所有配置热更新，无需重启服务

√ 　上游 DNS 服务器检测，优先使用速度最快的服务器

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
sudo ddns install
```

- 运行

```sh
sudo ddns start
```

- 测试

```sh
dig github.com @127.0.0.1
```

- 本机使用

```sh
sudo vim /etc/resolv.conf

#修改内容如下:
nameserver 127.0.0.1
```

## 2. hosts 文件

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
/dns
-----google.com
--------172.217.6.127
```

## 4.hydra 服务

为API服务设置解析域名
```go
hydra.Conf.API(":8081", api.WithDNS("www.ddns.com"))
```
DDNS实时收到解析信息

```sh
[2019/08/23 16:31:12.48855][i][9077e9414][缓存:www.ddns.com,1条]
```

## 5.上游 DNS

未在`hosts*`,或`注册中心`配置的域名，直接使用`上游DNS服务器`查询域名解析结果

- 打开`/etc/resolv.conf`文件，添加上游 DNS 服务器 IP

```sh
sudo vim /etc/resolv.conf 
```

```sh
nameserver 127.0.0.1
nameserver 114.114.114.114
nameserver 8.8.8.8
```

## 7. 优先级

注册中心 > 本地 HOSTS > 上游 DNS
