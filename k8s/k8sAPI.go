package k8s

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/metrics/pkg/client/clientset/versioned"
	"strconv"
)

/**
 *    Description: k8s 的API
 *    Date: 2024/10/27
 */

type API interface {
	init(k8sConfig, systemNamespace, appNamespace string) error
	exit()
	statefulSetCreate(namespace string, info *ContainerCreateInfo, isTry ...bool) error // 业务app 创建
	statefulSetDelete(name, namespace string, isTry ...bool) error                      // 业务app 删除
	statefulSetRestart(name, namespace string, isTry ...bool) error                     // 容器重启
	statefulSetRunOrStop(name, namespace, action string, isTry ...bool) error           // 停止或者启动容器
	watchPodEvents(namespace string)                                                    // 容器运行状态监听
	containerInfo(name, namespace string) (ContainerInfo, error)                        // 容器信息
	containerMetricStat(name, namespace string) (StatInfo, error)                       // 容器监控信息
}

type k8sApi struct {
	client          *kubernetes.Clientset
	metric          *versioned.Clientset
	systemNamespace string
	appNamespace    string
	exitCh          chan bool
}

func (api *k8sApi) init(k8sConfig, systemNamespace, appNamespace string) error {
	if k8sConfig == "" {
		return errors.New("k8s 权限认证文件为空")
	}
	config, err := clientcmd.BuildConfigFromFlags("", k8sConfig)
	if err != nil {
		return fmt.Errorf("error building kubeconfig: %v", err)
	}
	// 创建 Kubernetes 客户端
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("error creating Kubernetes client: %v", err)
	}
	api.client = client
	metric, err := versioned.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("error creating metric client: %v", err)
	}
	api.metric = metric
	api.systemNamespace = systemNamespace
	api.appNamespace = appNamespace
	api.exitCh = make(chan bool)
	return nil
}

func (api *k8sApi) exit() {
	close(api.exitCh)
}
func (api *k8sApi) statefulSetCreate(namespace string, info *ContainerCreateInfo, isTry ...bool) error {
	if info == nil {
		return fmt.Errorf("ContainerCreateInfo nil")
	}
	var err error
	// 定义: StatefulSet
	statefulSet := &v1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      info.Name,
			Namespace: namespace,
			Labels:    info.Label, // 核心: map[string]string{"node_ip": info.NodeName, "app": info.Name}
		},
		Spec: v1.StatefulSetSpec{
			Replicas:    proto.Int32(0),
			ServiceName: info.Name, // 绑定service 定义DNS域名, 用于DNS域名解析,后面创建的时候
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{ // 匹配由该StatefulSet管理的pod的筛选标签, 和PodTemplateSpec进行匹配成功才可以执行Template的操作
					"app": fmt.Sprintf("%s-%s", info.NodeName, info.Name),
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": fmt.Sprintf("%s-%s", info.NodeName, info.Name),
					},
				},
				Spec: corev1.PodSpec{
					HostNetwork: info.HostNetwork,
					NodeSelector: map[string]string{
						"kubernetes.io/hostname": info.NodeName,
					},
					Volumes:       info.Volumes,
					RestartPolicy: corev1.RestartPolicy(info.Restart), // statefulSet仅仅支持: Always
					Containers: []corev1.Container{
						{
							Name:         info.Name, // 容器名称,设置唯一可以和StatefulSet名称设置一个,因为我们设计都是按照单个pod启动,方便我们进行查看
							Image:        info.Image,
							Ports:        info.Port,
							Env:          info.Env, // 添加环境变量
							VolumeMounts: info.VolumeMounts,
							SecurityContext: &corev1.SecurityContext{
								Privileged: proto.Bool(info.Privileged), // 是否使用特权模式, 有的APP需要使用
							},
						},
					},
				},
			},
		},
	}

	if len(isTry) > 0 && isTry[0] {
		_, err = api.client.AppsV1().StatefulSets(namespace).Create(context.Background(), statefulSet, metav1.CreateOptions{DryRun: []string{"All"}})
	} else {
		_, err = api.client.AppsV1().StatefulSets(namespace).Create(context.Background(), statefulSet, metav1.CreateOptions{})
	}
	return err
}

func (api *k8sApi) statefulSetDelete(name, namespace string, isTry ...bool) error {
	if len(isTry) > 0 && isTry[0] {
		return api.client.AppsV1().StatefulSets(namespace).Delete(context.Background(), name, metav1.DeleteOptions{DryRun: []string{"All"}})
	} else {
		return api.client.AppsV1().StatefulSets(namespace).Delete(context.Background(), name, metav1.DeleteOptions{})
	}
}

