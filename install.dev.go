// +build !prod

package main

import "github.com/micro-plat/hydra/conf"

//bindConf 绑定启动配置， 启动时检查注册中心配置是否存在，不存在则引导用户输入配置参数并自动创建到注册中心
func init() {
	app.IsDebug = true
	app.Conf.API.SetMain(conf.NewAPIServerConf("9090"))
	tasks := conf.NewTasks()
	tasks.Append("@midnight", "/github/ip/check")
	app.Conf.CRON.SetTasks(tasks)
}
