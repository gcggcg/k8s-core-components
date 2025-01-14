package k8s

import (
	logger "github.com/alecthomas/log4go"
	"github.com/golang/protobuf/proto"
	corev1 "k8s.io/api/core/v1"
)

/**
 *    Description: k8s的管理模块
 *    Date: 2024/10/10
 */

var (
	DefaultK8SMgr ManagerAPI
)

type ManagerAPI interface {
	Init(conf, systemNamespace, appNamespace string) error
	Start()
	StatefulSetCreate(info *CreateReqInfo, isTry bool) error
	StatefulSetDelete(name string, isTry bool) error
	StatefulSetRunOrStop(name, action string, isTry bool) error
	StatefulSetRestart(name string, isTry bool) error
	ContainerInfo(name string, namespace string) (ContainerInfo, error)
	StatInfo(name string, namespace string) (StatInfo, error)
	GetCacheContainerInfo(name string, isSys bool) (ContainerInfo, error)
	GetCacheStatInfo(name string, isSys bool) (StatInfo, error)
	SetCacheContainerInfo(name string, namespace string, info ContainerInfo)
	SetCacheStatInfo(name string, namespace string, info StatInfo)
	DelCacheContainerMonitor(name string, namespace string)
	Stop()
}

type ManagerK8s struct {
	k8sConfig       string
	api             API
	systemNamespace string
	appNamespace    string
	eventExitCh     chan bool
	containerCache  *ContainerCache
}

func init() {
	DefaultK8SMgr = new(ManagerK8s)
}
func (manage *ManagerK8s) Init(conf, systemNamespace, appNamespace string) error {
	manage.k8sConfig = conf
	manage.systemNamespace = systemNamespace
	manage.appNamespace = appNamespace
	manage.api = new(k8sApi)
	manage.containerCache = new(ContainerCache)
	if err := manage.api.init(manage.k8sConfig, manage.systemNamespace, manage.appNamespace); err != nil {
		return logger.Error("init k8s api failed, error[%s]", err)
	}
	return nil
}
func (manage *ManagerK8s) Start() {
	logger.Info("==============k8s start=============")
	manage.eventExitCh = make(chan bool)
	go manage.api.watchPodEvents(manage.systemNamespace)
	go manage.api.watchPodEvents(manage.appNamespace)
}
func (manage *ManagerK8s) Stop() {
	logger.Info("==============k8s stop=============")
	close(manage.eventExitCh)
	manage.containerCache = nil
	manage.api.exit()
}

func (manage *ManagerK8s) StatefulSetCreate(info *CreateReqInfo, isTry bool) error {
	createInfo := &ContainerCreateInfo{Name: info.Name, NodeName: info.NodeName, Image: info.Image, HostNetwork: info.HostNetwork, Restart: info.Restart} // 只考虑两种要么是host模式要么是非host模式,
	// label 标签
	createInfo.Label = make(map[string]string)
	for _, label := range info.Label {
		createInfo.Label[label.Key] = label.Value
	}
	// 添加特殊标签
	createInfo.Label["node_ip"] = info.NodeName
	createInfo.Label["app"] = info.Name
	// env 环境变量
	for _, envInfo := range info.Env {
		if envInfo.Key != ENV_MACADDRESS && envInfo.Key != ENV_PRIVILEGED && envInfo.Key != ENV_ULIMIT_NAME {
			createInfo.Env = append(createInfo.Env, corev1.EnvVar{Name: envInfo.Key, Value: envInfo.Value})
		} else {
			// (通过环境变量)自定义属性
			switch envInfo.Key {
			case ENV_PRIVILEGED:
				if envInfo.Value == "true" || envInfo.Value == "1" {
					createInfo.Privileged = true
				}
			case ENV_ULIMIT_NAME:
				logger.Info("【容器: %s】set uLimit value [%s], k8s服务不支持,可直接配置主机系统级别的uLimit,默认会使用主机!", envInfo.Value)
			case ENV_MACADDRESS:
				logger.Info("【容器: %s】set macAddress value [%s], k8s服务不支持!", envInfo.Value)
			}
		}
	}
	// port 端口,非host模式下支持配置: Protocol 协议默认TCP, 仅仅支持: TCP,UDP,SCTP
	if !info.HostNetwork {
		for _, port := range info.Port {
			if port.Protocol == "TCP" || port.Protocol == "UDP" || port.Protocol == "SCTP" {
				createInfo.Port = append(createInfo.Port, corev1.ContainerPort{ContainerPort: int32(port.InnerPort), HostPort: int32(port.OuterPort), Protocol: corev1.Protocol(port.Protocol)})
			} else {
				return logger.Warn("【容器: %s】port protocol[%s] is not support, only support: TCP,UDP,SCTP", info.Name, port.Protocol)
			}
		}
	}
	// volume 卷映射
	volumeName := info.Name + "-data"
	for _, volume := range info.Volume {
		createInfo.Volumes = append(createInfo.Volumes, corev1.Volume{Name: volumeName, VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{Path: volume.OuterPath, Type: (*corev1.HostPathType)(proto.String(string(corev1.HostPathDirectoryOrCreate)))}}})
		createInfo.VolumeMounts = append(createInfo.VolumeMounts, corev1.VolumeMount{Name: volumeName, MountPath: volume.InnerPath})
	}
	logger.Info("【容器: %s】create container 命令执行中...目标服务器: %s ", info.Name, info.NodeName)
	return manage.api.statefulSetCreate(manage.appNamespace, createInfo, isTry)
}
func (manage *ManagerK8s) StatefulSetDelete(name string, isTry bool) error {
	logger.Info("【容器: %s】delete container 命令执行中...", name)
	return manage.api.statefulSetDelete(name, manage.appNamespace, isTry)
}

