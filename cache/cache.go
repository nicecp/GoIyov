package cache

import (
	"GoIyov/singleflight"
	"crypto/tls"
	"sync"
)

type Cache struct {
	m           sync.Map
	singleGroup *singleflight.Group
}

func NewCache() *Cache {
	return &Cache{
		singleGroup: &singleflight.Group{},
	}
}

func (cache *Cache) GetOrStore(key string, fn func() (interface{}, error)) (interface{}, error) {
	if val, ok := cache.m.Load(key); ok {
		return val.(tls.Certificate), nil
	}

	cert, err := cache.singleGroup.Do(key, fn)
	if err != nil {
		return nil, err
	}
	cache.m.Store(key, cert)
	return cert.(tls.Certificate), nil
}

func (cache *Cache) GetCache() sync.Map {
	return cache.m
}
