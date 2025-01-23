package test

import (
	"fmt"
	logger "github.com/alecthomas/log4go"
	"github.com/gcggcg/k8s-core-components/k8s"
	"strconv"
	"testing"
)

func TestGetPodStat(t *testing.T) {
	logger.Info("=================================TestGetPodStat====================================")
	info, err := k8s.DefaultK8SMgr.GetCacheStatInfo("etcd-yun-v1", false)
	if err != nil {
		logger.Error("【容器: test-create-mysql】get container stat info命令执行TestGetPodStat失败, error[%s]", err)
		return
	} else {
		total, _ := strconv.ParseFloat(fmt.Sprintf("%0.2f", float64(info.MemLoad.Total)/1024), 64)
		use, _ := strconv.ParseFloat(fmt.Sprintf("%0.2f", float64(info.MemLoad.Used)/1024), 64)
		logger.Info("【容器: test-create-mysql】get container stat info: %+v, 总内存: %v, 占用内存: %v", info, total, use)
	}
}
