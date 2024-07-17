package resp

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

type Server struct {
	host     string
	port     string
	listener net.Listener
	handler  CommandHandler
	quit     chan struct{}
}

func NewServer(host, port string) *Server {
	return &Server{
		host: host,
		port: port,
		quit: make(chan struct{}),
	}
}

func (s *Server) Handler(handler CommandHandler) {
	s.handler = handler
}

func (s *Server) Address() string {
	return fmt.Sprintf("%s:%s", s.host, s.port)
}

func (s *Server) Stop() {
	log.Println("[server] shutting down...")
	close(s.quit)
	s.listener.Close()
}

func (s *Server) ServeAndListen() error {
	var err error

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		s.Stop()
	}()

	s.listener, err = net.Listen("tcp", s.Address())
	if err != nil {
		return err
	}
	defer s.listener.Close()

	log.Println("[server] listening requests on: ", s.Address())

loop:
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-s.quit:
				break loop
			default:
				log.Printf("[server] failed to accept connection from %s: %s\n", conn.RemoteAddr(), err)
			}
		}

		go s.handle(conn)
	}

	return nil
}

func (s *Server) handle(conn net.Conn) {
	buf := make([]byte, 4096)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			continue
		}

		readBuf := buf[:n]
		splitString := strings.Split(string(readBuf), "*")

		var valNodes []ValueNode

		for i := 1; i < len(splitString); i++ {
			s := fmt.Sprintf("*%s", splitString[i])
			sr := strings.NewReader(s)
			valNode := parse(sr)
			valNodes = append(valNodes, valNode)
		}

		var responses []ValueNode

	outer:
		for _, valNode := range valNodes {
			if valNode.types != ValueNodeTypeArray {
				responses = append(responses, ValueNode{
					types: ValueNodeTypeSimpleError,
					val:   fmt.Sprintf("ERROR: expected * at the begining, got %s", valNode.types),
				})
				break
			}

			cmdStrs := []string{}

			for _, node := range valNode.nodes {
				if node.types != ValueNodeTypeBulkString {
					responses = append(responses, ValueNode{
						types: ValueNodeTypeSimpleError,
						val:   fmt.Sprintf("ERROR: expected $ at the begining, got %s\r\n", node.types),
					})
					break outer
				}

				cmdStrs = append(cmdStrs, node.val)
			}

			cmd := NewCommand(strings.ToLower(cmdStrs[0]), cmdStrs[1:]...)

			response := s.handler.Serve(cmd)

			responses = append(responses, response)
		}

		response := bytes.Buffer{}
		for _, r := range responses {
			response.Write(r.Marshal())
		}

		conn.Write(response.Bytes())

	}

}
