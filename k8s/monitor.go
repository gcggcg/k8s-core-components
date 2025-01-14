package k8s

import (
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
	}
	// ContainerCache 容器监控缓存
	ContainerCache struct {
		cache sync.Map
	}
)

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
		}
	}
	monitor := ContainerMonitor{
		containerInfo: &containerInfo,
	}
	cache.cache.Store(name, monitor)
}

func (cache *ContainerCache) setCacheStatInfo(name string, statInfo StatInfo) {
	if body, ok := cache.cache.Load(name); ok {
		monitor, success := body.(ContainerMonitor)
		if success {
			monitor.statInfo = &statInfo
			cache.cache.Store(name, monitor)
		}
	}
	monitor := ContainerMonitor{
		statInfo: &statInfo,
	}
	cache.cache.Store(name, monitor)
}
func (cache *ContainerCache) delCacheContainerMonitor(name string) {
	cache.cache.Delete(name)
}
