package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	_cache "github.com/albertojnk/go-chat/internal/cache"
	"github.com/google/uuid"
)

type Client struct {
	ID      uuid.UUID
	Name    string
	Address *net.UDPAddr
	cache   *_cache.Redis
	Config
}

type Config struct {
	Port       int
	Address    string
	BufferSize int
}

func (c *Client) NewUDP() {

	conn, err := net.Dial("udp", fmt.Sprintf("%s:%v", c.Config.Address, c.Config.Port))
	if err != nil {
		fmt.Printf("Some error %v", err)
		return
	}
	defer conn.Close()

	fmt.Fprintf(conn, "connected")

	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)
		fmt.Fprintf(conn, text)
	}
}
