// Copyright 2019 Drone IO, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package netrc

import (
	"context"
	"time"

	"github.com/drone/drone/core"

	lru "github.com/hashicorp/golang-lru"
)

type cacheEntry struct {
	expiry time.Time
	netrc  *core.Netrc
}

type cacher struct {
	base  core.NetrcService
	ttl   time.Duration
	cache *lru.Cache
}

func NewCache(base core.NetrcService, size int, ttl time.Duration) core.NetrcService {
	cache, _ := lru.New(size)
	return &cacher{
		base:  base,
		cache: cache,
		ttl:   ttl,
	}
}

func (c *cacher) Create(ctx context.Context, user *core.User, repo *core.Repository) (*core.Netrc, error) {
	now := time.Now()

	cached, ok := c.cache.Get(repo.ID)
	if ok {
		entry := cached.(*cacheEntry)
		if now.Before(entry.expiry) {
			return entry.netrc, nil
		}
		c.cache.Remove(repo.ID)
	}

	netrc, err := c.base.Create(ctx, user, repo)
	if err != nil {
		return nil, err
	}

	c.cache.Add(repo.ID, &cacheEntry{
		expiry: now.Add(c.ttl),
		netrc:  netrc,
	})

	return netrc, nil
}