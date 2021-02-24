#!/bin/sh

rootdir=$(pwd)
asset_name=static.zip

pkg=$1

cd $rootdir/web/ddnsweb

echo "1. 安装依赖：npm install"
npm  install

echo "2. 打包项目：npm run build"
npm run build

echo "3. 压缩：dist/static"
cd  $rootdir/web/ddnsweb/dist/static
rm -f $rootdir/$asset_name

zip -q -r $rootdir/$asset_name *


echo "5. 写入静态文件配置内容到assets.web.go"  

rm -f $rootdir/assets.web.go

echo "
//+build !none

package main

import (
	_ \"embed\"

	\"github.com/micro-plat/hydra\"
	\"github.com/micro-plat/hydra/conf/server/header\"
	\"github.com/micro-plat/hydra/conf/server/static\"
	\"github.com/micro-plat/hydra/global\"

)

//go:embed ${asset_name}
var archiveBytes []byte
var archiveName = \"${asset_name}\"
func init() {
	hydra.OnReady(func() {
		staticOpts:= []static.Option{}
		staticOpts = append(staticOpts,static.WithAutoRewrite(), static.WithEmbedBytes(archiveName, archiveBytes))
		serverStatic := hydra.Conf.GetWeb().Static(staticOpts...)
		if global.Def.IsDebug() {
			serverStatic.Header(header.WithCrossDomain())
		}
	})
}
" > $rootdir/assets.web.go 

sleep 1 


echo "5. 删除打包文件和压缩文件" 
rm -rf $rootdir/web/ddnsweb/dist/

echo "6. 编译项目"

buildtags=" -tags=none "
if [ "$pkg" != "none" ] ; then 
	buildtags=""
fi

mkdir -p $rootdir/out

cd $rootdir

echo "go build $buildtags -o $rootdir/out/ddns"
go build -mod=mod $buildtags -o $rootdir/out/ddns


#rm -rf $rootdir/$asset_name
#rm -rf $rootdir/assets.web.go

echo ""
echo "---------打包-success----------------" 
echo "---------目录:/out"
echo ""