package web

import (
	"path"
)

//EmbedArchive 归档文件
var EmbedArchive []byte

var EmbedExt string

func init() {
	for _, v := range AssetNames() {
		EmbedExt = path.Ext(v)
		EmbedArchive, _ = Asset(v)
	}
}
