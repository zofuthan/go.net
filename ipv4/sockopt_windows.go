// Copyright 2012 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ipv4

import (
	"net"
	"os"
	"syscall"
	"unsafe"
)

// Please refer to the online manual;
// http://msdn.microsoft.com/en-us/library/windows/desktop/ms738586(v=vs.85).aspx

func ipv4TOS(fd syscall.Handle) (int, error) {
	var v int32
	l := int32(4)
	err := syscall.Getsockopt(fd, ianaProtocolIP, syscall.IP_TOS, (*byte)(unsafe.Pointer(&v)), &l)
	if err != nil {
		return 0, os.NewSyscallError("getsockopt", err)
	}
	return int(v), nil
}

func setIPv4TOS(fd syscall.Handle, v int) error {
	vv := int32(v)
	err := syscall.Setsockopt(fd, ianaProtocolIP, syscall.IP_TOS, (*byte)(unsafe.Pointer(&vv)), 4)
	if err != nil {
		return os.NewSyscallError("setsockopt", err)
	}
	return nil
}

func ipv4TTL(fd syscall.Handle) (int, error) {
	var v int32
	l := int32(4)
	err := syscall.Getsockopt(fd, ianaProtocolIP, syscall.IP_TTL, (*byte)(unsafe.Pointer(&v)), &l)
	if err != nil {
		return 0, os.NewSyscallError("getsockopt", err)
	}
	return int(v), nil
}

func setIPv4TTL(fd syscall.Handle, v int) error {
	vv := int32(v)
	err := syscall.Setsockopt(fd, ianaProtocolIP, syscall.IP_TTL, (*byte)(unsafe.Pointer(&vv)), 4)
	if err != nil {
		return os.NewSyscallError("setsockopt", err)
	}
	return nil
}

func ipv4MulticastTTL(fd syscall.Handle) (int, error) {
	var v int32
	l := int32(4)
	err := syscall.Getsockopt(fd, ianaProtocolIP, syscall.IP_MULTICAST_TTL, (*byte)(unsafe.Pointer(&v)), &l)
	if err != nil {
		return 0, os.NewSyscallError("getsockopt", err)
	}
	return int(v), nil
}

func setIPv4MulticastTTL(fd syscall.Handle, v int) error {
	vv := int32(v)
	err := syscall.Setsockopt(fd, ianaProtocolIP, syscall.IP_MULTICAST_TTL, (*byte)(unsafe.Pointer(&vv)), 4)
	if err != nil {
		return os.NewSyscallError("setsockopt", err)
	}
	return nil
}

func ipv4ReceiveTTL(fd syscall.Handle) (bool, error) {
	// NOTE: Not supported yet on any Windows
	return false, syscall.EWINDOWS
}

func setIPv4ReceiveTTL(fd syscall.Handle, v bool) error {
	// NOTE: Not supported yet on any Windows
	return syscall.EWINDOWS
}

func ipv4ReceiveDestinationAddress(fd syscall.Handle) (bool, error) {
	// TODO(mikio): Implement this for XP and beyond
	return false, syscall.EWINDOWS
}

func setIPv4ReceiveDestinationAddress(fd syscall.Handle, v bool) error {
	// TODO(mikio): Implement this for XP and beyond
	return syscall.EWINDOWS
}

func ipv4HeaderPrepend(fd syscall.Handle) (bool, error) {
	// TODO(mikio): Implement this for XP and beyond
	return false, syscall.EWINDOWS
}

func setIPv4HeaderPrepend(fd syscall.Handle, v bool) error {
	// TODO(mikio): Implement this for XP and beyond
	return syscall.EWINDOWS
}

func ipv4ReceiveInterface(fd syscall.Handle) (bool, error) {
	// TODO(mikio): Implement this for Vista and beyond
	return false, syscall.EWINDOWS
}

func setIPv4ReceiveInterface(fd syscall.Handle, v bool) error {
	// TODO(mikio): Implement this for Vista and beyond
	return syscall.EWINDOWS
}

func ipv4MulticastInterface(fd syscall.Handle) (*net.Interface, error) {
	var a [4]byte
	l := int32(4)
	err := syscall.Getsockopt(fd, ianaProtocolIP, syscall.IP_MULTICAST_IF, (*byte)(unsafe.Pointer(&a[0])), &l)
	if err != nil {
		return nil, os.NewSyscallError("getsockopt", err)
	}
	return netIP4ToInterface(net.IPv4(a[0], a[1], a[2], a[3]))
}

func setIPv4MulticastInterface(fd syscall.Handle, ifi *net.Interface) error {
	ip, err := netInterfaceToIP4(ifi)
	if err != nil {
		return os.NewSyscallError("setsockopt", err)
	}
	var a [4]byte
	copy(a[:], ip.To4())
	err = syscall.Setsockopt(fd, ianaProtocolIP, syscall.IP_MULTICAST_IF, (*byte)(unsafe.Pointer(&a[0])), 4)
	if err != nil {
		return os.NewSyscallError("setsockopt", err)
	}
	return nil
}

func ipv4MulticastLoopback(fd syscall.Handle) (bool, error) {
	var v int32
	l := int32(4)
	err := syscall.Getsockopt(fd, ianaProtocolIP, syscall.IP_MULTICAST_LOOP, (*byte)(unsafe.Pointer(&v)), &l)
	if err != nil {
		return false, os.NewSyscallError("getsockopt", err)
	}
	return v == 1, nil
}

func setIPv4MulticastLoopback(fd syscall.Handle, v bool) error {
	vv := int32(boolint(v))
	err := syscall.Setsockopt(fd, ianaProtocolIP, syscall.IP_MULTICAST_LOOP, (*byte)(unsafe.Pointer(&vv)), 4)
	if err != nil {
		return os.NewSyscallError("setsockopt", err)
	}
	return nil
}

func joinIPv4Group(fd syscall.Handle, ifi *net.Interface, grp net.IP) error {
	mreq := &syscall.IPMreq{Multiaddr: [4]byte{grp[0], grp[1], grp[2], grp[3]}}
	if err := setSyscallIPMreq(mreq, ifi); err != nil {
		return err
	}
	err := syscall.Setsockopt(fd, ianaProtocolIP, syscall.IP_ADD_MEMBERSHIP, (*byte)(unsafe.Pointer(mreq)), int32(unsafe.Sizeof(*mreq)))
	if err != nil {
		return os.NewSyscallError("setsockopt", err)
	}
	return nil
}

func leaveIPv4Group(fd syscall.Handle, ifi *net.Interface, grp net.IP) error {
	mreq := &syscall.IPMreq{Multiaddr: [4]byte{grp[0], grp[1], grp[2], grp[3]}}
	if err := setSyscallIPMreq(mreq, ifi); err != nil {
		return err
	}
	err := syscall.Setsockopt(fd, ianaProtocolIP, syscall.IP_DROP_MEMBERSHIP, (*byte)(unsafe.Pointer(mreq)), int32(unsafe.Sizeof(*mreq)))
	if err != nil {
		return os.NewSyscallError("setsockopt", err)
	}
	return nil
}
