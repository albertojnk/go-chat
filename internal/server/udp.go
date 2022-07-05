package server

import (
	"fmt"
	"log"
	"net"
	"strings"

	_cache "github.com/albertojnk/go-chat/internal/cache"
	"github.com/albertojnk/go-chat/internal/client"
	"github.com/google/uuid"
)

type UDPServer struct {
	Conn     *net.UDPConn
	Messages chan string
	Clients  map[string]client.Client
	Config
	cache *_cache.Redis
}

type Config struct {
	Port       int
	Address    string
	BufferSize int
}

func (s *UDPServer) NewUDP() {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		Port: s.Config.Port,
		IP:   net.ParseIP(s.Config.Address),
	})
	if err != nil {
		log.Fatalf("error while connecting udp server: %v, err: %v", conn, err)
		panic(err)
	}
	defer conn.Close()

	s.Conn = conn

	log.Printf("server listening %s", s.Conn.LocalAddr().String())

	s.Messages = make(chan string, 20)
	s.Clients = make(map[string]client.Client, 0)

	go s.sendMessage()

	for {
		s.handleMessage()
	}
}

func (s *UDPServer) sendMessage() {
	for {
		msg := <-s.Messages
		for _, c := range s.Clients {
			_, err := s.Conn.WriteToUDP([]byte(msg), c.Address)
			if err != nil {
				log.Fatalf("error while writing to udp: %v, err: %v", s.Conn, err)
				panic(err)
			}
		}
	}

}
func (s *UDPServer) handleMessage() {
	message := make([]byte, s.BufferSize)
	rlen, remote, err := s.Conn.ReadFromUDP(message[:])
	if err != nil {
		log.Fatalf("error while reading from udp: %v, err: %v", s.Conn, err)
		panic(err)
	}

	data := strings.TrimSpace(string(message[:rlen]))
	if data == "connected" {
		s.Clients[remote.String()] = client.Client{
			ID:      uuid.New(),
			Name:    fmt.Sprintf("%v", remote.Port),
			Address: remote,
		}
	} else {
		s.Messages <- data
		fmt.Printf("received: %s from %s\n", data, remote)
	}
}
