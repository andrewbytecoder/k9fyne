package config

import "github.com/andrewbytecoder/k9fyne/controller"

type Cfg struct {
	SSH *controller.SSHClient
}

func NewConfig() *Cfg {
	return &Cfg{
		SSH: controller.NewSSHClient(),
	}
}

func (c *Cfg) GetSSH() *controller.SSHClient {
	return c.SSH
}
