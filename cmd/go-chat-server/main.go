package main

import (
	"log"
	"strconv"

	"github.com/albertojnk/go-chat/common"
	"github.com/albertojnk/go-chat/internal/cache"
	"github.com/albertojnk/go-chat/internal/server"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	cacheConfig := cache.NewRedisConfig()
	cache, err := cache.NewRedis(cacheConfig)
	if err != nil {
		log.Fatalf("error while connecting redis: %v, err: %v", cache, err)
	}

	strPort := common.GetEnv("UDP_LISTENING_PORT", "3000")
	port, err := strconv.Atoi(strPort)
	if err != nil {
		log.Fatalf("environment variable UDP_LISTENING_PORT is wrong: %s", strPort)
		port = 3000 //this makes 3000 the default port if there is a parse error
	}
	s := server.UDPServer{
		Config: server.Config{
			Port:       port,
			Address:    common.GetEnv("UDP_LISTENING_ADDRESS", "0.0.0.0"),
			BufferSize: 2048,
		},
	}
	s.NewUDP()
}
