package handlers

import (
	"encoding/json"
	"testing"

	"github.com/albertojnk/go-chat/internal/cache"
	"github.com/albertojnk/go-chat/internal/core/domains"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func getHandlers(t *testing.T) *NETHandler {
	cacheConfig := cache.NewRedisConfig()
	cache, err := cache.NewRedis(cacheConfig)
	if err != nil {
		t.Fatal(err)
	}

	hdl := NewNETHandler(cache)
	return hdl
}

func TestNETHandler_ChatWSHandler(t *testing.T) {
	hdl := getHandlers(t)

	router := gin.New()
	router.GET("/chat", HandlerAPI(hdl.ChatWSHandler))

	testValues := []domains.Message{
		{
			MessageType: domains.HANDSHAKE,
			UserName:    "foo",
			Content:     "handshake from test",
		},
		{
			MessageType: domains.MESSAGE,
			UserName:    "foo",
			Content:     "message from test",
		},
		{
			MessageType: domains.MESSAGE,
			UserName:    "foo",
			Content:     "message from test 2",
		},
		{
			MessageType: domains.MESSAGE,
			UserName:    "notconnected",
			Content:     "message from test 3",
		},
		{
			ID:          "2",
			MessageType: domains.DELETEMESSAGE,
			UserName:    "bar",
			Content:     "",
		},
		{
			ID:          "2",
			MessageType: domains.DELETEMESSAGE,
			UserName:    "foo",
			Content:     "",
		},
		{
			ID:          "4",
			MessageType: domains.GOODBYE,
			UserName:    "foo",
			Content:     "",
		},
	}

	ws, _, err := websocket.DefaultDialer.Dial("ws://localhost:9000/chat", nil)
	if err != nil {
		t.Fatal(err)
	}
	defer ws.Close()

	for _, values := range testValues {

		msgByte, _ := json.Marshal(values)
		ws.WriteMessage(websocket.TextMessage, msgByte)

		_, m, _ := ws.ReadMessage()

		msg := domains.Message{}
		json.Unmarshal(m, &msg)

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
				if msg.Content != values.Content || msg.ID == values.ID || msg.UserName != values.UserName {
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
			if values.UserName == "foo" && values.ID == "2" && msg.MessageType != domains.INVALIDMESSAGE {
				t.Errorf("Should've got %v but got %v", domains.DELETEMESSAGE, msg.MessageType)
			}
		}

		if values.MessageType == domains.GOODBYE {
			if msg.Content != values.Content || msg.ID == values.ID || msg.UserName != values.UserName {
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

		if msg.UserName != values.UserName {
			t.Errorf("Should've got %v but got %v", values.UserName, msg.UserName)
		}

	}
}
