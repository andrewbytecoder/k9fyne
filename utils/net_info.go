package utils

import (
	"fmt"
	"net"
	"strings"
)

// ParseAddress 解析地址字符串，返回主机和端口
func ParseAddress(addr string) (host string, port int, err error) {
	// 分割地址和端口
	if strings.Contains(addr, ":") {
		host, portStr, splitErr := net.SplitHostPort(addr)
		if splitErr != nil {
			return "", 0, fmt.Errorf("无效的地址格式: %v", splitErr)
		}

		// 尝试将端口转换为整数
		port, convErr := net.LookupPort("tcp", portStr)
		if convErr != nil {
			return "", 0, fmt.Errorf("无效的端口号: %v", convErr)
		}
		return host, port, nil
	}

	// 如果没有端口，返回主机名/IP 和默认端口 0
	return addr, 0, nil
}
