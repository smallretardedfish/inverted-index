package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	utils "github.com/smallretardedfish/inverted-index/pkg"
	"github.com/smallretardedfish/inverted-index/pkg/inverted_index"
	"github.com/smallretardedfish/inverted-index/pkg/server"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func run() error {
	args := make([]string, 2)
	copy(args, os.Args[1:])

	if len(args) == 0 {
		log.Println("no args given")
		args = append(args, "1", ".")
	} else if len(args) == 1 {
		args = append(args, ".")
	}

	num := args[0]
	n, err := strconv.Atoi(num)
	if err != nil {
		return fmt.Errorf("number of workers must be integer: %w", err)
	}

	dir := args[1]
	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("can't read from given dir: %s, error: %w", dir, err)
	}

	var fileSources []string
	for _, entry := range dirEntries {
		fileSources = append(fileSources, filepath.Join(dir, entry.Name()))
	}

	fmt.Printf("number of workers:%d\ndirectory with files to be processed:%s\n", n, dir)
	invIndex := inverted_index.NewMapInvertedIndex(inverted_index.FileSourceType)

	var (
		e error
	)

	t := utils.EstimateExecutionTime(func() {
		if err := invIndex.Build(n, fileSources); err != nil {
			e = err
		}
	})
	if e != nil {
		return e
	}
	log.Println(t)
	s := server.NewServer(":8000")
	invertedIndexServer := invIndexServer{
		invIndex:  invIndex,
		tcpServer: s,
	}

	log.Println("Starting Server...")
	if err := invertedIndexServer.Start(); err != nil {
		return err
	}

	return nil
}

type invIndexServer struct {
	invIndex  inverted_index.InvertedIndex
	tcpServer *server.Server
}

func (s *invIndexServer) Start() error {
	s.tcpServer.RegisterHandler("search", func(conn net.Conn) error {
		var size int64
		if err := binary.Read(conn, binary.LittleEndian, &size); err != nil {
			return err
		}
		buff := &bytes.Buffer{}
		n, err := io.CopyN(buff, conn, size)
		if err != nil {
			return err
		}
		word := string(buff.Bytes()[:n])
		res := s.invIndex.Search(word)

		response := "not found\n"

		if len(res) != 0 {
			response = fmt.Sprintf("%s: %s\n", word, strings.Join(res, ","))
		}
		respBytes := []byte(response)
		if err := binary.Write(conn, binary.LittleEndian, int64(len(respBytes))); err != nil {
			return err
		}

		if _, err := io.CopyN(conn, strings.NewReader(response), int64(len(respBytes))); err != nil {
			return err
		}

		return nil
	})

	return s.tcpServer.Start()
}

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}
