package test

import (
	"fmt"
	logger "github.com/alecthomas/log4go"
	"github.com/gcggcg/k8s-core-components/k8s"
	"testing"
)

func TestGetNamesByNamespace(t *testing.T) {
	logger.Info("=================================TestGetNamesByNamespace=====================================")
	names, err := k8s.DefaultK8SMgr.GetAppNamesByNamespace(false)
	if err != nil {
		fmt.Println(fmt.Sprintf("获取域名空间: %s, 所有容器资源信息异常: %v", appNamespace, err))
		return
	}
	for _, info := range names {
		logger.Info("【域名空间: %s】 获取的容器信息名称: %s ", appNamespace, info)
	}
}
