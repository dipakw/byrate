//go:build !windows

package app

import "syscall"

func setSockOptIPv6Only(fd uintptr) error {
	return syscall.SetsockoptInt(int(fd), syscall.IPPROTO_IPV6, syscall.IPV6_V6ONLY, 0)
}
