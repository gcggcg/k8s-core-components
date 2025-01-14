package test

import (
	logger "github.com/alecthomas/log4go"
	"github.com/gcggcg/k8s-core-components/k8s"
	"testing"
)

// 测试服务更新
func TestRunPod(t *testing.T) {
	logger.Info("=================================TestRunPod=================================")
	err := k8s.DefaultK8SMgr.StatefulSetRunOrStop("test-create-mysql", "start", false)
	if err != nil {
		logger.Error("【容器: test-create-mysql】action container 命令执行TestRunPod失败, error[%s]", err)
		return
	} else {
		logger.Info("【容器: test-create-mysql】action container 命令执行TestRunPod成功")
	}
}
