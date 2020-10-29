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

//Handle 保存动态域名信息
func (u *GithubHandler) Handle(ctx hydra.IContext) (r interface{}) {
	ctx.Log().Info("--------------定时获取github动态域名信息---------------")

	ctx.Log().Info("1.获取github domains")
	domians, err := GetGithubDomains()
	if err != nil {
		return err
	}

	ctx.Log().Info("2.保存域名")
	registry, err := registry.NewRegistry(hydra.G.RegistryAddr, ctx.Log())
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
