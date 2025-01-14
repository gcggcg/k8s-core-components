package test

import (
	logger "github.com/alecthomas/log4go"
	"github.com/gcggcg/k8s-core-components/k8s"
	"testing"
)

func TestCreatePod(t *testing.T) {
	logger.Info("=================================TestCreatePod=================================")
	err := k8s.DefaultK8SMgr.StatefulSetCreate(&k8s.CreateReqInfo{
		Name:        "test-create-mysql",
		NodeName:    "127.0.0.1",
		Image:       "mysql:5.7.18",
		HostNetwork: true,
		Label: []k8s.LabelInfo{
			{Key: "test-label", Value: "TestCreatePod"},
		},
		Env: []k8s.EnvInfo{
			{Key: "MYSQL_ROOT_PASSWORD", Value: "123root"},
			{Key: "PRIVILEGED", Value: "true"}, // 特权模式
		},
		Port: []k8s.PortInfo{
			{InnerPort: 3306, OuterPort: 33061, Protocol: "TCP"},
		},
		Volume: []k8s.VolumeInfo{
			{InnerPath: "/var/lib/mysql", OuterPath: "/opt/data/mysql"},
		},
		Restart: k8s.RESTART_Always,
	}, false)
	if err != nil {
		logger.Error("【容器: test-create-mysql】action container 命令执行TestCreatePod失败, error[%s]", err)
		return
	} else {
		logger.Info("【容器: test-create-mysql】action container 命令执行TestCreatePod成功")
	}
}
