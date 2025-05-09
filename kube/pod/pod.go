package pod

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"time"
)

type K9PodInfo struct {
	ctx     context.Context
	log     *zap.Logger
	kc      *kubernetes.Clientset
	podList *corev1.PodList
	podsMap map[string]*KubePodInfo // namespace to pods
}

func NewK9PodInfo(ctx context.Context, k *kubernetes.Clientset, log *zap.Logger) *K9PodInfo {
	return &K9PodInfo{
		ctx:     ctx,
		log:     log,
		kc:      k,
		podList: nil,
		podsMap: make(map[string]*KubePodInfo),
	}
}

type KubePodInfo struct {
	podMap map[string]corev1.Pod // Pod name 2 Pod
}

type KubePodInfoInterface interface {
	GetAllNamespace() []string
	GetPodInfoByNamespace(namespace string) (*corev1.PodList, error)
}

func (p *K9PodInfo) GetAllNamespace() []string {
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

func (p *K9PodInfo) GetPodInfoByNamespace(namespace string) (*corev1.PodList, error) {
	ctx, cancel := context.WithTimeout(p.ctx, 10*time.Second)
	defer cancel()
	// 3. 获取所有的 Namespace
	listItem, err := p.kc.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		p.log.Error("无法获取 Namespace 列表", zap.Error(err))
		return nil, fmt.Errorf("无法获取 Namespace 列表")
	}

	return listItem, nil
}
