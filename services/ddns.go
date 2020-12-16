package services

import (
	"encoding/json"
	"net/http"

	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/lib4go/errs"
	"github.com/micro-plat/lib4go/types"
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
	registry, err := registry.GetRegistry(hydra.G.RegistryAddr, ctx.Log())
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
	rgst, err := registry.GetRegistry(hydra.G.RegistryAddr, ctx.Log())
	if err != nil {
		return err
	}
	ps, _, err := rgst.GetChildren(registry.Join("/dns", domain))
	if err != nil {
		return err
	}
	return ps
}

//PlatNamesHandle 查询平台名及对应的域名信息
func (u *DdnsHandler) PlatNamesHandle(ctx hydra.IContext) (r interface{}) {
	ctx.Log().Info("--------------查询平台名及对应的域名信息---------------")

	ctx.Log().Info("1. 获取注册中心")
	rgst, err := registry.GetRegistry(hydra.G.RegistryAddr, ctx.Log())
	if err != nil {
		return err
	}

	ctx.Log().Info("2. 获取域名节点")
	domains, _, err := rgst.GetChildren("/dns")
	if err != nil {
		return err
	}

	ctx.Log().Info("3. 处理域名")
	result := make(map[string][]string, 0)
	for _, domain := range domains {
		val, _, err := rgst.GetValue(registry.Join("/dns", domain))
		if err != nil {
			return err
		}
		value := make(types.XMap, 0)
		err = json.Unmarshal(val, &value)
		if err != nil {
			ctx.Log().Infof("%s处理:%v", domain, err)
			continue
		}
		key := types.GetString(value.GetString("cn_plat_name"), domain)
		result[key] = append(result[key], domain)
	}

	return result
}
