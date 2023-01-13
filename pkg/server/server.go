package server

import (
	"errors"
	"fmt"
	"golang.org/x/exp/maps"
	"io"
	"log"
	"net"
)

type Server struct {
	routeMap   map[string]Handler
	listenAddr string
	ln         net.Listener
	quitCh     chan struct{}
	messageCh  chan []byte
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}
	defer ln.Close()
	s.ln = ln

	fmt.Println(maps.Keys(s.routeMap))
	go s.acceptConnections(maps.Keys(s.routeMap)...)
	<-s.quitCh
	close(s.messageCh)

	return nil
}

func (s *Server) Shutdown() error {
	s.quitCh <- struct{}{}
	return nil
}

func (s *Server) acceptConnections(routes ...string) {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Println("EOF")
				return
			}
			log.Println(err)
			continue
		}

		log.Println("new connection accepted:", conn.RemoteAddr())

		for _, route := range routes {
			go s.handle(conn, s.routeMap[route])
		}

	}
}

type Handler func(conn net.Conn) error

func (s *Server) handle(conn net.Conn, handler Handler) {
	defer conn.Close()
	for {
		if err := handler(conn); err != nil {
			log.Println("Server.handler:", err)
			return
		}
	}
}

func (s *Server) RegisterHandler(route string, handler Handler) {
	s.routeMap[route] = handler
}

func NewServer(listenAddr string) *Server {
	return &Server{
		routeMap:   make(map[string]Handler),
		listenAddr: listenAddr,
		ln:         nil,
		quitCh:     make(chan struct{}),
		messageCh:  make(chan []byte),
	}
}
