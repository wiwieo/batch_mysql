package server

import (
	"wiwieo/batch_mysql/adapter"
	"wiwieo/batch_mysql/cache"
)

type srv struct {
	AdCache         *cache.Cache
	DebrisTypeCache *cache.Cache
}

func NewSrv() *srv {
	var adCache adapter.Ad
	c, err := cache.NewCache(adCache, "ad")
	if err != nil {
		panic(err)
	}

	var dtCache adapter.DebrisType
	d, err := cache.NewCache(dtCache, "debris_type")
	if err != nil {
		panic(err)
	}

	return &srv{
		AdCache:         c,
		DebrisTypeCache: d,
	}
}
