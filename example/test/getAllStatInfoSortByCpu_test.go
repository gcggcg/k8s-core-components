package test

import (
	"fmt"
	logger "github.com/alecthomas/log4go"
	"github.com/gcggcg/k8s-core-components/k8s"
	"testing"
	"time"
)

func TestGetAllStatInfoOfSortByCpu(t *testing.T) {
	logger.Info("=================================TestGetAllStatInfoOfSortByCpu=====================================")
	desc := true
	appNames := []string{"test-create-mysql2", "etcd-yun-v1"}
	k8s.DefaultK8SMgr.InitStatByNamespace(appNames, false)
	list := k8s.DefaultK8SMgr.GetAllStatInfoOfSortByCpu(desc, appNamespace)
	for _, info := range list {
		logger.Info("【域名空间: %s】get all container stat info 总数据: %d, by desc:%v, info: %+v", appNamespace, len(list), desc, info)
	}
	fmt.Println(list)
	time.Sleep(3 * time.Second)
}
