package test

import (
	logger "github.com/alecthomas/log4go"
	"github.com/gcggcg/k8s-core-components/k8s"
	"testing"
)

func TestRestartPod(t *testing.T) {
	logger.Info("=================================TestRestartPod=================================")
	logger.Info("=================================TestGetPodStat====================================")
	err := k8s.DefaultK8SMgr.StatefulSetRestart("test-create-mysql", false)
	if err != nil {
		logger.Error("【容器: test-create-mysql】 container restart命令执行TestGetPodStat失败, error[%s]", err)
		return
	} else {
		logger.Info("【容器: test-create-mysql】 container restart命令执行TestRestartPod 成功!")
	}
}
