package service

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"k3s-admin/config"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

// 定义podcell类型 实现 DataCell接口
type podCell corev1.Pod

func (p podCell) GetCreation() time.Time {
	return p.CreationTimestamp.Time
}

func (p podCell) GetName() string {
	return p.Name
}

var Pod pod

type pod struct{}

// toCells 方法用于将pod类型数组转换成datacell类型
func (p pod) toCells(std []corev1.Pod) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = podCell(std[i])
	}
	return cells
}

// fromCell 将dataCell类型转换为pod
func (p pod) fromCells(std []DataCell) []corev1.Pod {
	pods := make([]corev1.Pod, len(std))
	for i := range std {
		// 接口断言
		pods[i] = corev1.Pod(std[i].(podCell))
	}
	return pods
}

// 获取pod列表，过滤 排序 分页
type PodsResp struct {
	Total int          `json:"total"`
	Items []corev1.Pod `json:"items"`
}

func (p pod) GetPods(filterName, namespace string, limit, page int) (*PodsResp, error) {
	// 获取pod列表
	podList, err := K8s.ClientSet.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Println("corev1 get pods failed ", err)
		return nil, err
	}

	// 实例化 dataSelector对象
	selectableData := &dataSelector{
		GenericDataList: p.toCells(podList.Items),
		dataSelectQuery: &DataSelectQuery{
			FilterQuery: &FilterQuery{Name: filterName},
			PaginateQuery: &PaginateQuery{
				Limit: limit,
				Page:  page,
			},
		},
	}

	selectableData.Filter().Sort()
	resp := &PodsResp{
		Total: selectableData.Len(),
		Items: p.fromCells(selectableData.GenericDataList),
	}

	return resp, nil
}

// 获取pod详情
func (p pod) GetPodDetail(podName, namespace string) (*corev1.Pod, error) {
	pod, err := K8s.ClientSet.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		fmt.Println("get pod failed ", err)
		return nil, err
	}
	return pod, nil
}

// 删除pod
func (p pod) DeletePod(podName, namespace string) error {
	err := K8s.ClientSet.CoreV1().Pods(namespace).Delete(context.TODO(), podName, metav1.DeleteOptions{})
	if err != nil {
		fmt.Println("del pod failed ", err)
		return err
	}
	return nil
}

// 更新pod
func (p pod) UpdatePod(podName, namespace, content string) error {
	pod := &corev1.Pod{}
	err := json.Unmarshal([]byte(content), pod)
	if err != nil {
		fmt.Println("update pod  unmarshal failed ", err)
		return err
	}

	_, err = K8s.ClientSet.CoreV1().Pods(namespace).Update(context.TODO(), pod, metav1.UpdateOptions{})
	if err != nil {
		fmt.Println("update pod failed ", err)
		return err
	}
	return nil
}

// 获取pod 容器列表 觉得直接序列化到列表展示页即可

// 获取容器log
func (p pod) GetPodLog(containerName, podName, namespace string) (string, error) {
	lineLimit := int64(config.PodLogTailLine)
	options := &corev1.PodLogOptions{
		Container: containerName,
		Follow:    false,
		TailLines: &lineLimit,
	}

	// 获取req实例
	req := K8s.ClientSet.CoreV1().Pods(namespace).GetLogs(podName, options)
	// 发起req请求,返回一个 io.readCloser类型 （类似response.Body）
	podLogs, err := req.Stream(context.TODO())
	if err != nil {
		fmt.Println("req stream failed ", err)
		return "", nil
	}
	defer podLogs.Close()

	// 将resp body 写入缓冲区
	buf := &bytes.Buffer{}
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		fmt.Println("io copy failed ", err)
		return "", nil
	}

	return buf.String(), nil
}

// 获取容器实时日志
// https://cloud.tencent.com/developer/ask/sof/297098
func (p pod) GetPodLogSync(containerName, podName, namespace string) error {
	lineLimit := int64(config.PodLogTailLine)
	options := &corev1.PodLogOptions{
		Container: containerName,
		Follow:    true,
		TailLines: &lineLimit,
	}

	// 获取req实例
	req := K8s.ClientSet.CoreV1().Pods(namespace).GetLogs(podName, options)
	// 发起req请求,返回一个 io.readCloser类型 （类似response.Body）
	podLogs, err := req.Stream(context.TODO())
	if err != nil {
		fmt.Println("req stream failed ", err)
		return nil
	}

	// 测试 设置读取20s
	ctx, _ := context.WithTimeout(context.TODO(), time.Duration(30)*time.Second)
	go func() {
		// 一定要在内层关闭 否则永远也拿不到log
		defer podLogs.Close()

		reader := bufio.NewScanner(podLogs)
		for reader.Scan() {
			select {
			case <-ctx.Done():
				fmt.Println("超时退出")
				return
			default:
				line := reader.Text()
				fmt.Println("当前拿到line ", line)
			}
		}
	}()
	return nil
}
