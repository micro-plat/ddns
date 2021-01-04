
if [ -f web/static.go ] ; then 
  echo "1. 编译项目"
  go build
	echo  web/static.go"静态二进制文件已存在,不进行打包"
	echo ""
	echo ""
	exit 1 
fi 

cd web/ddnsweb

echo "1. 打包项目：npm run build"
npm run build

echo "2. 压缩：dist/static"
cd dist/static
zip -q -r ../../../static.zip *

echo "3. 生成资源文件:web/static.go" 
cd ../../../
go-bindata -o=./static.go -pkg=web static.zip
sleep 1s

echo "4. 删除打包文件和压缩文件" 
rm -rf ddnsweb/dist/
rm -rf static.zip
cd ..

echo "5. 编译项目"
go build

echo "6. 完成"
exit