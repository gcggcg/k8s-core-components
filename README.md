### 组件说明

> 该组件是基于K8S的API进行二次封装，方便进行开发使用，核心思想就是将复杂的逻辑处理进行接口化。同时可插拔式接口设计可以根据需求进行自动封装新接口，以及动态服务扩展，从而实现自研k8s控制台服务进行定制化开发设计。组件设计的核心思想是作为三方库可以进行模块化的改动升级。

### 功能说明

*   支持容器的创建
*   支持容器的删除
*   支持容器的启动
*   支持容器的停止
*   支持容器的重启
*   支持容器的信息查询
*   支持容器的状态资源查询
*   支持容器的异步实时监控

### 版本说明

1.  v1.x 版本核心是支持有状态类型的容器StatefulSet服务(偏业务性,可以自行修改设计)。

### 使用实例

```go
// 初始化全局k8s核心模块以及日志模块
func init() {
	// 加载日志配置
	logger.LoadConfiguration("../conf/log.xml")
	// 初始化k8s管理器
	if err := k8s.DefaultK8SMgr.Init("../conf/k8s.conf", "plate-system", "plate-app"); err != nil {
		panic(logger.Error("init k8s manager failed, error[%s]", err))
	}
	// 启动k8s管理器
	k8s.DefaultK8SMgr.Start()
}
// 创建容器并且启动容器
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


```

### 注意事项
*   该组件依赖k8s的api,需要k8s集群环境支持;
*   执行组件前优先按照初始化的k8s管理器进行先创建相关的namespace;