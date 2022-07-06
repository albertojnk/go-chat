package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/albertojnk/go-chat/common"
	"github.com/albertojnk/go-chat/internal/cache"
	"github.com/albertojnk/go-chat/internal/core/domains"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func (h *NETHandler) ChatWSHandler(ctx *context.Context, c *gin.Context) error {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		common.HandleError(err, "upgrader.Upgrade")
		return err
	}
	defer ws.Close()

	client := &Client{}
	go client.NewUDP()
	for {
		//Read Message from client
		mt, message, err := ws.ReadMessage()
		if err != nil {
			common.HandleError(err, "ws.ReadMessage")
			return err
		}

		fmt.Println(string(message))
		sendToClient(client, message)

		//Response message to client
		err = ws.WriteMessage(mt, message)
		if err != nil {
			common.HandleError(err, "ws.WriteMessage")
			return err
		}
	}
}

func sendToClient(c *Client, msg []byte) {
	_, err := c.Conn.Write(msg)
	common.HandleError(err, "sendToClient")

}

type Client struct {
	UserName string
	Messages chan []byte
	Address  *net.UDPAddr
	Conn     *net.UDPConn
	cache    *cache.Redis
	Config
}

type Config struct {
	Port       int
	Address    string
	BufferSize int
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
	c.Messages = make(chan []byte, 20)
	c.Address = conAddr

	go c.sendMessage()

	for {
		c.handleMessage()
	}
}

func (c *Client) sendMessage() {
	for {
		msg := <-c.Messages

		message := domains.Message{}
		err := json.Unmarshal(msg, &message)
		if err != nil {
			common.HandleError(err, "sendMessage Unmarshal")
		}

		// send to server
		log.Printf("\n[%v] %s: %s", message.Time.Format("15:04:05"), message.UserName, message.Content)
	}
}

func (c *Client) handleMessage() {
	message := make([]byte, c.BufferSize)
	rlen, err := c.Conn.Read(message[0:])
	if err != nil {
		common.HandleError(err, "handleMessage c.Conn.Read")
	}

	data := strings.TrimSpace(string(message[0:rlen]))
	fmt.Println("message: " + data)
	msg := domains.Message{}
	err = json.Unmarshal([]byte(data), &msg)
	if err != nil {
		common.HandleError(err, "handleMessage json.Unmarshal")
	}

	msg.Time = time.Now()
	msg.ID = common.GenerateUUID()
	msgByte, err := json.Marshal(msg)
	if err != nil {
		common.HandleError(err, "handleMessage json.Marshal")
	}

	c.Messages <- msgByte
}
