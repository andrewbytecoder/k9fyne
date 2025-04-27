package topo

import (
	"github.com/andrewbytecoder/k9fyne/kube/pod"
	"github.com/andrewbytecoder/k9fyne/widgets"
)

type K9Info struct {
	pod.KubePodInfoInterface
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
		lc.KubePodInfoInterface = podInfoInterface
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
		tutorial.Data = k.KubePodInfoInterface
		return nil
	default:
		return nil
	}
}
