package kube

import (
	"fmt"
	"strings"
)

// Pod kubectl get pod -A -o wide
type Pod struct {
	Namespace      string `json:"NAMESPACE"`
	Name           string `json:"NAME"`
	Ready          string `json:"READY"`
	Status         string `json:"STATUS"`
	Restarts       string `json:"RESTARTS"`
	Age            string `json:"AGE"`
	Ip             string `json:"IP"`
	Node           string `json:"NODE"`
	NominatedNode  string `json:"NOMINATED NODE"`
	ReadinessGates string `json:"READINESS GATES"`
}

type PodsInfo struct {
	PodMap          map[string]*Pod
	Namespaces2Pods map[string][]*Pod
}

func NewPodsInfo() *PodsInfo {
	return &PodsInfo{
		PodMap:          make(map[string]*Pod),
		Namespaces2Pods: make(map[string][]*Pod),
	}
}

func parsePods(input string) ([]*Pod, error) {
	var pods []*Pod

	// 按行分割输入数据
	lines := strings.Split(input, "\n")
	if len(lines) < 2 {
		return nil, fmt.Errorf("invalid input format")
	}

	// 获取表头
	headers := strings.Fields(lines[0])

	// 遍历每一行数据（跳过表头）
	for _, line := range lines[1:] {
		if strings.TrimSpace(line) == "" {
			continue // 跳过空行
		}

		// 分割字段
		fields := strings.Fields(line)
		if len(fields) != len(headers) {
			return nil, fmt.Errorf("mismatched fields in line: %s", line)
		}

		// 映射字段到结构体
		pod := &Pod{
			Namespace:      fields[0],
			Name:           fields[1],
			Ready:          fields[2],
			Status:         fields[3],
			Restarts:       fields[4],
			Age:            fields[5],
			Ip:             fields[6],
			Node:           fields[7],
			NominatedNode:  fields[8],
			ReadinessGates: fields[9],
		}

		// 添加到结果列表
		pods = append(pods, pod)
	}

	return pods, nil
}

func (p *PodsInfo) ParsePods(input string) error {
	pods, err := parsePods(input)
	if err != nil {
		return err
	}
	for _, pod := range pods {
		p.PodMap[pod.Name] = pod
		p.Namespaces2Pods[pod.Namespace] = append(p.Namespaces2Pods[pod.Namespace], pod)
	}
	
	return nil
}
