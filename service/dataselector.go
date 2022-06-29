package service

// 数组分页 过滤 排序

import (
	//nwv1 "k8s.io/api/networking/v1"
	"sort"
	"strings"
	"time"
)

type dataSelector struct {
	GenericDataList []DataCell
	dataSelectQuery *DataSelectQuery
}

func (d *dataSelector) Len() int {
	return len(d.GenericDataList)
}

func (d *dataSelector) Swap(i, j int) {
	d.GenericDataList[i], d.GenericDataList[j] = d.GenericDataList[j], d.GenericDataList[i]
}

func (d *dataSelector) Less(i, j int) bool {
	return d.GenericDataList[j].GetCreation().Before(d.GenericDataList[i].GetCreation())
}

// 重写sort方法，排序 dataSelector
func (d *dataSelector) Sort() *dataSelector {
	sort.Sort(d)
	return d
}

// 过滤 比较元素的Name, 如果包含则返回
func (d *dataSelector) Filter() *dataSelector {
	if d.dataSelectQuery.FilterQuery.Name == "" {
		// 过滤字段为空，标识返回全量
		return d
	}

	var filter []DataCell
	for _, val := range d.GenericDataList {
		// 包含指定filter的添加到 filter 切片中
		if strings.Contains(val.GetName(), d.dataSelectQuery.FilterQuery.Name) {
			filter = append(filter, val)
		}
	}

	d.GenericDataList = filter
	return d
}

// 分页逻辑
func (d *dataSelector) Paginate() *dataSelector {
	limit := d.dataSelectQuery.PaginateQuery.Limit
	page := d.dataSelectQuery.PaginateQuery.Page

	// 参数合法性校验
	if limit <= 0 || page <= 0 {
		return d
	}

	start, end := limit*(page-1), limit*page

	if len(d.GenericDataList) < end {
		end = len(d.GenericDataList)
	}

	// 切分
	d.GenericDataList = d.GenericDataList[start:end]
	return d
}

// datacell 接口 用于各宗资源的list类型转换，转换后克使用 dataSelector的自定义排序方法
type DataCell interface {
	// 获取创建时间，根据时间排序
	GetCreation() time.Time
	// 获取名字，根据名字过滤
	GetName() string
}

// 定义过滤和分页属性，过滤 Name, 分页 Limit 和page
type DataSelectQuery struct {
	FilterQuery   *FilterQuery
	PaginateQuery *PaginateQuery
}

type FilterQuery struct {
	Name string
}

type PaginateQuery struct {
	Limit int
	Page  int
}
