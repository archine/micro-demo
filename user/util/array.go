package util

import "google.golang.org/grpc/resolver"

// Exists 判断集合中是否存在某个元素
func Exists(l []resolver.Address, addr string) bool {
	for _, address := range l {
		if address.Addr == addr {
			return true
		}
	}
	return false
}

func Remove(s []resolver.Address, serverName string) ([]resolver.Address, bool) {
	for i, address := range s {
		if address.ServerName == serverName {
			if i == len(s)-1 {
				return s[:i], true
			}
			if i == 0 {
				return s[i+1:], true
			}
			return append(s[:i], s[i+1:]...), true
		}
	}
	return nil, false
}
