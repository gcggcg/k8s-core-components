package test

import (
	logger "github.com/alecthomas/log4go"
	"github.com/gcggcg/k8s-core-components/k8s"
	"testing"
)

func TestStopPod(t *testing.T) {
	logger.Info("=================================TestStopPod=================================")
	err := k8s.DefaultK8SMgr.StatefulSetRunOrStop("test-create-mysql", "stop", true)
	if err != nil {
		logger.Error("【容器: test-create-mysql】action container 命令执行TestStopPod失败, error[%s]", err)
		return
	} else {
		logger.Info("【容器: test-create-mysql】action container 命令执行TestStopPod成功")
	}
}
