package topo

import (
	"github.com/andrewbytecoder/k9fyne/kube/daemonsets"
	"github.com/andrewbytecoder/k9fyne/kube/deployment"
	"github.com/andrewbytecoder/k9fyne/kube/pod"
	pullimage "github.com/andrewbytecoder/k9fyne/kube/pull_image"
	"github.com/andrewbytecoder/k9fyne/kube/service"
	"github.com/andrewbytecoder/k9fyne/kube/statefulsets"
	"github.com/andrewbytecoder/k9fyne/widgets"
)

type K9Info struct {
	kubePodInfoInterface          pod.KubePodInfoInterface
	kubeServiceInfoInterface      service.KubeServiceInfoInterface
	kubeDeploymentInfoInterface   deployment.KubeDeploymentInfoInterface
	kubeDaemonSetsInfoInterface   daemonsets.KubeDaemonSetsInfoInterface
	kubeStatefulSetsInfoInterface statefulsets.KubeStatefulSetsInfoInterface
	kubePullImageInfoInterface    pullimage.KubePullImageInfoInterface
}

var K9InfoHandler *K9Info

// An Option configures a Logger.
type Option interface {
	apply(*K9Info)
}

// optionFunc wraps a func so it satisfies the Option interface.
type optionFunc func(*K9Info)

func (f optionFunc) apply(k *K9Info) {
	f(k)
}

// SetPodInfoInterface 设置pod info interface
func SetPodInfoInterface(podInfoInterface pod.KubePodInfoInterface) Option {
	return optionFunc(func(lc *K9Info) {
		lc.kubePodInfoInterface = podInfoInterface
	})
}

// SetServiceInfoInterface 设置service info interface
func SetServiceInfoInterface(serviceInfoInterface service.KubeServiceInfoInterface) Option {
	return optionFunc(func(lc *K9Info) {
		lc.kubeServiceInfoInterface = serviceInfoInterface
	})
}

// SetDeploymentInfoInterface 设置deployment
func SetDeploymentInfoInterface(infoInterface deployment.KubeDeploymentInfoInterface) Option {
	return optionFunc(func(lc *K9Info) {
		lc.kubeDeploymentInfoInterface = infoInterface
	})
}

// SetDaemonSetsInfoInterface 设置deployment
func SetDaemonSetsInfoInterface(infoInterface daemonsets.KubeDaemonSetsInfoInterface) Option {
	return optionFunc(func(lc *K9Info) {
		lc.kubeDaemonSetsInfoInterface = infoInterface
	})
}

// SetStatefulSetsInfoInterface 设置deployment
func SetStatefulSetsInfoInterface(infoInterface statefulsets.KubeStatefulSetsInfoInterface) Option {
	return optionFunc(func(lc *K9Info) {
		lc.kubeStatefulSetsInfoInterface = infoInterface
	})
}

// SetPullImageInfoInterface 设置deployment
func SetPullImageInfoInterface(infoInterface pullimage.KubePullImageInfoInterface) Option {
	return optionFunc(func(lc *K9Info) {
		lc.kubePullImageInfoInterface = infoInterface
	})
}

func (k *K9Info) WithOptions(opts ...Option) *K9Info {
	for _, opt := range opts {
		opt.apply(k)
	}
	return k
}

func NewK9Info(op ...Option) *K9Info {
	k9Info := &K9Info{}

	for _, opt := range op {
		opt.apply(k9Info)
	}
	return k9Info
}

func (k *K9Info) FetchData(tutorial *widgets.Tutorial) error {
	switch tutorial.Title {
	case "Pod info":
		tutorial.Data = k.kubePodInfoInterface
		return nil
	case "Service info":
		tutorial.Data = k.kubeServiceInfoInterface
		return nil
	case "Deployment info":
		tutorial.Data = k.kubeDeploymentInfoInterface
		return nil
	case "DaemonSets info":
		tutorial.Data = k.kubeDaemonSetsInfoInterface
		return nil
	case "StatefulSets info":
		tutorial.Data = k.kubeStatefulSetsInfoInterface
		return nil
	case "PullImage info":
		tutorial.Data = k.kubePullImageInfoInterface
		return nil
	default:
		return nil
	}
}
