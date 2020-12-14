package services

import (
	"fmt"

	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/registry"
)

// GithubHandler handle
type GithubHandler struct {
}

// NewGithubHandler 构建GithubHandler
func NewGithubHandler() *GithubHandler {
	return &GithubHandler{}
}

//CheckHandle 检查域名信息
func (u *GithubHandler) CheckHandle(ctx hydra.IContext) (r interface{}) {
	ctx.Log().Info("--------------检查github是否可以访问---------------")

	ctx.Log().Info("1.检查域名是否可用")
	err := Check()
	if err != nil {
		return err
	}
	return "success"
}

//RequestHandle 保存动态域名信息
func (u *GithubHandler) RequestHandle(ctx hydra.IContext) (r interface{}) {
	ctx.Log().Info("--------------定时获取github动态域名信息---------------")

	ctx.Log().Info("1.获取github domains")
	domians, err := GetGithubDomains()
	if err != nil {
		return err
	}

	ctx.Log().Info("2.保存域名")
	registry, err := registry.GetRegistry(hydra.G.RegistryAddr, ctx.Log())
	if err != nil {
		return fmt.Errorf("无法保存域名:%w", err)
	}

	for _, v := range domians {
		if err := checkAndCreate(v, registry); err != nil {
			return err
		}
	}

	ctx.Log().Info("3. 返回结果")
	return "success"
}
