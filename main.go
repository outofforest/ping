package main

import (
	"fmt"
	"net"
	"syscall"

	"github.com/pkg/errors"
)

func main() {
	const (
		eth0 = "enp1s0f0np0"
		eth1 = "enp162s0f0np0"
	)

	send := []byte{
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0x00, 0xe0, 0xed, 0xe4, 0xfc, 0xb2,
		0x00, 0x2e,
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17,
		0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
		0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27,
		0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d,
	}
	receive := make([]byte, len(send)+100)

	fd0, err := Open(eth0)
	if err != nil {
		panic(err)
	}

	fd1, err := Open(eth1)
	if err != nil {
		panic(err)
	}

	if _, err := Write(fd0, send); err != nil {
		panic(err)
	}

	n, err := Read(fd1, receive)
	if err != nil {
		panic(err)
	}

	if err := Close(fd0); err != nil {
		panic(err)
	}

	if err := Close(fd1); err != nil {
		panic(err)
	}

	fmt.Println(receive[:n])
}

// Open opens the socket.
func Open(ifName string) (int, error) {
	ethPAll := (uint16(syscall.ETH_P_ALL) << 8) | (uint16(syscall.ETH_P_ALL) >> 8)
	fd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, int(ethPAll))
	if err != nil {
		return 0, errors.WithStack(err)
	}

	ifInfo, err := net.InterfaceByName(ifName)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	var haddr [8]byte
	copy(haddr[:], ifInfo.HardwareAddr)
	addr := syscall.SockaddrLinklayer{
		Protocol: ethPAll,
		Ifindex:  ifInfo.Index,
		Halen:    uint8(len(haddr)),
		Addr:     haddr,
	}

	if err := syscall.Bind(fd, &addr); err != nil {
		return 0, errors.WithStack(err)
	}

	//nolint:staticcheck
	if err := syscall.SetLsfPromisc(ifName, false); err != nil {
		return 0, errors.WithStack(err)
	}

	return fd, nil
}

// Write writes to the socket.
func Write(fd int, packet []byte) (int, error) {
	n, err := syscall.Write(fd, packet)
	if err != nil {
		return 0, errors.WithStack(err)
	}
	return n, nil
}

// Read reads from the socket.
func Read(fd int, packet []byte) (int, error) {
	n, err := syscall.Read(fd, packet)
	if err != nil {
		return 0, errors.WithStack(err)
	}
	return n, nil
}

// Close closes the socket.
func Close(fd int) error {
	return errors.WithStack(syscall.Close(fd))
}