func (manage *ManagerK8s) StatefulSetRunOrStop(name, action string, isTry bool) error {
	logger.Info("【容器: %s】action container 命令执行中... 容器操作: [%s]", name, action)
	return manage.api.statefulSetRunOrStop(name, manage.appNamespace, action, isTry)
}

func (manage *ManagerK8s) StatefulSetRestart(name string, isTry bool) error {
	logger.Info("【容器: %s】restart container 命令执行中... ", name)
	return manage.api.statefulSetRestart(name, manage.appNamespace, isTry)
}

func (manage *ManagerK8s) ContainerInfo(name string, namespace string) (ContainerInfo, error) {
	logger.Info("【容器: %s】 get container info 命令执行中... ", name)
	return manage.api.containerInfo(name, namespace)
}

func (manage *ManagerK8s) StatInfo(name string, namespace string) (StatInfo, error) {
	logger.Info("【容器: %s】 get container stat info 命令执行中... ", name)
	return manage.api.containerMetricStat(name, namespace)
}

func (manage *ManagerK8s) GetCacheContainerInfo(name string, isSys bool) (ContainerInfo, error) {
	logger.Info("【容器: %s】 get container cache info 命令执行中... ", name)
	var namespace string
	if isSys {
		namespace = manage.systemNamespace
	} else {
		namespace = manage.appNamespace
	}
	cacheInfo, ok := manage.containerCache.getCacheContainerInfo(name + "_" + namespace)
	if ok {
		return cacheInfo, nil
	} else {
		if info, err := manage.ContainerInfo(name, namespace); err != nil {
			return ContainerInfo{}, err
		} else {
			manage.containerCache.setCacheContainerInfo(name+"_"+namespace, info)
			return info, nil
		}
	}
}
func (manage *ManagerK8s) GetCacheStatInfo(name string, isSys bool) (StatInfo, error) {
	logger.Info("【容器: %s】 get container cache stat  命令执行中... ", name)
	var namespace string
	if isSys {
		namespace = manage.systemNamespace
	} else {
		namespace = manage.appNamespace
	}
	cacheInfo, ok := manage.containerCache.getCacheStatInfo(name + "_" + namespace)
	if ok {
		return cacheInfo, nil
	} else {
		if info, err := manage.StatInfo(name, namespace); err != nil {
			return StatInfo{}, err
		} else {
			manage.containerCache.setCacheStatInfo(name+"_"+namespace, info)
			return info, nil
		}
	}
}
func (manage *ManagerK8s) SetCacheContainerInfo(name string, namespace string, info ContainerInfo) {
	logger.Info("【容器: %s】 set container cache info 命令执行中... ", name)
	manage.containerCache.setCacheContainerInfo(name+"_"+namespace, info)
}

func (manage *ManagerK8s) SetCacheStatInfo(name string, namespace string, info StatInfo) {
	logger.Info("【容器: %s】 set container cache info 命令执行中... ", name)
	manage.containerCache.setCacheStatInfo(name+"_"+namespace, info)
}
func (manage *ManagerK8s) DelCacheContainerMonitor(name string, namespace string) {
	logger.Info("【容器: %s】 delete container cache 命令执行中... ", name)
	manage.containerCache.delCacheContainerMonitor(name + "_" + namespace)
}
