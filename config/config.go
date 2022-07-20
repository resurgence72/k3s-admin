package config

import "time"

const (
	ListenAddr = ":9090"
	Kubeconfig = "D:\\kubeconfig"
	// 容器日志行数
	PodLogTailLine = 20
	
	//数据库配置
	BbHost = "127.0.0.1"
	DbPort = "3306"
	DbUser = "root"
	DbPwd = "123"
	LogMode = false

	// 连接池配置
	MaxIdleConns = 10 // 最大空闲连接
	MaxOpenConns = 100 // 最大连接数
	MaxLifeTime = 30 * time.Second // 最大生存时间

)
