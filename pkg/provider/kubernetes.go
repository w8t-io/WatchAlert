package provider

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"time"
	"watchAlert/internal/global"
)

type KubernetesClient struct {
	Cli *kubernetes.Clientset
	Ctx context.Context
}

func NewKubernetesClient(ctx context.Context, kubeConfigContent string) (KubernetesClient, error) {
	// 如果配置内容为空，则去默认目录下取配置文件的内容
	if kubeConfigContent == "" {
		kubeConfigContent = os.Getenv("HOME") + "/.kube/config"
	}

	// 如果默认的配置文件Path实际是一个目录，那么跳过
	if _, err := os.Stat(kubeConfigContent); err == nil {
		content, err := os.ReadFile(kubeConfigContent)
		if err != nil {
			global.Logger.Sugar().Error(err.Error())
			return KubernetesClient{}, err
		}
		kubeConfigContent = string(content)
	}

	// 构建配置
	configBytes := []byte(kubeConfigContent)
	config, err := clientcmd.RESTConfigFromKubeConfig(configBytes)
	if err != nil {
		global.Logger.Sugar().Error(err.Error())
		return KubernetesClient{}, err
	}

	// 新建客户端
	cs, err := kubernetes.NewForConfig(config)
	if err != nil {
		global.Logger.Sugar().Error(err.Error())
	}

	return KubernetesClient{
		Cli: cs,
		Ctx: ctx,
	}, nil
}

func (a KubernetesClient) GetWarningEvent(reason string, scope int) (*corev1.EventList, error) {
	// 获取所有命名空间的事件列表
	list, err := a.Cli.CoreV1().Events(corev1.NamespaceAll).List(a.Ctx, metav1.ListOptions{})
	if err != nil {
		return list, err
	}

	// 创建一个新的 EventList 用于存储指定的 Reason 事件
	warningEvents := &corev1.EventList{}
	cutoffTime := time.Now().Add(-time.Duration(scope) * time.Minute)
	for _, event := range list.Items {
		// 检查事件的 Reason 和事件发生时间
		eventTime := event.LastTimestamp.Time // 使用 LastTimestamp 作为事件时间
		if event.Reason == reason && eventTime.After(cutoffTime) {
			warningEvents.Items = append(warningEvents.Items, event)
		}
	}

	return warningEvents, nil
}
