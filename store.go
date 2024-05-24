package main

import (
	"context"
	"fmt"
	"slices"
	"strconv"
)

type CacheStore struct {
	i       int
	entries map[string]*Entry
}

type Entry struct {
	Id       string      `json:"id"`
	Config   CacheConfig `json:"config"`
	Done     bool        `json:"done"`
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
		entries: make(map[string]*Entry),
	}
}

func (s *CacheStore) List() []*Entry {
	ret := make([]*Entry, 0)
	for _, entry := range s.entries {
		ret = append(ret, entry)
	}
	slices.SortFunc(ret, func(a, b *Entry) int {
		ai, _ := strconv.Atoi(a.Id)
		bi, _ := strconv.Atoi(b.Id)
		return ai - bi
	})
	return ret
}

func (s *CacheStore) Get(id string) (Entry, bool) {
	entry, ok := s.entries[id]
	return *entry, ok
}

// Stop
func (s *CacheStore) Cancel(id string) bool {
	entry, ok := s.entries[id]

	// not found
	if !ok {
		return false
	}

	// stop
	if !entry.Done {
		entry.cancel()
	}
	// delete entry
	delete(s.entries, id)

	return true
}

// Start
func (s *CacheStore) Request(index *IndexV2, config CacheConfig) string {
	id := fmt.Sprintf("%d", s.i)
	ctx, cancel := context.WithCancel(context.Background())
	entry := Entry{
		Id:       id,
		Config:   config,
		Done:     false,
		Progress: 0,
		cache:    nil,
		cancel:   cancel,
	}
	s.entries[id] = &entry
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
				p := 100 * (y - yL) / (yR - yL)
				entry.Progress = p
				// fmt.Println("progress", p)
			},
		)
		entry.cache = &cache
		entry.Done = true
	}()
	s.i += 1
	return id
}
