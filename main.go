package nbsgetname

import (
	"bytes"
	"encoding/binary"
	"log"
	"math/rand"
	"net"
	"syscall"
	"time"
)

const (
	NBSS_PORT   = 137
	NBNS_PORT   = 137
	NBSS_HEADER = 4
)

type Buffer struct {
	data  []byte
	start int
}

func (b *Buffer) PrependBytes(n int) []byte {
	length := cap(b.data) + n
	newData := make([]byte, length)
	copy(newData, b.data)
	b.start = cap(b.data)
	b.data = newData
	return b.data[b.start:]
}

func NewBuffer() *Buffer {
	return &Buffer{}
}

// 反转字符串
func Reverse(s string) (result string) {
	for _, v := range s {
		result = string(v) + result
	}
	return
}
func nbns(buffer *Buffer) {
	rand.Seed(time.Now().UnixNano())
	tid := rand.Intn(0x7fff)
	b := buffer.PrependBytes(12)
	binary.BigEndian.PutUint16(b, uint16(tid))        // 0x0000 标识
	binary.BigEndian.PutUint16(b[2:], uint16(0x0010)) // 标识
	binary.BigEndian.PutUint16(b[4:], uint16(1))      // 问题数
	binary.BigEndian.PutUint16(b[6:], uint16(0))      // 资源数
	binary.BigEndian.PutUint16(b[8:], uint16(0))      // 授权资源记录数
	binary.BigEndian.PutUint16(b[10:], uint16(0))     // 额外资源记录数
	// 查询问题
	b = buffer.PrependBytes(1)
	b[0] = 0x20
	b = buffer.PrependBytes(32)
	copy(b, []byte{0x43, 0x4b})
	for i := 2; i < 32; i++ {
		b[i] = 0x41
	}

	b = buffer.PrependBytes(1)
	// terminator
	b[0] = 0
	// type 和 classIn
	b = buffer.PrependBytes(4)
	binary.BigEndian.PutUint16(b, uint16(33))
	binary.BigEndian.PutUint16(b[2:], 1)
}

func GetNetbiosNameFromIp(ip string) string {

	// Replace with the IP address of the remote host
	remoteIP := net.ParseIP(ip)

	// Create a UDP socket
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_UDP)
	if err != nil {
		log.Fatal(err)
	}
	defer syscall.Close(fd)

	buffer := NewBuffer()
	nbns(buffer)
	packet := buffer.data

	// Send the packet to the remote host
	sa := &syscall.SockaddrInet4{Port: NBNS_PORT}
	copy(sa.Addr[:], remoteIP.To4())
	if err := syscall.Sendto(fd, packet, 0, sa); err != nil {
		log.Fatal(err)
	}
	log.Println("Sent NBNS request to", remoteIP)
	// Receive the response
	response := make([]byte, 512) // Adjust buffer size as needed
	n, _, err := syscall.Recvfrom(fd, response, 0)
	if err != nil {
		log.Fatal(err)
	}

	// Parse the response
	if n < NBSS_HEADER+18 {
		log.Fatal("Invalid response length")
	}

	// Parse the name from the response and clean it
	name := bytes.TrimRight(response[NBSS_HEADER+53:], "\x00")
	newName := parseNetBIOSName(name)

	return newName
}

func parseNetBIOSName(nameBytes []byte) string {
	newNamebytes := []byte{}
	for _, b := range nameBytes {
		//remove if value is 20
		if b == 0x20 {
			return string(newNamebytes)
		}
		newNamebytes = append(newNamebytes, b)
	}
	return string(newNamebytes)
}
