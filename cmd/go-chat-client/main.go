package main

import (
	"log"
	"strconv"

	"github.com/albertojnk/go-chat/common"
	"github.com/albertojnk/go-chat/internal/client"
)

func main() {

	addr := common.GetEnv("UDP_LISTENING_ADDRESS", "0.0.0.0")
	strPort := common.GetEnv("UDP_LISTENING_PORT", "3000")
	port, err := strconv.Atoi(strPort)
	if err != nil {
		log.Fatalf("environment variable UDP_LISTENING_PORT is wrong: %s", strPort)
		port = 3000 //this makes 3000 the default port if there is a parse error
	}

	c := client.Client{
		Config: client.Config{
			Port:       port,
			Address:    addr,
			BufferSize: 2048,
		},
	}

	c.NewUDP()

	return
}
