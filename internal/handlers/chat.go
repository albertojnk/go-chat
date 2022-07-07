package handlers

import (
	"context"
	"encoding/json"
	"net"
	"strings"
	"time"

	"github.com/albertojnk/go-chat/common"
	"github.com/albertojnk/go-chat/internal/core/domains"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type Client struct {
	UserName string
	Messages chan domains.Message
	Address  *net.UDPAddr
	Conn     *net.UDPConn
	WS       *websocket.Conn
	Config
}

type Config struct {
	Port       int
	Address    string
	BufferSize int
}

func (h *NETHandler) ChatWSHandler(ctx *context.Context, c *gin.Context) error {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		common.HandleError(err, "upgrader.Upgrade")
		return err
	}

	client := &Client{
		WS: ws,
		Config: Config{
			BufferSize: 2048,
		},
	}

	defer func(c *Client) {
		msg := domains.Message{
			UserName:    c.UserName,
			MessageType: domains.GOODBYE,
			Time:        time.Now(),
		}

		msgByte, err := json.Marshal(msg)
		common.HandleError(err, "ChatWSHandler defer json.Marshal")

		sendToClient(c, msgByte)
		ws.Close()
	}(client)

	go client.NewUDP()

	for {
		//Read Message from web
		_, message, err := ws.ReadMessage()
		if err != nil {
			common.HandleError(err, "ws.ReadMessage")
			return err
		}

		sendToClient(client, message)
	}

}

func (c *Client) NewUDP() error {
	broadcast, err := net.ResolveUDPAddr("udp", "0.0.0.0:3000")
	if err != nil {
		common.HandleError(err, "NewUDP() net.ResolveUDPAddr")
		return err
	}

	conn, err := net.DialUDP("udp", nil, broadcast)
	if err != nil {
		common.HandleError(err, "NewUDP() net.DialUDP")
		return err
	}
	defer conn.Close()

	// get resolve client addr
	conAddr, err := net.ResolveUDPAddr("udp", conn.LocalAddr().String())
	if err != nil {
		common.HandleError(err, "NewUDP() net.ResolveUDPAddr")
		return err
	}

	c.Conn = conn
	c.Messages = make(chan domains.Message, 20)
	c.Address = conAddr

	go c.sendMessageToWeb()

	for {
		c.handleMessage()
	}
}

func (c *Client) sendMessageToWeb() {
	for {
		msg := <-c.Messages

		err := c.WS.WriteJSON(msg)
		common.HandleError(err, "sendMessage ws.WriteMessage")
	}
}

func (c *Client) handleMessage() {
	message := make([]byte, c.BufferSize)
	rlen, err := c.Conn.Read(message[0:])
	common.HandleError(err, "handleMessage c.Conn.Read")

	data := strings.TrimSpace(string(message[0:rlen]))

	msg := domains.Message{}
	err = json.Unmarshal([]byte(data), &msg)
	common.HandleError(err, "handleMessage json.Unmarshal")

	c.UserName = msg.UserName
	msg.Time = time.Now()
	msg.ID = common.GenerateUUID()

	c.Messages <- msg
}

func sendToClient(c *Client, data []byte) {
	msg := domains.Message{}
	err := json.Unmarshal(data, &msg)
	common.HandleError(err, "handleMessage json.Unmarshal")

	c.UserName = msg.UserName
	msg.Address = c.Address
	msg.Time = time.Now()
	msg.ID = common.GenerateUUID()
	msgByte, err := json.Marshal(msg)
	common.HandleError(err, "handleMessage json.Marshal")

	_, err = c.Conn.Write(msgByte)
	common.HandleError(err, "sendMessage")
}

func (c *Client) Close() {
	c.Conn.Close()
}
