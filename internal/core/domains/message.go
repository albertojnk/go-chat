package domains

import (
	"net"
	"time"
)

type Message struct {
	ID          string       `json:"id"`
	UserName    string       `json:"username"`
	Content     string       `json:"content"`
	Time        time.Time    `json:"time"`
	Address     *net.UDPAddr `json:"address"`
	MessageType `json:"message_type"`
	Clients     map[string]Client `json:"clients"`
}

type MessageType string

const (
	HANDSHAKE      MessageType = "HANDSHAKE"
	GOODBYE        MessageType = "GOODBYE"
	MESSAGE        MessageType = "MESSAGE"
	DELETEMESSAGE  MessageType = "DELETEMESSAGE"
	INVALIDMESSAGE MessageType = "INVALIDMESSAGE"
)
