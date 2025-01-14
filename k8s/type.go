package k8s

import (
	corev1 "k8s.io/api/core/v1"
)

const (
	RESTART_Always    = "Always"    // Always: 容器总是重启，除非被手动停止。
	RESTART_Never     = "Never"     // Never: 容器不会自动重启。
	RESTART_OnFailure = "OnFailure" // OnFailure: 容器在退出时，如果退出码为0，则不会重启。
	ENV_ULIMIT_NAME   = "ULIMIT"
	ENV_PRIVILEGED    = "PRIVILEGED"
	ENV_MACADDRESS    = "MACADDRESS"
	RunningStatus     = "Running"
	Action            = "start"
)

// ContainerCreateInfo app 创建的容器信息
type (
	ContainerCreateInfo struct {
		Name     string
		NodeName string
		Image    string
		// 是否是"host模式服务"
		HostNetwork  bool
		Label        map[string]string
		Env          []corev1.EnvVar
		Port         []corev1.ContainerPort
		Volumes      []corev1.Volume      // 挂载卷,映射主机路径
		VolumeMounts []corev1.VolumeMount // 挂载卷,映射容器路径
		Restart      string
		Privileged   bool
	}
)

type (
	CreateReqInfo struct {
		Name     string
		NodeName string
		Image    string
		// 是否是"host模式服务"
		HostNetwork bool
		Label       []LabelInfo
		Env         []EnvInfo
		Port        []PortInfo
		Volume      []VolumeInfo
		Restart     string
	}
	LabelInfo struct {
		Key   string
		Value string
	}
	EnvInfo struct {
		Key   string
		Value string
	}
	PortInfo struct {
		Protocol  string
		InnerPort uint16
		OuterPort uint16
	}
	VolumeInfo struct {
		InnerPath string
		OuterPath string
	}
	NodeInfo struct {
		Name          string
		Addr          string
		Status        string // "Healthy" 为正常
		DockerVersion string
		Total         int //容器总数
		Running       int //运行容器数
		Paused        int //暂停容器数
		Stopped       int //停止容器数
	}
)

//node, err := clientset.CoreV1().Nodes().Get(context.Background(), nodeName, metav1.GetOptions{})
//if err != nil {
//log.Fatalf("failed to get node: %v", err)
//}
//
//// 检查节点的 Conditions 是否存在 Ready 状态
//isReady := false
//for _, condition := range node.Status.Conditions {
//if condition.Type == v1.NodeReady && condition.Status == v1.ConditionTrue {
//isReady = true
//break
//}
//}
