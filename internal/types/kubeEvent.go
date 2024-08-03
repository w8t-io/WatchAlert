package types

var EventResourceTypeList = []string{
	"Pods",
	"Nodes",
	"PVC/PV",
	"HPA",
}

type reason struct {
	Type   string `json:"type"`
	TypeCN string `json:"typeCN"`
}

var EventReasonLMapping = map[string][]reason{
	"Pods": {
		{
			"Failed",
			"容器启动失败",
		},
		{
			"Unhealthy",
			"容器健康状况不佳",
		},
		{
			"CrashLoopBackOff",
			"容器反复崩溃和重启",
		},
		{
			"FailedMount",
			"挂载卷失败",
		},
		{
			"FailedAttachVolume",
			"附加卷到节点失败",
		},
		{
			"DeadlineExceeded",
			"Pod超过其运行期限",
		},
	},
	"Nodes": {
		{
			"NodeNotReady",
			"节点处于不可用状态",
		},
		{
			"NodeUnderMemoryPressure",
			"节点处于内存压力下",
		},
		{
			"NodeUnderDiskPressure",
			"节点处于磁盘压力下",
		},
		{},
	},
	"PVC/PV": {
		{
			"FailedBinding",
			"PVC/PV绑定失败",
		},
	},
	"HPA": {
		{
			"FailedRescale",
			"调整副本数失败",
		},
		{
			"FailedGetResourceMetric",
			"获取资源指标失败",
		},
		{
			"FailedGetExternalMetric",
			"获取外部指标失败",
		},
	},
}
