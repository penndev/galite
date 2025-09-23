package api

import (
	"fmt"
	"net"
	"net/url"
)

func replaceDomainWithIP(originalURL string) (string, error) {
	// 解析URL
	parsedURL, err := url.Parse(originalURL)
	if err != nil {
		return "", fmt.Errorf("解析URL失败: %v", err)
	}

	// 如果已经是IP地址，直接返回原URL
	if net.ParseIP(parsedURL.Hostname()) != nil {
		return originalURL, nil
	}

	// 解析域名获取IP地址
	ips, err := net.LookupIP(parsedURL.Hostname())
	if err != nil {
		return "", fmt.Errorf("域名解析失败: %v", err)
	}

	if len(ips) == 0 {
		return "", fmt.Errorf("未找到域名的IP地址")
	}

	// 优先使用IPv4地址
	var ip net.IP
	for _, candidate := range ips {
		if candidate.To4() != nil {
			ip = candidate
			break
		}
	}

	// 如果没有IPv4地址，使用第一个找到的地址
	if ip == nil {
		ip = ips[0]
	}

	// 重建URL，替换主机部分为IP地址
	// 保留端口号
	host := ip.String()
	if parsedURL.Port() != "" {
		host = net.JoinHostPort(host, parsedURL.Port())
	}

	parsedURL.Host = host
	return parsedURL.String(), nil
}
