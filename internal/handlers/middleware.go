package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

//HandlerAPI .
func HandlerAPI(f func(ctx *context.Context, c *gin.Context) error) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		ctx := context.Background()
		if err := f(&ctx, c); err != nil {
			log.Println(fmt.Errorf("%v: %s", err, debug.Stack()))
			http.Error(c.Writer, "Internal Server Error", http.StatusInternalServerError)
		}
	})
}

// HandlerPage provide abstractionn for page handler receive the context of application.
func HandlerPage(f func(ctx *context.Context, c *gin.Context) error) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		ctx := context.Background()

		if err := f(&ctx, c); err != nil {
			log.Println(fmt.Errorf("%v: %s", err, debug.Stack()))
			c.Redirect(302, "/error")
		}
	})
}
