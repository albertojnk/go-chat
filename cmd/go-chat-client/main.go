package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/albertojnk/go-chat/common"
	"github.com/albertojnk/go-chat/internal/cache"
	"github.com/albertojnk/go-chat/internal/handlers"
	"github.com/gin-gonic/gin"
)

func main() {

	cacheConfig := cache.NewRedisConfig()
	cache, err := cache.NewRedis(cacheConfig)
	if err != nil {
		log.Fatalf("error while connecting redis: %v, err: %v", cache, err)
	}

	ctx := context.Background()
	router := gin.New()

	router.Static("/assets", common.GetEnv("WEB_ASSETS", "../../web/assets"))
	router.LoadHTMLGlob(common.GetEnv("WEB_PAGES", "../../web/pages/*"))

	hdl := handlers.NewNETHandler(cache)

	router.GET("/", handlers.HandlerPage(hdl.LoginHandler))
	router.GET("/chat/:username", handlers.HandlerPage(hdl.ChatHandler))

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
