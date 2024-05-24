package main

import (
	"context"
	"fmt"
)

type CacheStore struct {
	i       int
	entries []*Entry
}

type Entry struct {
	Id       string      `json:"id"`
	Config   CacheConfig `json:"config"`
	Status   string      `json:"status"`
	Progress int         `json:"progress"`
	cache    *Cache
	cancel   func()
}

type CacheConfig struct {
	X       int `json:"x"`
	Y       int `json:"y"`
	K       int `json:"k"`
	Bin     int `json:"bin"`
	FreqLow int `json:"freqLow"`
	FreqUp  int `json:"freqUp"`
}

func NewCacheStore() CacheStore {
	return CacheStore{
		i:       0,
		entries: make([]*Entry, 0),
	}
}

func (s *CacheStore) List() []*Entry {
	return s.entries
}

func (s CacheStore) Get(id string) (Entry, bool) {
	for _, entry := range s.entries {
		if entry.Id == id {
			return *entry, true
		}
	}
	var empty Entry
	return empty, false
}

// Stop
func (s CacheStore) Cancel(id string) bool {
	for _, entry := range s.entries {
		if entry.Id == id {
			entry.cancel()
			return true
		}
	}
	return false
}

// Start
func (s *CacheStore) Request(index *IndexV2, config CacheConfig) string {
	id := fmt.Sprintf("%d", s.i)
	ctx, cancel := context.WithCancel(context.Background())
	entry := Entry{
		Id:       id,
		Config:   config,
		Status:   "pending",
		Progress: 0,
		cache:    nil,
		cancel:   cancel,
	}
	s.entries = append(s.entries, &entry)
	go func() {
		xindex := index.xindex[config.X]
		yindex := index.yindex[config.Y]
		cache := NewCache(
			ctx,
			xindex,
			yindex,
			Config{
				k:            config.K,
				bin:          config.Bin,
				freqLow:      config.FreqLow,
				freqUp:       config.FreqUp,
				localFreqLow: 0,
				localFreqUp:  0,
				xL:           0,
				xR:           xindex.N,
				yL:           0,
				yR:           yindex.N,
			},
			func(y, yL, yR int) {
				entry.Progress = 100 * (y - yL) / (yR - yL)
			},
		)
		entry.cache = &cache
		entry.Status = "done"
	}()
	s.i += 1
	return id
}
