package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"

	_cache "github.com/albertojnk/go-chat/internal/cache"
	"github.com/albertojnk/go-chat/internal/core/domains"
)

type UDPServer struct {
	Conn     *net.UDPConn
	Messages chan []byte
	Clients  map[string]domains.Client
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
	s.Messages = make(chan []byte, 20)
	s.Clients = make(map[string]domains.Client, 0)

	go s.sendMessage()

	for {
		s.handleMessage()
	}
}

func (s *UDPServer) sendMessage() {
	for {
		data := <-s.Messages

		msg := domains.Message{}
		err := json.Unmarshal(data, &msg)
		if err != nil {
			log.Fatalf("error while decrypting message: %v, err: %v", msg, err)
			panic(err)
		}

		for _, c := range s.Clients {
			if c.UserName != msg.UserName {
				_, err := s.Conn.WriteToUDP(data, c.Address)
				if err != nil {
					log.Fatalf("error while writing to udp: %v, err: %v", s.Conn, err)
					panic(err)
				}
			}
		}
	}

}
func (s *UDPServer) handleMessage() {
	message := make([]byte, s.BufferSize)
	rlen, _, err := s.Conn.ReadFromUDP(message[0:])
	if err != nil {
		log.Fatalf("error while reading from udp: %v, err: %v", s.Conn, err)
		panic(err)
	}

	data := strings.TrimSpace(string(message[0:rlen]))
	fmt.Println(data)
	msg := domains.Message{}
	err = json.Unmarshal([]byte(data), &msg)
	if err != nil {
		log.Fatalf("error while decrypting message: %v, err: %v", msg, err)
		panic(err)
	}

	switch msg.MessageType {
	case domains.HANDSHAKE:
		s.Clients[msg.UserName] = domains.Client{
			UserName: msg.UserName,
			Address:  msg.Address,
		}
	case domains.MESSAGE:
		s.Messages <- message[:rlen]
		fmt.Printf("[%v] %s: %s \n", msg.Time.Format("15:04:05"), msg.UserName, msg.Content)
	}

}
