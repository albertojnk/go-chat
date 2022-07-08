package handlers

import (
	"context"
	"net/http"

	"github.com/albertojnk/go-chat/internal/core/domains"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

func (h *NETHandler) LoginHandler(ctx *context.Context, c *gin.Context) error {

	body := gin.H{}

	c.HTML(http.StatusOK, "index.html", body)

	return nil
}

func (h *NETHandler) ChatHandler(ctx *context.Context, c *gin.Context) error {

	body := gin.H{}

	// get users
	users, err := h.cache.GetValue(domains.USERSKEY)
	if err != nil && err != redis.Nil {
		return err
	}
	// get messages
	msgs, err := h.cache.GetValue(domains.MESSAGESKEY)
	if err != nil && err != redis.Nil {
		return err
	}

	body["users"] = users
	body["messages"] = msgs

	c.HTML(http.StatusOK, "chat.html", body)

	return nil
}
