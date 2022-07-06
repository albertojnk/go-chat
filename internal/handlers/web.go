package handlers

import (
	"context"
	"net/http"

	"github.com/albertojnk/go-chat/internal/core/domains"
	"github.com/gin-gonic/gin"
)

func LoginHandler(ctx *context.Context, c *gin.Context) error {

	body := gin.H{}

	c.HTML(http.StatusOK, "index.html", body)

	return nil
}

func ChatHandler(ctx *context.Context, c *gin.Context) error {

	body := gin.H{}

	msgs := []domains.Message{}

	body["messages"] = msgs

	c.HTML(http.StatusOK, "chat.html", body)

	return nil
}
