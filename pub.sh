#!/usr/bin/expect 

# 自动发布，执行命令：expect pub.sh
#安装expect组件请使用命令:sudo apt install expect


set timeout -1  
set uname "yanfa"
set host "192.168.106.190"
set pwd  "-A0l1ao!@##@!\r"
set dt  [exec date "+%Y%m%d%H%M%S"]


#编译文件-------------
spawn go build


#上传文件-------------
spawn echo "上传文件..."
spawn scp ./ddns $uname@$host:/tmp
expect {
    "password" {send $pwd;}
}
expect eof


#远程更新---------------
spawn echo "远程更新..."
spawn ssh -t  $uname@$host "cd /srv/ddns/bin;sudo ./ddns stop;sudo cp ./ddns ./ddns_${dt} ;sudo rm -rf ./ddns;sudo cp /tmp/ddns ./;sleep 3;sudo ./ddns start;"
expect {
    "password" {send $pwd;exp_continue}
    "密码" {send $pwd;}
}
expect eof

#删除本地文件
spawn rm -rf ./ddns
exit