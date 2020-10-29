package services

import (
	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/registry"
)

// DdnsHandler Handler
type DdnsHandler struct {
}

// NewDdnsHandler 构建DdnsHandler
func NewDdnsHandler() *DdnsHandler {
	return &DdnsHandler{}
}

//Handle 保存动态域名信息
func (u *DdnsHandler) Handle(ctx hydra.IContext) (r interface{}) {
	ctx.Log().Info("--------------保存动态域名信息---------------")

	ctx.Log().Info("1. 检查必须参数")
	var domain Domain
	if err := ctx.Request().Bind(&domain); err != nil {
		return err
	}

	ctx.Log().Info("2. 检查并创建解析信息")
	registry, err := registry.NewRegistry(hydra.G.RegistryAddr, ctx.Log())
	if err != nil {
		return err
	}
	return checkAndCreate(&domain, registry)
}
