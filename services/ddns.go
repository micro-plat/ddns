package services

import (
	"github.com/micro-plat/hydra/component"
	"github.com/micro-plat/hydra/context"
)

type DdnsHandler struct {
	container component.IContainer
}

func NewDdnsHandler(container component.IContainer) (u *DdnsHandler) {
	return &DdnsHandler{container: container}
}

//Handle 保存动态域名信息
func (u *DdnsHandler) Handle(ctx *context.Context) (r interface{}) {
	ctx.Log.Info("--------------保存动态域名信息---------------")

	ctx.Log.Info("1.获取参数")
	domians, err := RemoteChromedp()
	if err != nil {
		return err
	}

	for _, v := range domians {
		ctx.Log.Infof("2. 获取分布式锁,v%+v", v)
		lk := ctx.NewDLock(v.Domain)
		if err := lk.Lock(); err != nil {
			return err
		}
		defer lk.Unlock()

		ctx.Log.Info("3. 检查并创建解析信息")
		if err := checkAndCreate(v, u.container.GetRegistry()); err != nil {
			return err
		}
	}

	return "success"
}
