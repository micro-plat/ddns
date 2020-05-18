package services

import (
	"github.com/micro-plat/hydra/component"
	"github.com/micro-plat/hydra/context"
)

// GithubHandler handle
type GithubHandler struct {
	container component.IContainer
}

// NewGithubHandler 构建GithubHandler
func NewGithubHandler(container component.IContainer) *GithubHandler {
	return &GithubHandler{container: container}
}

//Handle 保存动态域名信息
func (u *GithubHandler) Handle(ctx *context.Context) (r interface{}) {
	ctx.Log.Info("--------------定时获取github动态域名信息---------------")

	ctx.Log.Info("1.获取github domains")
	domians, err := GetGithubDomains()
	if err != nil {
		return err
	}
	ctx.Log.Infof("2.获取分布式锁,v%+v", domians)
	for _, v := range domians {
		lk := ctx.NewDLock(v.Domain)
		if err := lk.Lock(); err != nil {
			return err
		}
		defer lk.Unlock()
		if err := checkAndCreate(v, u.container.GetRegistry()); err != nil {
			return err
		}
	}

	ctx.Log.Info("3. 返回结果")
	return "success"
}
