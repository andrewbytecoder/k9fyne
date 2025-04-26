package kube

type Controller interface {
	Execute(command string) ([]byte, error)
}

type Command []string
