package handlers

import "github.com/albertojnk/go-chat/internal/cache"

type NETHandler struct {
	cache *cache.Redis
}

func NewNETHandler(cache *cache.Redis) *NETHandler {
	return &NETHandler{
		cache: cache,
	}
}
