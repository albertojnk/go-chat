package domains

import (
	"net"

	"github.com/albertojnk/go-chat/internal/cache"
)

type Client struct {
	UserName string       `json:"username"`
	Address  *net.UDPAddr `json:"-"`
	Conn     *net.UDPConn `json:"-"`
	cache    *cache.Redis `json:"-"`
	Config   `json:"-"`
}

type Config struct {
	Port       int
	Address    string
	BufferSize int
}
