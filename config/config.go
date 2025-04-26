package config

import (
	"github.com/andrewbytecoder/k9fyne/controller"
	"go.uber.org/zap"
)

type Cfg struct {
	log *zap.Logger
	SSH *controller.SSHClient
}

func NewConfig(l *zap.Logger) *Cfg {
	return &Cfg{
		log: l,
		SSH: controller.NewSSHClient(l),
	}
}

func (c *Cfg) GetSSH() *controller.SSHClient {
	return c.SSH
}
