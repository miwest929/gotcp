package main

import (
	"fmt"
	//      "golang.org/x/sys/unix"
	"net"
	"os"
	"syscall"
)

type TCPServer struct {
	port string
}

func NewTCPServer(port string) *TCPServer {
	return &TCPServer{port: port}
}

func (server *TCPServer) createSocketaddr(proto uint16) syscall.Sockaddr {
	ifi, err := net.InterfaceByName("eth0")
	if err != nil {
		fmt.Printf("error retrieving eth0 interface\n")
	}

	return &syscall.SockaddrLinklayer{Protocol: uint16(proto), Ifindex: int(ifi.Index)}
}

func (server *TCPServer) StartAndListen() error {
	fmt.Println("Initializing TCP server...")
	protocol := (syscall.ETH_P_ALL<<8)&0xff00 | syscall.ETH_P_ALL>>8
	fd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, int(protocol))
	defer syscall.Close(fd)
	if err != nil {
		fmt.Printf("error creating socket: %s\n", err)
		os.Exit(1)
	}

	//(fd int, sa Sockaddr) (err error)
	fmt.Println("Waiting for incoming packets...")
	syscall.Bind(fd, server.createSocketaddr(uint16(protocol)))

	var buffer []byte
	for {
		n, _, e := syscall.Recvfrom(fd, buffer, syscall.MSG_PEEK)
		if e != nil {
			fmt.Printf("recvfrom failed: %s\n", err)
			os.Exit(1)
		}
		fmt.Printf("Received %d bytes\n", n)
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: tcp-server <port>")
		os.Exit(1)
	}
	portArg := os.Args[1]

	server := NewTCPServer(portArg)
	server.StartAndListen()
}
