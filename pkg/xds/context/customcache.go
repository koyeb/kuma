package context

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

func NewRedisCache(address string, defaultTTL time.Duration) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	ctx := context.Background()
	_, err := client.Get(ctx, "whatever").Result()
	if err != redis.Nil {
		return nil, errors.Wrapf(err, "could not contact redis at address %s", address)
	}

	return &RedisCache{
		client:     client,
		defaultTTL: defaultTTL,
	}, nil
}

type RedisCache struct {
	client     *redis.Client
	defaultTTL time.Duration
}

func (gc *GlobalContext) CacheKeyName(_ string) string {
	// shared between all meshes
	return "global-context"
}

func (gc *BaseMeshContext) CacheKeyName(meshName string) string {
	return fmt.Sprintf("base-context-mesh-%s", meshName)
}

type cachable interface {
	CacheKeyName(meshName string) string
}

func (rc *RedisCache) Set(ctx context.Context, meshName string, item cachable) error {
	if rc == nil {
		return nil
	}

	key := item.CacheKeyName(meshName)

	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)

	err := e.Encode(item)
	if err != nil {
		return err
	}

	err = rc.client.Set(ctx, key, b.Bytes(), rc.defaultTTL).Err()
	if err != nil {
		return err
	}

	return nil
}

func (rc *RedisCache) Get(ctx context.Context, meshName string, receiver cachable) (bool, error) {
	if rc == nil {
		return false, nil
	}

	key := receiver.CacheKeyName(meshName)

	data, err := rc.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}

		return false, err
	}

	b := bytes.Buffer{}
	b.Write(data)
	d := gob.NewDecoder(&b)
	err = d.Decode(receiver)
	if err != nil {
		return false, err
	}

	return true, nil
}
