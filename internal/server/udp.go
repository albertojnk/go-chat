package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/albertojnk/go-chat/common"
	_cache "github.com/albertojnk/go-chat/internal/cache"
	"github.com/albertojnk/go-chat/internal/core/domains"
	"github.com/go-redis/redis"
)

type UDPServer struct {
	Conn     *net.UDPConn
	Messages chan []byte
	Clients  map[string]domains.Client
	Cache    *_cache.Redis
	Config
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
			_, err := s.Conn.WriteToUDP(data, c.Address)
			if err != nil {
				err := s.Cache.DeleteKeys(msg.UserName)
				common.HandleError(err, "sendMessage s.Cache.DeleteKeys")
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

	data := strings.TrimSpace(string(message[:rlen]))
	fmt.Println(data)
	msg := domains.Message{}
	err = json.Unmarshal([]byte(data), &msg)
	if err != nil {
		log.Fatalf("error while decrypting message: %v, err: %v", msg, err)
		panic(err)
	}

	switch msg.MessageType {
	case domains.HANDSHAKE:
		// set new user to redis, proliferate msg
		client := domains.Client{
			UserName: msg.UserName,
			Address:  msg.Address,
		}
		users, err := s.getUsers()
		if err != redis.Nil {
			common.HandleError(err, "HANDSHAKE s.Cache.GetUsers")
		}

		s.Clients[msg.UserName] = client
		users[msg.UserName] = client

		bUsers, err := json.Marshal(users)
		common.HandleError(err, "MESSAGE Marshal")

		err = s.Cache.SetValue(domains.USERSKEY, string(bUsers), domains.CACHEDURATION)
		common.HandleError(err, "HANDSHAKE s.Cache.SetValue")

	case domains.MESSAGE:
		// add message to redis
		msgs, err := s.getMessages()
		if err != redis.Nil {
			common.HandleError(err, "MESSAGE s.getMessages")
		}

		if len(msgs) == 20 {
			msgs = msgs[1:]
			msgs = append(msgs, msg)
		}

		bMsgs, err := json.Marshal(msgs)
		common.HandleError(err, "MESSAGE Marshal")

		err = s.Cache.SetValue(domains.MESSAGESKEY, string(bMsgs), domains.CACHEDURATION)
		common.HandleError(err, "MESSAGE s.Cache.SetValue")

	case domains.DELETEMESSAGE:
		// delete message from redis
		msgs, err := s.deleteMessage(msg.ID)
		common.HandleError(err, "DELETEMESSAGE s.deleteMessage")

		bMsgs, err := json.Marshal(msgs)
		common.HandleError(err, "DELETEMESSAGE Marshal")

		err = s.Cache.SetValue(domains.MESSAGESKEY, string(bMsgs), domains.CACHEDURATION)
		common.HandleError(err, "DELETEMESSAGE s.Cache.SetValue")

	case domains.GOODBYE:
		users, err := s.deleteUser(msg.UserName)
		if err != redis.Nil {
			common.HandleError(err, "GOODBYE s.deleteUser")
		}

		msg.Address = nil

		bUsers, err := json.Marshal(users)
		common.HandleError(err, "GOODBYE Marshal")

		err = s.Cache.SetValue(domains.MESSAGESKEY, string(bUsers), domains.CACHEDURATION)
		common.HandleError(err, "GOODBYE s.Cache.SetValue")

		delete(s.Clients, msg.UserName)
	}

	s.Messages <- message[:rlen]
}

func (s *UDPServer) populateFromRedis() {

	// messages
	{
	}
	// users
	{
	}
}

func (s *UDPServer) getUsers() (map[string]domains.Client, error) {
	users := map[string]domains.Client{}
	val, err := s.Cache.GetValue(domains.USERSKEY)
	if err != nil {
		return users, err
	}

	err = json.Unmarshal([]byte(val), &users)
	if err != nil {
		return users, err
	}

	return users, err
}

func (s *UDPServer) getMessages() ([]domains.Message, error) {
	msgs := []domains.Message{}
	val, err := s.Cache.GetValue(domains.MESSAGESKEY)
	if err != nil {
		return msgs, err
	}

	err = json.Unmarshal([]byte(val), &msgs)
	if err != nil {
		return msgs, err
	}

	return msgs, err
}

func (s *UDPServer) deleteMessage(id string) ([]domains.Message, error) {
	msgs := []domains.Message{}
	val, err := s.Cache.GetValue(domains.MESSAGESKEY)
	if err != nil {
		return msgs, err
	}

	err = json.Unmarshal([]byte(val), &msgs)
	if err != nil {
		return msgs, err
	}

	newMsgs := []domains.Message{}
	for index, msg := range msgs {
		if msg.ID == id {
			newMsgs = append(msgs[:index], msgs[index+1:]...)
		}
	}

	if len(newMsgs) > 0 {
		return newMsgs, err
	}

	return msgs, err
}

func (s *UDPServer) deleteUser(username string) (map[string]domains.Client, error) {
	clients := map[string]domains.Client{}
	val, err := s.Cache.GetValue(domains.USERSKEY)
	if err != nil {
		return clients, err
	}

	err = json.Unmarshal([]byte(val), &clients)
	if err != nil {
		return clients, err
	}

	delete(clients, username)

	return clients, err
}
