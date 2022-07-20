package main

import (
	"github.com/gin-gonic/gin"
	"k3s-admin/config"
	"k3s-admin/controller"
	"k3s-admin/service"
	"net/http"
	"os"
)

func main() {
	// 初始化数据库
	//db.Init()

	// 初始化k8s-client
	err := service.K8s.Init()
	if err != nil {
		os.Exit(-1)
	}

	r := gin.Default()
	// 初始化router方法
	controller.Router.InitApiRouter(r)

	// 启动ws
	go func() {
		http.HandleFunc("/ws", service.Terminal.WsHandler)
		http.ListenAndServe(":8081", nil)
	}()

	// 启动gin server
	_ = r.Run(config.ListenAddr)

	// 关闭数据库连接
	//db.Close()
}
