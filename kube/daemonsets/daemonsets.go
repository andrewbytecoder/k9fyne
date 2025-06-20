package daemonsets

import (
	"context"
	"go.uber.org/zap"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"time"
)

type K9DaemonSetsInfo struct {
	ctx     context.Context
	log     *zap.Logger
	kc      *kubernetes.Clientset
	podList *corev1.PodList
	podsMap map[string]*KubeDaemonSetsInfo // namespace to pods
}

func NewK9DaemonSetsInfo(ctx context.Context, k *kubernetes.Clientset, log *zap.Logger) *K9DaemonSetsInfo {
	return &K9DaemonSetsInfo{
		ctx:     ctx,
		log:     log,
		kc:      k,
		podList: nil,
		podsMap: make(map[string]*KubeDaemonSetsInfo),
	}
}

type KubeDaemonSetsInfo struct {
	podMap map[string]corev1.Pod // Pod name 2 Pod
}

type KubeDaemonSetsInfoInterface interface {
	GetAllNamespace() []string
	GetDaemonSetsInfoByNamespace(namespace string) (*appsv1.DaemonSetList, error)
}

func (p *K9DaemonSetsInfo) GetAllNamespace() []string {
	ctx, cancel := context.WithTimeout(p.ctx, 10*time.Second)
	defer cancel()
	// 3. 获取所有的 Namespace
	listItem, err := p.kc.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		p.log.Error("无法获取 Namespace 列表", zap.Error(err))
		return nil
	}
	allNamespace := make([]string, 0)
	for _, item := range listItem.Items {
		allNamespace = append(allNamespace, item.Name)
	}
	return allNamespace
}

func (p *K9DaemonSetsInfo) GetDaemonSetsInfoByNamespace(namespace string) (*appsv1.DaemonSetList, error) {
	ctx, cancel := context.WithTimeout(p.ctx, 10*time.Second)
	defer cancel()

	return p.kc.AppsV1().DaemonSets(namespace).List(ctx, metav1.ListOptions{})
}
