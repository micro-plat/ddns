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

//CheckHandle 检查域名信息
func (u *GithubHandler) CheckHandle(ctx hydra.IContext) (r interface{}) {
	ctx.Log().Info("--------------检查github是否可以访问---------------")

	ctx.Log().Info("1.检查域名是否可用")
	if err := Check(); err != nil {
		return err
	}
	return "success"
}

//RequestHandle 保存动态域名信息
func (u *GithubHandler) RequestHandle(ctx hydra.IContext) (r interface{}) {
	ctx.Log().Info("--------------保存github域名解析信息---------------")

	ctx.Log().Info("1.获取github域名信息")
	domians, err := GetGithubDomains()
	if err != nil {
		return err
	}

	ctx.Log().Info("2.保存域名")
	for _, v := range domians {
		if err := local.R.CreateOrUpdateGithub(v.Domain, v.IP, v.Value); err != nil {
			return err
		}
	}
	return "success"
}
