package statefulsets

import (
	"context"
	"go.uber.org/zap"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"time"
)

type K9StatefulSetsInfo struct {
	ctx     context.Context
	log     *zap.Logger
	kc      *kubernetes.Clientset
	podList *corev1.PodList
	podsMap map[string]*KubeStatefulSetsInfo // namespace to pods
}

func NewK9StatefulSetsInfo(ctx context.Context, k *kubernetes.Clientset, log *zap.Logger) *K9StatefulSetsInfo {
	return &K9StatefulSetsInfo{
		ctx:     ctx,
		log:     log,
		kc:      k,
		podList: nil,
		podsMap: make(map[string]*KubeStatefulSetsInfo),
	}
}

type KubeStatefulSetsInfo struct {
	podMap map[string]corev1.Pod // Pod name 2 Pod
}

type KubeStatefulSetsInfoInterface interface {
	GetAllNamespace() []string
	GetStatefulSetsInfoByNamespace(namespace string) (*appsv1.StatefulSetList, error)
}

func (p *K9StatefulSetsInfo) GetAllNamespace() []string {
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

func (p *K9StatefulSetsInfo) GetStatefulSetsInfoByNamespace(namespace string) (*appsv1.StatefulSetList, error) {
	ctx, cancel := context.WithTimeout(p.ctx, 10*time.Second)
	defer cancel()

	return p.kc.AppsV1().StatefulSets(namespace).List(ctx, metav1.ListOptions{})
}
