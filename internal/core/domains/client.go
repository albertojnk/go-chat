package domains

import (
	"net"

	"github.com/albertojnk/go-chat/internal/cache"
)

type Client struct {
	UserName string
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
