package main

import (
	"context"
	"fmt"
	"time"
)

type CacheStore struct {
	i       int
	entries []*Entry
}

type Entry struct {
	// c      Cache
	Id       string `json:"id"`
	Config   string `json:"config"`
	Status   string `json:"status"`
	Progress int    `json:"progress"`
	Cache    int    `json:"cache"`
	cancel   func()
}

func NewCacheStore() CacheStore {
	return CacheStore{
		i:       0,
		entries: make([]*Entry, 0),
	}
}

func (s *CacheStore) List() []*Entry {
	fmt.Println("list", s)
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
	fmt.Println("cancel", id)
	for _, entry := range s.entries {
		if entry.Id == id {
			fmt.Println("call cancel()")
			entry.cancel()
			return true
		}
	}
	return false
}

// Start
func (s *CacheStore) Request(config string) string {
	fmt.Println("request", config, s)
	id := fmt.Sprintf("%d", s.i)
	ctx, cancel := context.WithCancel(context.Background())
	entry := Entry{
		Id:       id,
		Config:   config,
		Status:   "pending",
		Progress: 0,
		Cache:    0,
		cancel:   cancel,
	}
	s.entries = append(s.entries, &entry)
	go func() {
		cache := HeavyTask(ctx, func(i int) {
			entry.Progress = i
			fmt.Println("progress", i)
		})
		entry.Cache = cache
		entry.Status = "done"
	}()
	s.i += 1
	return id
}

func HeavyTask(ctx context.Context, onProgress func(int)) int {
	fmt.Println("heavy task start")
	for i := 0; i < 30; i++ {
		select {
		case <-ctx.Done():
			fmt.Println("canceled", i)
			return -1
		default:
			fmt.Println("not canceled", i)
		}
		time.Sleep(time.Second)
		onProgress(i)
	}
	fmt.Println("heavy task end")
	return 111
}
