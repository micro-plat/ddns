package services

import (
	"github.com/micro-plat/hydra/component"
	"github.com/micro-plat/hydra/context"
)

// DdnsHandler Handler
type DdnsHandler struct {
	container component.IContainer
}

// NewDdnsHandler 构建DdnsHandler
func NewDdnsHandler(container component.IContainer) *DdnsHandler {
	return &DdnsHandler{container: container}
}

//Handle 保存动态域名信息
func (u *DdnsHandler) Handle(ctx *context.Context) (r interface{}) {
	ctx.Log.Info("--------------保存动态域名信息---------------")

	ctx.Log.Info("1. 检查必须参数")
	var domain Domain
	if err := ctx.Request.Bind(&domain); err != nil {
		return context.NewError(context.ERR_NOT_ACCEPTABLE, err)
	}

	ctx.Log.Info("2. 获取分布式锁")
	lk := ctx.NewDLock(domain.Domain)
	if err := lk.Lock(); err != nil {
		return err
	}
	defer lk.Unlock()

	ctx.Log.Info("3. 检查并创建解析信息")
	return checkAndCreate(&domain, u.container.GetRegistry())
}
