package test

import (
	logger "github.com/alecthomas/log4go"
	"github.com/gcggcg/k8s-core-components/k8s"
)

func init() {
	// 加载日志配置
	logger.LoadConfiguration("../conf/log.xml")
	// 初始化k8s管理器
	if err := k8s.DefaultK8SMgr.Init("../conf/k8s.conf", "plate-system", "plate-app"); err != nil {
		panic(logger.Error("init k8s manager failed, error[%s]", err))
	}
	// 启动k8s管理器
	k8s.DefaultK8SMgr.Start()
}
