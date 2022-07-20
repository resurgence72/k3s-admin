package service

import (
	"context"
	"encoding/json"
	"fmt"
	appv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strconv"
	"time"
)

type deployCell appv1.Deployment

func (d deployCell) GetCreation() time.Time {
	return d.CreationTimestamp.Time
}

func (d deployCell) GetName() string {
	return d.Name
}

var Deployment deployment

type deployment struct{}

// toCells 方法用于将deploy类型数组转换成datacell类型
func (d deployment) toCells(std []appv1.Deployment) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = deployCell(std[i])
	}
	return cells
}

// fromCell 将dataCell类型转换为pod
func (d deployment) fromCells(std []DataCell) []appv1.Deployment {
	deploys := make([]appv1.Deployment, len(std))
	for i := range std {
		// 接口断言
		deploys[i] = appv1.Deployment(std[i].(deployCell))
	}
	return deploys
}

// 更新deploy副本数
func (d deployment) ScaleDeployment(deployName, namespace string, scaleNum int) (int32, error) {
	scale, err := K8s.ClientSet.AppsV1().Deployments(namespace).GetScale(context.TODO(), deployName, metav1.GetOptions{})
	if err != nil {
		fmt.Println("获取scale失败", err)
		return 0, err
	}

	if int32(scaleNum) == scale.Spec.Replicas {
		return scale.Spec.Replicas, nil
	}

	// 修改副本数
	scale.Spec.Replicas = int32(scaleNum)
	// 更新副本数
	updateScale, err := K8s.ClientSet.AppsV1().Deployments(namespace).UpdateScale(context.TODO(), deployName, scale, metav1.UpdateOptions{})
	if err != nil {
		fmt.Println("更新scale失败", err)
		return 0, err
	}

	return updateScale.Spec.Replicas, nil
}

type kv map[string]interface{}

// 重启deploy 本质上就是修改yaml, 并apply
func (d deployment) RestartDeployment(deployName, namespace string) error {
	// 组装数据
	patchData := kv{
		"spec": kv{
			"template": kv{
				"spec": kv{
					"containers": []kv{
						{
							"name": deployName,
							"env": []map[string]string{{
								"name":  "RESTART_",
								"value": strconv.FormatInt(time.Now().Unix(), 10),
							}},
						},
					},
				},
			},
		},
	}

	// 序列化为字节
	patchByte, err := json.Marshal(patchData)
	if err != nil {
		fmt.Println("marshal 失败", err)
		return err
	}

	_, err = K8s.ClientSet.AppsV1().Deployments(namespace).Patch(
		context.TODO(),
		deployName,
		// pt 类型为固定的4个
		"application/strategic-merge-patch+json",
		patchByte, metav1.PatchOptions{},
	)

	return nil
}
