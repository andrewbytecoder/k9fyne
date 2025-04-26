package kube

type KubeController struct {
	Command
}

func NewKubeController() *KubeController {
	return &KubeController{
		Command: Command{"kubectl"},
	}
}

func (c *KubeController) WithNamespace(namespace string) *KubeController {
	c.Command = append(c.Command, "-n", namespace)
	return c
}

func (c *KubeController) WithCommonds(commands ...string) *KubeController {
	c.Command = append(c.Command, commands...)
	return c
}
func (c *KubeController) WithWide() *KubeController {
	c.Command = append(c.Command, "-o", "wide")
	return c
}
