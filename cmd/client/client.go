package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
		log.Println(err)
		return
	}
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Printf(">> ")
	for scanner.Scan() {
		word := scanner.Text()
		wordBytes := []byte(word)
		size := len([]byte(word))
		err = binary.Write(conn, binary.LittleEndian, int64(size))
		if err != nil {
			log.Println(err)
			return
		}
		err = binary.Write(conn, binary.LittleEndian, wordBytes)
		if err != nil {
			log.Println(err)
			return
		}

		var respSize int64
		if err := binary.Read(conn, binary.LittleEndian, &respSize); err != nil {
			return
		}

		//reading response
		buff := &bytes.Buffer{}
		n, err := io.CopyN(buff, conn, respSize)
		if err != nil {
			return
		}
		fmt.Println(string(buff.Bytes()[:n]))

		fmt.Printf(">> ")
	}
}
