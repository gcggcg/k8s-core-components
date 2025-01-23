package k8s

import (
	logger "github.com/alecthomas/log4go"
	"sort"
	"strings"
	"sync"
	"time"
)

/**
 *    Description: 容器监控管理CPU, 内存使用率,极其相关的容器信息
 *    Date: 2024/11/23
 */

type (
	StatInfo struct {
		Name    string   `json:"name"`
		CpuLoad LoadInfo `json:"CpuLoad"`
		MemLoad LoadInfo `json:"MemLoad"`
	}
	LoadInfo struct {
		Used  uint64  //内存单位byte
		Total uint64  //内存单位byte
		Ratio float64 //单位%
	}
	ContainerInfo struct {
		Name         string
		HostIP       string
		PodIP        string
		Status       string
		ReStartCount int
		NewStartAt   time.Time
	}
	ContainerMonitor struct {
		statInfo      *StatInfo
		containerInfo *ContainerInfo
		statsExitCh   chan bool
	}
	// ContainerCache 容器监控缓存
	ContainerCache struct {
		cache sync.Map
	}
)

func (cache *ContainerCache) addMonitor(name, namespace string) {
	var (
		monitor ContainerMonitor
		success bool
	)
	if body, ok := cache.cache.Load(name + "_" + namespace); ok {
		monitor, success = body.(ContainerMonitor)
		if !success || monitor.statInfo == nil {
			return
		}
	}
	timer := time.NewTimer(time.Second * 3)
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			if stat, err := DefaultK8SMgr.StatInfo(monitor.statInfo.Name, namespace); err != nil {
				logger.Info("【容器: %v】定时获取容器资源信息异常: %v", monitor.statInfo.Name, err)
				DefaultK8SMgr.DelContainerStatInfo(monitor.statInfo.Name, namespace)
				return
			} else {
				DefaultK8SMgr.SetCacheStatInfo(monitor.statInfo.Name, namespace, StatInfo{
					Name: monitor.statInfo.Name,
					CpuLoad: LoadInfo{
						Ratio: stat.CpuLoad.Ratio,
						Used:  stat.CpuLoad.Used,
						Total: stat.CpuLoad.Total,
					},
					MemLoad: LoadInfo{
						Ratio: stat.MemLoad.Ratio,
						Used:  stat.MemLoad.Used,
						Total: stat.MemLoad.Total,
					},
				})
			}
			timer.Reset(time.Second * 3)
		case <-monitor.statsExitCh:
			logger.Info("【容器: %v】容器正常的进行下线, 停止采集容器相关资源信息！", monitor.statInfo.Name)
			return
		}
	}
}
func (cache *ContainerCache) getCacheContainerInfo(name string) (ContainerInfo, bool) {
	if body, ok := cache.cache.Load(name); ok {
		monitor, success := body.(ContainerMonitor)
		if success && monitor.containerInfo != nil {
			return *monitor.containerInfo, true
		}
	}
	return ContainerInfo{}, false
}
func (cache *ContainerCache) getCacheStatInfo(name string) (StatInfo, bool) {
	if body, ok := cache.cache.Load(name); ok {
		monitor, success := body.(ContainerMonitor)
		if success && monitor.statInfo != nil {
			return *monitor.statInfo, true
		}
	}
	return StatInfo{}, false
}
func (cache *ContainerCache) setCacheContainerInfo(name string, containerInfo ContainerInfo) {
	if body, ok := cache.cache.Load(name); ok {
		monitor, success := body.(ContainerMonitor)
		if success {
			monitor.containerInfo = &containerInfo
			cache.cache.Store(name, monitor)
			return
		}
	}
	monitor := ContainerMonitor{
		containerInfo: &containerInfo,
		statsExitCh:   make(chan bool),
	}
	cache.cache.Store(name, monitor)
}

func (cache *ContainerCache) setCacheStatInfo(name string, statInfo StatInfo) {
	if body, ok := cache.cache.Load(name); ok {
		monitor, success := body.(ContainerMonitor)
		if success {
			monitor.statInfo = &statInfo
			cache.cache.Store(name, monitor)
			return
		}
	}
	monitor := ContainerMonitor{
		statInfo:    &statInfo,
		statsExitCh: make(chan bool),
	}
	cache.cache.Store(name, monitor)
}

func (cache *ContainerCache) delContainerStatInfo(name string) {
	if body, ok := cache.cache.Load(name); ok {
		monitor, success := body.(ContainerMonitor)
		if success && monitor.statInfo != nil {
			monitor.statInfo = nil
			cache.cache.Store(name, monitor)
		}
	}
}
func (cache *ContainerCache) delCacheContainerMonitor(name string) {
	if body, ok := cache.cache.Load(name); ok {
		monitor, success := body.(ContainerMonitor)
		if success && monitor.statInfo != nil {
			close(monitor.statsExitCh)
		}
	}
	cache.cache.Delete(name)
}

// getAllStats 获取所有nameSpace 下的app的StatInfo
func (cache *ContainerCache) getAllStats(nameSpace string) []*StatInfo {
	stats := make([]*StatInfo, 0)
	cache.cache.Range(func(key, value any) bool {
		if strings.HasSuffix(key.(string), nameSpace) {
			if statInfo, ok := cache.getCacheStatInfo(key.(string)); ok {
				stats = append(stats, &statInfo)
			} else {
				stats = append(stats, &StatInfo{Name: strings.TrimRight(key.(string), "_"+nameSpace), CpuLoad: LoadInfo{Ratio: 0, Used: 0}, MemLoad: LoadInfo{Ratio: 0, Used: 0}})
			}
		}
		return true
	})
	return stats
}
func (cache *ContainerCache) getAllStatInfoOfSortByCpu(desc bool, nameSpace string) []*StatInfo {
	stats := cache.getAllStats(nameSpace)
	if desc {
		sort.Sort(CpuDescSorter(stats))
	} else {
		sort.Sort(CpuAscSorter(stats))
	}
	return stats
}

func (cache *ContainerCache) getAllStatInfoOfSortByMem(desc bool, nameSpace string) []*StatInfo {
	stats := cache.getAllStats(nameSpace)
	if desc {
		sort.Sort(MemDescSorter(stats))
	} else {
		sort.Sort(MemAscSorter(stats))
	}
	return stats
}
