echo "start dig"
for k in $( seq 1 1000000 )
do
    echo "dig :$k"
    dig @192.168.5.94 github.com
    dig @192.168.5.94 baidu.com
    dig @192.168.5.94 zhuanlan.zhihu.com
    dig @192.168.5.94 github.global.ssl.fastly.net
    dig @192.168.5.94 cnblogs.com
    dig @192.168.5.94 studygolang.com
    dig @192.168.5.94 fanyi.baidu.com
    dig @192.168.5.94 jianshu.com
    dig @192.168.5.94 developer.github.com
    dig @192.168.5.94 api.coupon.17ebs.18jiayou0.com
    dig @192.168.5.94 api.sso.18jiayou0.com
    dig @192.168.5.94 api.sso.18jiayou1.com
    dig @192.168.5.94 goproxy.cn
    dig @192.168.5.94 webapi.sso.18jiayou1.com
    dig @192.168.5.94 api.bss.sso.18jiayou0.com
    dig @192.168.5.94 www.sina.com.cn
    dig @192.168.5.94 www.qq.com
    dig @192.168.5.94 www.163.com
    dig @192.168.5.94 www.ctrip.com
    dig @192.168.5.94 www.jd.com
    dig @192.168.5.94 cd.58.com
done