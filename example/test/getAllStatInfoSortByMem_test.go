package test

import (
	logger "github.com/alecthomas/log4go"
	"github.com/gcggcg/k8s-core-components/k8s"
	"testing"
)

func TestGetAllStatInfoOfSortByMem(t *testing.T) {
	logger.Info("=================================TestGetAllStatInfoOfSortByMem=====================================")
	desc := true
	appNames := []string{"test-create-mysql2", "etcd-yun-v1"}
	k8s.DefaultK8SMgr.InitStatByNamespace(appNames, false)
	list := k8s.DefaultK8SMgr.GetAllStatInfoOfSortByMem(desc, appNamespace)
	for _, info := range list {
		logger.Info("【域名空间: %s】get all container stat info 总数据: %d, by desc:%v, info: %+v", appNamespace, len(list), desc, info)
	}
}
