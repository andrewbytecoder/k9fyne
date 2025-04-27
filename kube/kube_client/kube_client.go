package kubeclient

import (
	"go.uber.org/zap"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func NewKubeClient(configFile string, log *zap.Logger) (*kubernetes.Clientset, error) {
	// 1. 加载 kubeconfig 配置文件
	var kubeconfig string
	kubeconfig = filepath.Join(configFile)

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Info("无法加载 kubeconfig")
		return nil, err
	}

	// 2. 创建 Kubernetes 客户端
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Info("无法创建 Kubernetes 客户端")
		return nil, err
	}

	return clientSet, nil
}
