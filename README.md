# ddns
提供简单、快速的DNS缓存服务器


√　支持从`/etc/hosts`中读取配置

√　支持从`注册中心` `/dns`节点下读取配置

√　支持从`/etc/names.conf`读取上游DNS服务器

√　动态更新`/etc/hosts`、`/etc/names.conf`、`/dns`自动加载,不需要重启服务

√　缓存所有解析结果，缩短响应时长

√　基于[hydra](https://github.com/micro-plat/hydra)实现

#### 适用场景：

* windows只能配置两个DNS服务器限制
* 局域网内部通信时使用自定义域名


