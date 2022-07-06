package router

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/albertojnk/go-chat/common"
	"github.com/albertojnk/go-chat/internal/handlers"
	"github.com/gin-gonic/gin"
)

func Serve() {

	ctx := context.Background()
	router := gin.New()

	router.Static("/assets", "../../web/assets")
	router.LoadHTMLGlob("../../web/pages/*")

	hdl := handlers.NewNETHandler()

	router.GET("/", handlers.HandlerPage(handlers.LoginHandler))
	router.GET("/chat/:username", handlers.HandlerPage(handlers.ChatHandler))

	router.GET("/chat", handlers.HandlerAPI(hdl.ChatWSHandler))

	port := common.GetEnv("PORT", "9000")
	s := &http.Server{
		Addr:    fmt.Sprintf(":%v", port),
		Handler: router,
	}

	done := make(chan struct{})
	doneTLS := make(chan struct{})
	go func() {
		<-ctx.Done()
		if err := s.Shutdown(ctx); err != nil {
			log.Fatal(err.Error())
		}
		close(done)
	}()

	close(doneTLS)
	log.Printf("Serving web application at http://0.0.0.0:%v", port)
	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err.Error())
	}

	<-doneTLS
	<-done

	return
}
