package conf

import (
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/conf"
)

//TypeNodeName 分类节点名
const TypeNodeName = "dnss"

//Dnss dns服务器dns列表配置信息
type Dnss struct {
	Dnss []string `json:"dnss,omitempty" toml:"dnss,omitempty"`
}

//NewDnss 构建任务列表
func NewDnss(dns ...string) *Dnss {
	if len(dns) > 0 {
		return &Dnss{Dnss: dns}
	}
	return &Dnss{Dnss: []string{"114.114.114.114", "8.8.8.8"}}
}

//GetSubConf .
func GetSubConf(cnf conf.IServerConf) (dnss *Dnss, err error) {
	dnss = &Dnss{Dnss: []string{}}
	_, err = cnf.GetSubObject(TypeNodeName, dnss)
	if err == conf.ErrNoSetting {
		return &Dnss{Dnss: []string{"114.114.114.114", "8.8.8.8"}}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("ddns配置错误:%v", err)
	}

	if b, err := govalidator.ValidateStruct(dnss); !b {
		return nil, fmt.Errorf("ddns配置有误:%v", err)
	}
	return dnss, nil
}
