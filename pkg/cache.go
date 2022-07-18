package pkg

import (
	"fmt"
	lru "github.com/hashicorp/golang-lru"
	"github.com/lovoo/goka/storage"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

// CachedStorage wraps a Goka Storage with an LRU cache.
type CachedStorage struct {
	storage.Storage
	cache *lru.Cache
}

func NewCachedStorage(size int, st storage.Storage) (storage.Storage, error) {
	c, err := lru.New(size)
	if err != nil {
		return nil, fmt.Errorf("error creating cache: %v", err)
	}
	return &CachedStorage{st, c}, nil
}

func (cs *CachedStorage) Has(key string) (bool, error) {
	if cs.cache.Contains(key) {
		return true, nil
	}
	return cs.Storage.Has(key)
}

func (cs *CachedStorage) Get(key string) ([]byte, error) {
	if val, has := cs.cache.Get(key); has {
		return val.([]byte), nil
	}
	val, err := cs.Storage.Get(key)
	cs.cache.Add(key, val)
	return val, err
}

func (cs *CachedStorage) Set(key string, value []byte) error {
	defer cs.cache.Remove(key)
	return cs.Storage.Set(key, value)
}

func (cs *CachedStorage) Delete(key string) error {
	defer cs.cache.Remove(key)
	return cs.Storage.Delete(key)
}

func CachedStorageBuilder(size int, path string, opts *opt.Options) storage.Builder {
	builder := storage.BuilderWithOptions(path, opts)
	return func(topic string, partition int32) (storage.Storage, error) {
		st, err := builder(topic, partition)
		if err != nil {
			return nil, err
		}
		return NewCachedStorage(size, st)
	}
}
