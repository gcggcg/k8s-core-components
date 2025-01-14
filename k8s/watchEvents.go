package k8s

import (
	"context"
	logger "github.com/alecthomas/log4go"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"strings"
	"time"
)

/**
 *    Description: 监控器事件
 *    Date: 2024/11/28
 */

func (api *k8sApi) watchPodEvents(namespace string) {
	init := true
	// TODO: 待处理 k8s状态监控相关的缓存处理,所有事件执行后都会异步触发到监控模块进行缓存创建实时更新
	for {
		var restart bool
		watchChan, err := api.client.CoreV1().Pods(namespace).Watch(context.TODO(), metav1.ListOptions{Watch: true})
		if err != nil {
			err = logger.Error("域名: %s, 创建监控出现异常: :%v ", namespace, err)
			if init {
				panic(err)
			}
		} else {
			init = false
		}
		for {
			if restart {
				logger.Error("域名: %s, 监听事件出现异常: %v,等待重新创建监听器", namespace)
				watchChan.Stop()
				break
			}
			select {
			case <-api.exitCh:
				watchChan.Stop()
			case event := <-watchChan.ResultChan():
				if event.Type == watch.Error {
					restart = true
					break
				}
				pod, ok := event.Object.(*corev1.Pod)
				if !ok {
					continue
				}
				switch event.Type {
				case watch.Added:
					logger.Info("添加事件: 容器: %s,所在域名空间: %s,最新状态: %v", pod.Name, pod.Namespace, pod.Status.Phase)
				case watch.Modified:
					logger.Info("更新事件: 容器: %s,所在域名空间: %s,最新状态: %v", pod.Name, pod.Namespace, pod.Status.Phase)
				case watch.Deleted:
					logger.Warn("删除事件: 容器: %s,所在域名空间: %s,最新状态: %v", pod.Name, pod.Namespace, pod.Status.Phase)
				default:
					continue
				}
				if pod.Status.Phase != RunningStatus && len(pod.Status.Conditions) > 0 && pod.Status.Conditions[0].Message != "" {
					logger.Warn("容器: %s,所在域名空间: %s, 部署异常信息: %v", pod.Name, pod.Namespace, pod.Status.Conditions[0].Message)
				}
				// 删除事件, 删除缓存
				if event.Type == watch.Deleted {
					logger.Warn("容器: %s,所在域名空间: %s, 删除事件, 删除缓存", pod.Name, pod.Namespace)
					DefaultK8SMgr.DelCacheContainerMonitor(pod.Name, namespace)
				} else {
					// 信息变更缓存更新
					DefaultK8SMgr.SetCacheContainerInfo(pod.Name, namespace, ContainerInfo{
						Name:         strings.TrimRight(pod.Name, "-0"),
						HostIP:       pod.Status.HostIP,
						PodIP:        pod.Status.PodIP,
						Status:       string(pod.Status.Phase),
						ReStartCount: int(pod.Status.ContainerStatuses[0].RestartCount),
						NewStartAt:   pod.Status.ContainerStatuses[0].State.Running.StartedAt.Time,
					})
				}
			}
		}
		time.Sleep(1 * time.Minute)
	}
}
