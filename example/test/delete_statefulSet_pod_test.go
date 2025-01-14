package test

import (
	logger "github.com/alecthomas/log4go"
	"github.com/gcggcg/k8s-core-components/k8s"
	"testing"
)

func TestDeletePod(t *testing.T) {
	logger.Info("=================================TestDeletePod=================================")
	err := k8s.DefaultK8SMgr.StatefulSetDelete("test-create-mysql", false)
	if err != nil {
		logger.Error("【容器: test-create-mysql】action container 命令执行TestDeletePod失败, error[%s]", err)
		return
	} else {
		logger.Info("【容器: test-create-mysql】action container 命令执行TestDeletePod成功")
	}
}
