package conf

import (
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/conf"
)

//TypeNodeName 分类节点名
const TypeNodeName = "names"

//Names dns服务器dns列表配置信息
type Names struct {
	IPS []string `json:"ips,omitempty" toml:"ips,omitempty"`
}

//NewNames 构建任务列表
func NewNames(ips ...string) *Names {
	if len(ips) > 0 {
		return &Names{IPS: ips}
	}
	return &Names{}
}

//GetNamesConf .
func GetNamesConf(cnf conf.IServerConf) (names *Names, err error) {
	names = &Names{IPS: []string{}}
	_, err = cnf.GetSubObject(TypeNodeName, names)
	if err == conf.ErrNoSetting {
		return &Names{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("ddns配置错误:%v", err)
	}

	if b, err := govalidator.ValidateStruct(names); !b {
		return nil, fmt.Errorf("ddns配置有误:%v", err)
	}
	return names, nil
}