func (api *k8sApi) statefulSetRunOrStop(name, namespace, action string, isTry ...bool) error {
	statefulSet, err := api.client.AppsV1().StatefulSets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	if action == Action {
		statefulSet.Spec.Replicas = proto.Int32(1)
	} else {
		statefulSet.Spec.Replicas = proto.Int32(0)
	}
	if len(isTry) > 0 && isTry[0] {
		_, err = api.client.AppsV1().StatefulSets(namespace).Update(context.Background(), statefulSet, metav1.UpdateOptions{DryRun: []string{"All"}})
	} else {
		_, err = api.client.AppsV1().StatefulSets(namespace).Update(context.Background(), statefulSet, metav1.UpdateOptions{})
	}
	return err
}
func (api *k8sApi) statefulSetRestart(name, namespace string, isTry ...bool) error {
	podName := fmt.Sprintf("%s-0", name)
	if len(isTry) > 0 && isTry[0] {
		return api.client.CoreV1().Pods(namespace).Delete(context.Background(), podName, metav1.DeleteOptions{DryRun: []string{"All"}})
	} else {
		return api.client.CoreV1().Pods(namespace).Delete(context.Background(), podName, metav1.DeleteOptions{})
	}
}

func (api *k8sApi) containerInfo(name, namespace string) (ContainerInfo, error) {
	podName := fmt.Sprintf("%s-0", name)
	if pod, err := api.client.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{}); err != nil {
		return ContainerInfo{}, err
	} else {
		return ContainerInfo{
			Name:         name,
			HostIP:       pod.Status.HostIP,
			PodIP:        pod.Status.PodIP,
			Status:       string(pod.Status.Phase),
			ReStartCount: int(pod.Status.ContainerStatuses[0].RestartCount),
			NewStartAt:   pod.Status.ContainerStatuses[0].State.Running.StartedAt.Time,
		}, nil
	}
}

func (api *k8sApi) getNodeInfo() {

}
func (api *k8sApi) containerMetricStat(name, namespace string) (StatInfo, error) {
	podName := fmt.Sprintf("%s-0", name)
	var (
		isSys       bool
		totalCPUNum uint64
		totalMemNum uint64
	)
	if namespace == api.systemNamespace {
		isSys = true
	} else if namespace == api.appNamespace {
		isSys = false
	} else {
		return StatInfo{}, errors.New("get containerMetricStat namespace error")
	}
	containerInfo, err := DefaultK8SMgr.GetCacheContainerInfo(name, isSys)
	if err != nil {
		return StatInfo{}, err
	}
	if nodeInfo, err := api.client.CoreV1().Nodes().Get(context.TODO(), containerInfo.HostIP, metav1.GetOptions{}); err != nil {
		return StatInfo{}, err
	} else {
		totalCPU := nodeInfo.Status.Capacity[corev1.ResourceCPU]
		totalMemory := nodeInfo.Status.Capacity[corev1.ResourceMemory]
		// 转换为数字值
		totalCPUNum = uint64(totalCPU.Value())
		totalMemNum = uint64(totalMemory.Value())
	}
	if podMetric, err := api.metric.MetricsV1beta1().PodMetricses(namespace).Get(context.TODO(), podName, metav1.GetOptions{}); err != nil {
		return StatInfo{}, err
	} else {
		resourceCPU := podMetric.Containers[0].Usage[corev1.ResourceCPU]
		resourceMemory := podMetric.Containers[0].Usage[corev1.ResourceMemory]
		podUseCPURatio, _ := strconv.ParseFloat(fmt.Sprintf("%0.2f", float64(resourceCPU.Value())/float64(totalCPUNum)*100), 64)
		podUseMemoryRatio, _ := strconv.ParseFloat(fmt.Sprintf("%0.2f", float64(resourceMemory.Value())/float64(totalMemNum)*100), 64)
		return StatInfo{
			Name: name,
			CpuLoad: LoadInfo{
				Total: totalCPUNum,
				Used:  uint64(resourceCPU.Value()),
				Ratio: podUseCPURatio,
			},
			MemLoad: LoadInfo{
				Total: totalMemNum,
				Used:  uint64(resourceMemory.Value()),
				Ratio: podUseMemoryRatio,
			},
		}, nil
	}
}
