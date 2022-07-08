package server

import (
	"encoding/json"
	"net"
	"strings"
	"testing"

	"github.com/albertojnk/go-chat/common"
	"github.com/albertojnk/go-chat/internal/cache"
	"github.com/albertojnk/go-chat/internal/core/domains"
	"github.com/joho/godotenv"
)

func getServer(t *testing.T) *UDPServer {
	godotenv.Load()

	cacheConfig := cache.NewRedisConfig()
	cache, err := cache.NewRedis(cacheConfig)
	if err != nil {
		t.Fatal(err)
	}

	cache.FlushDB()

	s := UDPServer{
		Config: Config{
			Port:       3000,
			Address:    "0.0.0.0",
			BufferSize: 2048,
		},
		Cache: cache,
	}
	return &s
}

func setup(t *testing.T) *net.UDPConn {
	broadcast, err := net.ResolveUDPAddr("udp", common.GetEnv("UDP_SERVER_CONN_ADDR", "0.0.0.0:3000"))
	if err != nil {
		t.Fatal(err)
	}

	udpClient, err := net.DialUDP("udp", nil, broadcast)
	if err != nil {
		t.Fatal(err)
	}

	return udpClient
}

func TestUDPServer_NewUDP(t *testing.T) {
	server := getServer(t)
	go server.NewUDP()

	udpClient := setup(t)
	conAddr, err := net.ResolveUDPAddr("udp", udpClient.LocalAddr().String())
	if err != nil {
		t.Fatal(err)
	}

	testValues := []domains.Message{
		{
			MessageType: domains.HANDSHAKE,
			UserName:    "foo",
			Content:     "handshake from test",
			Address:     conAddr,
		},
		{
			ID:          "1",
			MessageType: domains.MESSAGE,
			UserName:    "foo",
			Content:     "message from test",
			Address:     conAddr,
		},
		{
			ID:          "2",
			MessageType: domains.MESSAGE,
			UserName:    "foo",
			Content:     "message from test 2",
			Address:     conAddr,
		},
		{
			ID:          "3",
			MessageType: domains.MESSAGE,
			UserName:    "notconnected",
			Content:     "message from test 3",
			Address:     conAddr,
		},
		{
			ID:          "2",
			MessageType: domains.DELETEMESSAGE,
			UserName:    "bar",
			Content:     "",
			Address:     conAddr,
		},
		{
			ID:          "2",
			MessageType: domains.DELETEMESSAGE,
			UserName:    "foo",
			Content:     "",
			Address:     conAddr,
		},
		{
			ID:          "4",
			MessageType: domains.GOODBYE,
			UserName:    "foo",
			Content:     "",
			Address:     conAddr,
		},
	}

	for _, values := range testValues {

		byteMsg, _ := json.Marshal(values)
		udpClient.Write(byteMsg)

		for {
			message := make([]byte, 2048)
			rlen, _, _ := udpClient.ReadFromUDP(message[0:])

			data := strings.TrimSpace(string(message[0:rlen]))

			msg := domains.Message{}
			json.Unmarshal([]byte(data), &msg)

			if values.MessageType == domains.HANDSHAKE {
				if msg.Content != "handshake from server" {
					t.Errorf("Should've got 'handshake from server' but got %v", msg.Content)
				}
			}

			if values.MessageType == domains.MESSAGE {
				if values.UserName == "notconnected" {
					if msg.MessageType != domains.INVALIDMESSAGE {
						t.Errorf("Should've got %v but got %v", domains.INVALIDMESSAGE, msg.MessageType)
					}
					if msg.Content != "invalid message" {
						t.Errorf("Should've got 'invalid message' but got %v", msg.Content)
					}
				} else {
					if msg.Content != values.Content || msg.ID != values.ID || msg.UserName != values.UserName {
						t.Errorf("Should've got %v but got %v",
							domains.Message{
								Content:  values.Content,
								UserName: values.UserName,
								ID:       values.ID,
							},
							domains.Message{
								Content:  msg.Content,
								UserName: msg.UserName,
								ID:       msg.ID,
							},
						)
					}
				}
			}

			if values.MessageType == domains.DELETEMESSAGE {
				if values.UserName == "bar" && values.ID == "2" && msg.MessageType != domains.INVALIDMESSAGE {
					t.Errorf("Should've got %v but got %v", domains.INVALIDMESSAGE, msg.MessageType)
				}
				if values.UserName == "foo" && values.ID == "2" && msg.MessageType != domains.DELETEMESSAGE {
					t.Errorf("Should've got %v but got %v", domains.DELETEMESSAGE, msg.MessageType)
				}
			}

			if values.MessageType == domains.GOODBYE {
				// as it is if this triggers there is no one connected
				users, _ := server.Cache.GetValue(domains.USERSKEY)
				if users != "" {
					t.Errorf("Should've got nothing but got %v", users)
				}

				messages, _ := server.Cache.GetValue(domains.MESSAGESKEY)
				if messages != "" {
					t.Errorf("Should've got nothing but got %v", messages)
				}
			}

			if msg.UserName != values.UserName {
				t.Errorf("Should've got %v but got %v", values.UserName, msg.UserName)
			}

			break
		}

	}

}
