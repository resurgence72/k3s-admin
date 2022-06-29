package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type router struct{}

var Router router

// 初始化路由规则，创建测试api接口
func (r *router) InitApiRouter(router *gin.Engine) {
	router.GET("/testapi", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "testapi success",
			"data": nil,
		})
	})

	// 获取pod列表
	router.GET("/api/v1/pods", Pod.GetPods)
	// 获取pod详情
	router.GET("/api/v1/pod", Pod.GetPodDetail)
	// 删除pod
	router.DELETE("/api/v1/pod", Pod.DeletePod)
	// 更新pod
	router.PATCH("/api/v1/pod", Pod.UpdatePod)

	// 获取pod日志
	router.GET("/api/v1/pod-log", Pod.PodLog)
	// 获取pod实时日志
	router.GET("/api/v1/pod-log-sync", Pod.PodLogSync)
}
