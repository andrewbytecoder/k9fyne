package config

import (
	"github.com/andrewbytecoder/k9fyne/kube"
	"go.uber.org/zap"
)

type Cfg struct {
	log *zap.Logger
	SSH *kube.Client
}

func NewConfig(l *zap.Logger) *Cfg {
	return &Cfg{
		log: l,
		SSH: kube.NewSSHClient(l),
	}
}

func (c *Cfg) GetSSH() *kube.Client {
	return c.SSH
}
