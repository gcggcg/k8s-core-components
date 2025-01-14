package test

import (
	logger "github.com/alecthomas/log4go"
	"github.com/gcggcg/k8s-core-components/k8s"
	"testing"
)

func TestGetPodInfo(t *testing.T) {
	logger.Info("=================================TestGetPodInfo=================================")
	info, err := k8s.DefaultK8SMgr.GetCacheContainerInfo("test-create-mysql", false)
	if err != nil {
		logger.Error("【容器: test-create-mysql】get container container info命令执行TestGetPodStat失败, error[%s]", err)
		return
	} else {
		logger.Info("【容器: test-create-mysql】get container container info: %+v", info)
	}
}
