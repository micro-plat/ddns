package services

import (
	"github.com/micro-plat/ddns/local"
	"github.com/micro-plat/hydra"
)

// GithubHandler handle
type GithubHandler struct {
}

// NewGithubHandler 构建GithubHandler
func NewGithubHandler() *GithubHandler {
	return &GithubHandler{}
}

//RequestHandle 保存动态域名信息
func (u *GithubHandler) RequestHandle(ctx hydra.IContext) (r interface{}) {
	ctx.Log().Info("--------------github域名解析---------------")

	ctx.Log().Info("1.获取github最快的IP信息")
	domians, err := GetGithubDomains()
	if err != nil {
		return err
	}

	ctx.Log().Info("2.保存域名")
	for _, v := range domians {
		ctx.Log().Infof("保存%s %s", v.Domain, v.IP)
		if err := local.R.CreateOrUpdate(v.Domain, v.IP, true, v.Value); err != nil {
			return err
		}
	}
	return "success"
}
