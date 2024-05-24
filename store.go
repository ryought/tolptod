package main

import "fmt"

type CacheStore struct {
	i       int
	entries []Entry
}

type Entry struct {
	// c      Cache
	Id     string `json:"id"`
	Config string `json:"config"`
	Status string `json:"status"`
}

func NewCacheStore() CacheStore {
	return CacheStore{
		i:       0,
		entries: make([]Entry, 0),
	}
}

func (s *CacheStore) List() []Entry {
	fmt.Println("list", s)
	return s.entries
}

func (s CacheStore) Get(id string) (Entry, bool) {
	for _, entry := range s.entries {
		if entry.Id == id {
			return entry, true
		}
	}
	var empty Entry
	return empty, false
}

func (s *CacheStore) Request(config string) string {
	fmt.Println("request", config, s)
	id := fmt.Sprintf("%d", s.i)
	entry := Entry{
		Id:     id,
		Config: config,
		Status: "pending",
	}
	s.entries = append(s.entries, entry)
	s.i += 1
	return id
}
