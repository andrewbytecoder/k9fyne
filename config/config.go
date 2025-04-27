package config

import (
	"context"
	kubeclient "github.com/andrewbytecoder/k9fyne/kube/kube_client"
	"go.uber.org/zap"
)

type Cfg struct {
	log *zap.Logger
	SSH *kubeclient.Client
}

func NewConfig(l *zap.Logger, ctx context.Context) *Cfg {
	return &Cfg{
		log: l,
		SSH: kubeclient.NewSSHClient(l, ctx),
	}
}

func (c *Cfg) GetSSH() *kubeclient.Client {
	return c.SSH
}
