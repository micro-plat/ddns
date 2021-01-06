cd web/ddnsweb


echo "1. 安装依赖：npm install"
npm  install

echo "2. 打包项目：npm run build"
npm run build

echo "3. 压缩：dist/static"
cd dist/static
zip -q -r ../../../static.zip *

echo "4. 生成资源文件:web/static.go" 
cd ../../../
go-bindata -o=./static.go -pkg=web static.zip
sleep 1s

echo "5. 写入静态文件配置内容到web/web.go" 
echo '
package web

import (
	"path"

	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/conf/server/static"
)

func init() {
	hydra.OnReady(func() {
		for _, v := range AssetNames() {
			ext := path.Ext(v)
			embed, _ := Asset(v)
			hydra.Conf.GetWeb().Static(static.WithArchiveByEmbed(embed, ext))
		}
	})
}
' > ./web.go

echo "5. 删除打包文件和压缩文件" 
rm -rf ddnsweb/dist/
rm -rf static.zip
cd ..


echo "6. 编译项目"
go build  -o out/ddnsserver

echo ""
echo "---------打包-success----------------" 
echo "---------目录:/out"
echo ""