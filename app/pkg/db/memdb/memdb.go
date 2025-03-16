package memdb

import (
	"sync"
	"url_shortener/pkg/urlgenerator"
)

type URLStorage struct {
	urlsize     uint
	longToShort map[string]string
	shortToLong map[string]string
	mu          *sync.Mutex
}

func NewURLStorage(urlsize uint) *URLStorage {
	longToShort := make(map[string]string)
	shortToLong := make(map[string]string)
	return &URLStorage{urlsize: urlsize, longToShort: longToShort, shortToLong: shortToLong, mu: &sync.Mutex{}}
}

func (s *URLStorage) AddURL(url string) string {
	shortURL := urlgenerator.Generate(s.urlsize)
	s.mu.Lock()
	if s.longToShort[url] != "" {
		return s.longToShort[url] //race
	}
	if s.shortToLong[shortURL] != "" {
		s.mu.Unlock()
		return s.AddURL(url) //collision
	}
	s.longToShort[url] = shortURL
	s.shortToLong[shortURL] = url
	s.mu.Unlock()
	return shortURL
}

func (s *URLStorage) GetOriginal(url string) string {
	return s.shortToLong[url]
}

func (s *URLStorage) GetShort(url string) string {
	return s.longToShort[url]
}

type DB struct {
	storage *URLStorage
}

func NewDB(urlsize uint) *DB {
	return &DB{storage: NewURLStorage(urlsize)}
}

func (db *DB) MakeShort(url string) (string, error) {
	short := db.storage.GetShort(url)
	if short != "" {
		return short, nil
	}
	short = db.storage.AddURL(url)
	return short, nil
}

func (db *DB) GetOriginal(url string) (string, error) {
	originalURL := db.storage.GetOriginal(url)
	return originalURL, nil
}

func (db *DB) URLSize() uint {
	return db.storage.urlsize
}
