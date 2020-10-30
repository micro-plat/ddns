package services

import (
	"net/http"

	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/lib4go/errs"
)

// DdnsHandler Handler
type DdnsHandler struct {
}

// NewDdnsHandler 构建DdnsHandler
func NewDdnsHandler() *DdnsHandler {
	return &DdnsHandler{}
}

//RequestHandle 保存动态域名信息
func (u *DdnsHandler) RequestHandle(ctx hydra.IContext) (r interface{}) {
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

//QueryHandle 查询域名信息
func (u *DdnsHandler) QueryHandle(ctx hydra.IContext) (r interface{}) {
	ctx.Log().Info("--------------查询域名信息---------------")

	ctx.Log().Info("1. 检查必须参数")
	var domain string
	if domain := ctx.Request().GetString("d"); domain == "" {
		return errs.NewError(http.StatusNotAcceptable, "域名不能为空")
	}

	ctx.Log().Info("2. 查询解析信息")
	rgst, err := registry.NewRegistry(hydra.G.RegistryAddr, ctx.Log())
	if err != nil {
		return err
	}
	ps, _, err := rgst.GetChildren(registry.Join("/dns", domain))
	if err != nil {
		return err
	}
	return ps
}
