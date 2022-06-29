package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"k3s-admin/service"
	"net/http"
)

var Pod pod

type pod struct{}

// 获取pod列表 支持分页 过滤 排序
func (p *pod) GetPods(c *gin.Context) {
	// get 请求
	req := &struct {
		FilterName string `form:"filter"`
		Namespace  string `form:"namespace"`
		Limit      int    `form:"limit"`
		Page       int    `form:"page"`
	}{}

	err := c.BindQuery(req)
	if err != nil {
		fmt.Println("getpods failed ", err)
		c.JSON(http.StatusOK, gin.H{
			"msg":  "参数绑定失败" + err.Error(),
			"data": nil,
		})
		return
	}

	resp, err := service.Pod.GetPods(req.FilterName, req.Namespace, req.Limit, req.Page)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "get pods failed" + err.Error(),
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":  "success",
		"data": resp,
	})
}

// 获取pod详情
func (p pod) GetPodDetail(c *gin.Context) {
	req := &struct {
		Name      string `form:"name"`
		Namespace string `form:"namespace"`
	}{}
	err := c.BindQuery(req)
	if err != nil {
		fmt.Println("get pod bind query failed ", err)
		c.JSON(http.StatusOK, gin.H{
			"msg":  "get pod bind query failed" + err.Error(),
			"data": nil,
		})
		return
	}

	pod, err := service.Pod.GetPodDetail(req.Name, req.Namespace)
	if err != nil {
		fmt.Println("get pod failed ", err)
		c.JSON(http.StatusOK, gin.H{
			"msg":  "get pod failed" + err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "get pod detail success",
		"data": pod,
	})
}

// 删除pod
func (p pod) DeletePod(c *gin.Context) {
	req := &struct {
		Name      string `form:"name"`
		Namespace string `form:"namespace"`
	}{}
	err := c.BindQuery(req)
	if err != nil {
		fmt.Println("delete pod bind query failed ", err)
		c.JSON(http.StatusOK, gin.H{
			"msg":  "delete pod bind query failed" + err.Error(),
			"data": nil,
		})
		return
	}

	err = service.Pod.DeletePod(req.Name, req.Namespace)
	if err != nil {
		fmt.Println("delete pod failed ", err)
		c.JSON(http.StatusOK, gin.H{
			"msg":  "delete pod failed" + err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "delete pod success",
		"data": nil,
	})
}

func (p pod) UpdatePod(c *gin.Context) {
	req := &struct {
		Name      string `form:"name"`
		Namespace string `form:"namespace"`
		Content   string `form:"content"`
	}{}
	err := c.BindQuery(req)
	if err != nil {
		fmt.Println("update pod bind query failed ", err)
		c.JSON(http.StatusOK, gin.H{
			"msg":  "update pod bind query failed" + err.Error(),
			"data": nil,
		})
		return
	}

	err = service.Pod.UpdatePod(req.Name, req.Namespace, req.Content)
	if err != nil {
		fmt.Println("update pod failed ", err)
		c.JSON(http.StatusOK, gin.H{
			"msg":  "update pod failed" + err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "update pod success",
		"data": nil,
	})
}

// 获取pod日志
func (p pod) PodLog(c *gin.Context) {
	req := &struct {
		Name      string `form:"name"`
		Container string `form:"container"`
		Namespace string `form:"namespace"`
	}{}
	err := c.BindQuery(req)
	if err != nil {
		fmt.Println("pod log bind query failed ", err)
		c.JSON(http.StatusOK, gin.H{
			"msg":  "pod log bind query failed" + err.Error(),
			"data": nil,
		})
		return
	}

	logs, err := service.Pod.GetPodLog(req.Container, req.Name, req.Namespace)
	if err != nil {
		fmt.Println("pod log failed ", err)
		c.JSON(http.StatusOK, gin.H{
			"msg":  "pod log failed" + err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "pod log success",
		"data": logs,
	})
}

// 打印pod实时日志
func (p pod) PodLogSync(c *gin.Context) {
	req := &struct {
		Name      string `form:"name"`
		Container string `form:"container"`
		Namespace string `form:"namespace"`
	}{}
	err := c.BindQuery(req)
	if err != nil {
		fmt.Println("pod log bind query failed ", err)
		c.JSON(http.StatusOK, gin.H{
			"msg":  "pod log bind query failed" + err.Error(),
			"data": nil,
		})
		return
	}

	err = service.Pod.GetPodLogSync(req.Container, req.Name, req.Namespace)
	if err != nil {
		fmt.Println("pod log failed ", err)
		c.JSON(http.StatusOK, gin.H{
			"msg":  "pod log failed" + err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "pod log sync success",
		"data": nil,
	})
}
