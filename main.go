package main

import (
	"github.com/gin-gonic/gin"
	"k3s-admin/config"
	"k3s-admin/controller"
	"k3s-admin/service"
	"os"
)

func main() {
	// 初始化k8s-client
	err := service.K8s.Init()
	if err != nil {
		os.Exit(-1)
	}


	r := gin.Default()
	// 初始化router方法
	controller.Router.InitApiRouter(r)

	// 启动gin server
	_ = r.Run(config.ListenAddr)
}
